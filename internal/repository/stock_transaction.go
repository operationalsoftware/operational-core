package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type StockTrxRepository struct{}

func NewStockTrxRepository() *StockTrxRepository {
	return &StockTrxRepository{}
}

func (r *StockTrxRepository) GetStockLevels(ctx context.Context, exec db.PGExecutor, input *model.GetStockLevelsInput) ([]model.StockLevel, error) {

	query := `
WITH RankedStock AS (
	SELECT
		ste.stock_transaction_id,
		ste.account,
		ste.stock_code,
		ste.location,
		ste.bin,
		ste.lot_number,
		ste.running_total AS stock_level,
		st.timestamp,
		ROW_NUMBER() OVER (
			PARTITION BY ste.account, ste.stock_code, ste.location, ste.bin, COALESCE(ste.lot_number, 'NO_LOT')
			ORDER BY st.timestamp DESC, ste.stock_transaction_id DESC
		) AS rn
	FROM 
		stock_transaction_entry ste
	JOIN
		stock_transaction st ON ste.stock_transaction_id = st.stock_transaction_id
	WHERE
		($1::text IS NULL OR ste.account = $1::text)
		AND
		($2::text IS NULL OR ste.stock_code = $2::text)
		AND
		($3::text IS NULL OR ste.location = $3::text)
		AND
		($4::text IS NULL OR ste.bin = $4::text)
		AND
		($5::text IS NULL OR ste.lot_number = $5::text)
		AND
		($6::timestamp IS NULL OR st.timestamp <= $6::timestamp)  -- If LTETimestamp provided, only match up until it
)

SELECT
	account,
	stock_code,
	location,
	bin,
	lot_number,
	stock_level,
	timestamp
FROM
	RankedStock
WHERE 
	rn = 1
	AND
	stock_level <> 0
ORDER BY
	stock_transaction_id DESC
LIMIT $7 OFFSET $8
    `

	// Prepare params with NULL handling
	accountParam := input.Account
	stockCodeParam := input.StockCode
	locationParam := input.Location
	binParam := input.Bin
	lotNumberParam := input.LotNumber
	lteTimestampParam := input.LTETimestamp
	// accountParam := db.NullStringToParam(input.Account)
	// stockCodeParam := db.NullStringToParam(input.StockCode)
	// locationParam := db.NullStringToParam(input.Location)
	// binParam := db.NullStringToParam(input.Bin)
	// lotNumberParam := db.NullStringToParam(input.LotNumber)
	// lteTimestampParam := db.NullTimeToParam(input.LTETimestamp)

	limit := 1000
	if input.PageSize > 0 {
		limit = input.PageSize
	}
	offset := 0
	if input.Page > 0 {
		offset = (input.Page - 1) * input.PageSize
	}

	// Prepare the query and execute it with the provided filters
	rows, err := exec.Query(ctx, query,
		accountParam,
		stockCodeParam,
		locationParam,
		binParam,
		lotNumberParam,
		lteTimestampParam,

		limit, offset, // Pagination parameters
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.StockLevel
	for rows.Next() {
		var sl model.StockLevel
		err := rows.Scan(
			&sl.Account,
			&sl.StockCode,
			&sl.Location,
			&sl.Bin,
			&sl.LotNumber,
			&sl.StockLevel,
			&sl.Timestamp,
		)
		if err != nil {
			return nil, err
		}

		// temp do some rounding of Decimals to correct floating point representation errors
		// round to 8 d.p. to align with SyteLine
		sl.StockLevel = sl.StockLevel.Round(8)

		results = append(results, sl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *StockTrxRepository) PostStockTransactions(
	ctx context.Context,
	exec pgx.Tx,
	transactions *model.PostStockTransactionsInput,
	userID int,
) error {

	insertTrxQuery := `
INSERT INTO
	stock_transaction (
		stock_transaction_type,
		transaction_by,
		transaction_note,
		timestamp
	)
VALUES (
	$1, $2, $3, COALESCE($4, NOW())
)
	RETURNING 
		stock_transaction_id, timestamp
	`

	// Combined insert query for both From and To transactions
	insertEntriesQuery := `
-- $1 = stock transaction id
-- $2 = timestamp
-- $3 = stock_code
-- $4 = quantity
-- $5 = from_account
-- $6 = from_location
-- $7 = from_bin
-- $8 = from_lot_number
-- $9 = to_account
-- $10 = to_location
-- $11 = to_bin
-- $12 = to_lot_number

INSERT INTO stock_transaction_entry (
	stock_transaction_id,
	quantity,
	stock_code, 
	account, 
	location, 
	bin,
	lot_number,
	running_total
)
VALUES 
	-- From entry
	($1, -1 * $3, $2, $4, $5, $6, $7,
	COALESCE($12, 0) - $3),

	-- To entry
	($1, $3, $2, $8, $9, $10, $11,
	COALESCE($13, 0) + $3)
`

	// Update future RunningTotal for the From and To accounts
	updateQuery := `
	UPDATE
	    stock_transaction_entry ste
	SET
	    running_total = running_total + $1
	FROM
	    stock_transaction st
	WHERE
	    ste.stock_transaction_id = st.stock_transaction_id
	    AND ste.account = $4
	    AND ste.stock_code = $3
	    AND ste.location = $5
	    AND ste.bin = $6
	    AND (ste.lot_number = $7 OR ste.lot_number IS NULL)
	    AND st.timestamp > $2`

	for _, t := range *transactions {
		// Fetch previous running_total
		var fromRunningTotal, toRunningTotal decimal.Decimal

		err := exec.QueryRow(ctx, `
SELECT running_total
	FROM 
		stock_transaction_entry
		INNER JOIN stock_transaction
		USING (stock_transaction_id)
	WHERE 
		account = $1
		AND stock_code = $2
		AND location = $3
		AND bin = $4
		AND (lot_number IS NOT DISTINCT FROM $5)
		AND timestamp <= $6
	ORDER BY
		timestamp DESC, stock_transaction_id DESC
	LIMIT 1
`,
			t.FromAccount,
			t.StockCode,
			t.FromLocation,
			t.FromBin,
			t.FromLotNumber,
			t.Timestamp).Scan(&fromRunningTotal)
		if err != nil && err != pgx.ErrNoRows {
			return fmt.Errorf("failed to fetch running_total for FROM entry: %v", err)
		}

		err = exec.QueryRow(ctx, `
SELECT running_total
	FROM stock_transaction_entry
		INNER JOIN stock_transaction
		USING (stock_transaction_id)
	WHERE account = $1
		AND stock_code = $2
		AND location = $3
		AND bin = $4
		AND (lot_number IS NOT DISTINCT FROM $5)
		AND timestamp <= $6
	ORDER BY timestamp DESC, stock_transaction_id DESC
	LIMIT 1
`,
			t.ToAccount,
			t.StockCode,
			t.ToLocation,
			t.ToBin,
			t.ToLotNumber,
			t.Timestamp).Scan(&toRunningTotal)
		if err != nil && err != pgx.ErrNoRows {
			return fmt.Errorf("failed to fetch running_total for TO entry: %v", err)
		}

		// newFromRunningTotal := prevFromRunningTotal.Sub(t.Qty)

		var stockTrxID int
		var stockTrxTimestamp time.Time
		err = exec.QueryRow(ctx, insertTrxQuery, "STOCK", userID, "This is test note", t.Timestamp).Scan(&stockTrxID, &stockTrxTimestamp)
		if err != nil {
			return fmt.Errorf("failed to create stock transaction: %v", err)
		}

		// IMPORTANT: posting from and to the same Account, Location, Bin and
		// LotNumber is not allowed as it is not compatible with how running totals
		// are calculated
		if t.FromAccount == t.ToAccount &&
			t.FromLocation == t.ToLocation &&
			t.FromBin == t.ToBin &&
			((t.FromLotNumber != nil && t.ToLotNumber != nil && t.FromLotNumber == t.ToLotNumber) ||
				(t.FromLotNumber != nil == false && t.ToLotNumber != nil == false)) {
			return fmt.Errorf("Account, Location, Bin and LotNumber cannot be the same for 'From' and 'To'")
		}

		if t.Qty.Equal(decimal.Zero) {
			fmt.Printf("Skipping qty 0 transaction of %s\n", t.StockCode)
			continue
		}

		_, err = exec.Exec(ctx, insertEntriesQuery,
			stockTrxID,
			// t.Timestamp,
			t.StockCode,
			t.Qty,
			// From transaction
			t.FromAccount,
			t.FromLocation,
			t.FromBin,
			t.FromLotNumber,

			// To transaction
			t.ToAccount,
			t.ToLocation,
			t.ToBin,
			t.ToLotNumber,
			fromRunningTotal,
			toRunningTotal,
		)

		if err != nil {
			return fmt.Errorf("failed to insert transactions: %v", err)
		}

		_, err = exec.Exec(ctx, updateQuery,
			t.Qty.Mul(decimal.NewFromInt(-1)), // negative for "From"
			t.Timestamp,
			t.FromAccount,
			t.StockCode,
			t.FromLocation,
			t.FromBin,
			t.FromLotNumber,
		)

		if err != nil {
			return fmt.Errorf("failed to update future RunningTotal for From account: %v", err)
		}

		_, err = exec.Exec(ctx, updateQuery,
			t.Qty, // positive for "To"
			t.Timestamp,
			t.ToAccount,
			t.StockCode,
			t.ToLocation,
			t.ToBin,
			t.ToLotNumber,
		)

		if err != nil {
			return fmt.Errorf("failed to update future RunningTotal for To account: %v", err)
		}
	}

	return nil
}

func (r *StockTrxRepository) GetStockTransactions(ctx context.Context, exec db.PGExecutor, input *model.GetTransactionsInput) ([]model.StockTransactionEntry, error) {

	// Base query
	query := `
WITH matched_tx_ids AS (
	SELECT DISTINCT ste.stock_transaction_id
	FROM stock_transaction_entry ste
	JOIN stock_transaction st ON st.stock_transaction_id = ste.stock_transaction_id
	WHERE
		($1::text IS NULL OR ste.account = $1::text) AND
		($2::text IS NULL OR ste.stock_code = $2::text) AND
		($3::text IS NULL OR ste.location = $3::text) AND
		($4::text IS NULL OR ste.bin = $4::text) AND
		($5::text IS NULL OR ste.lot_number = $5::text) AND
		($6::timestamp IS NULL OR st.timestamp <= $6::timestamp)
)

SELECT 
	ste.stock_transaction_entry_id,
	st.stock_transaction_type,
	ste.account,
	ste.stock_code,
	ste.location,
	ste.bin,
	ste.quantity,
	ste.lot_number,
	ste.running_total,
	st.transaction_by,
	u.username AS transaction_by_username,
	st.timestamp,
	ste.stock_transaction_id
FROM stock_transaction_entry ste
JOIN stock_transaction st ON st.stock_transaction_id = ste.stock_transaction_id
LEFT JOIN app_user u ON u.user_id = st.transaction_by
JOIN matched_tx_ids m ON m.stock_transaction_id = ste.stock_transaction_id
ORDER BY ste.stock_transaction_id DESC, ste.stock_transaction_entry_id ASC
LIMIT $7 OFFSET $8
	`

	// Pagination
	limit := input.PageSize
	if limit == 0 {
		limit = 1000
	}
	offset := (input.Page - 1) * limit

	account := input.Account
	stockCode := input.StockCode
	location := input.Location
	bin := input.Bin
	lotNumber := input.LotNumber
	lteTimestamp := input.LTETimestamp

	// Execute
	rows, err := exec.Query(ctx, query,
		account,
		stockCode,
		location,
		bin,
		lotNumber,
		lteTimestamp,
		limit,
		offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect results
	var transactions []model.StockTransactionEntry
	for rows.Next() {
		var st model.StockTransactionEntry
		err := rows.Scan(
			&st.StockTransactionEntryID, // 1
			&st.StockTransactionType,    // 2
			&st.Account,                 // 3
			&st.StockCode,               // 4
			&st.Location,                // 5
			&st.Bin,                     // 6
			&st.Quantity,                // 7
			&st.LotNumber,               // 8
			&st.RunningTotal,            // 9
			&st.TransactionBy,           // 10
			&st.TransactionByUsername,   // 11
			&st.Timestamp,               // 12
			&st.StockTransactionID,
		)
		if err != nil {
			return nil, err
		}

		// temp do some rounding of Decimals to correct floating point representation errors
		// round to 8 d.p. to align with SyteLine
		st.Quantity = st.Quantity.Round(8)
		st.RunningTotal = st.RunningTotal.Round(8)

		transactions = append(transactions, st)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil

}

// func GetTransactions(ctx context.Context, exec db.PGExecutor, input *m.GetTransactionsInput) ([]m.StockTransactionEntry, error) {
// 	query := `
// WITH WideTransactions AS (
//     -- Create a wide table with both "out" and "in" transactions in the same row
//     SELECT
//         outTxn.StockTransactionID AS OutTransactionID,
//         outTxn.Account AS OutAccount,
//         outTxn.StockCode AS OutStockCode,
//         outTxn.Location AS OutLocation,
//         outTxn.Bin AS OutBin,
//         outTxn.LotNumber AS OutLotNumber,
//         outTxn.Qty AS OutQty,
//         outTxn.RunningTotal AS OutRunningTotal,
//         outTxn.TransactionBy AS OutTransactionBy,
//         outTxn.Timestamp AS OutTimestamp,
//         inTxn.StockTransactionID AS InTransactionID,
//         inTxn.Account AS InAccount,
//         inTxn.StockCode AS InStockCode,
//         inTxn.Location AS InLocation,
//         inTxn.Bin AS InBin,
//         inTxn.LotNumber AS InLotNumber,
//         inTxn.Qty AS InQty,
//         inTxn.RunningTotal AS InRunningTotal,
//         inTxn.TransactionBy AS InTransactionBy,
//         inTxn.Timestamp AS InTimestamp,
// 		u.Username AS TransactionByUsername
//     FROM StockTransaction outTxn
//     INNER JOIN StockTransaction inTxn
//     ON outTxn.StockTransactionID + 1 = inTxn.StockTransactionID
//     AND outTxn.StockTransactionID % 2 = 1  -- Ensure we only join "out" to "in"
// 	INNER JOIN User u
// 	ON u.UserID = outTxn.TransactionBy
//     WHERE
// 		(
// 			-- Apply filters to "out" transaction
// 			(outTxn.Account = ? OR ? IS NULL)
// 			AND (outTxn.StockCode = ? OR ? IS NULL)
// 			AND (outTxn.Location = ? OR ? IS NULL)
// 			AND (outTxn.Bin = ? OR ? IS NULL)
// 			AND (outTxn.LotNumber = ? OR ? IS NULL)
// 			AND (? IS NULL OR strftime('%Y-%m-%d %H:%M', outTxn.Timestamp) <= strftime('%Y-%m-%d %H:%M', ?))
// 		) OR (
// 			-- Apply filters to "in" transaction
// 			(inTxn.Account = ? OR ? IS NULL)
// 			AND (inTxn.StockCode = ? OR ? IS NULL)
// 			AND (inTxn.Location = ? OR ? IS NULL)
// 			AND (inTxn.Bin = ? OR ? IS NULL)
// 			AND (inTxn.LotNumber = ? OR ? IS NULL)
// 			AND (? IS NULL OR strftime('%Y-%m-%d %H:%M', inTxn.Timestamp) <= strftime('%Y-%m-%d %H:%M', ?))
// 		)
// ),
// SplitTransactions AS (
//     -- Split the wide table into two rows: one for "out" and one for "in"
//     SELECT
//         OutTransactionID AS StockTransactionID,
//         OutAccount AS Account,
//         OutStockCode AS StockCode,
//         OutLocation AS Location,
//         OutBin AS Bin,
//         OutLotNumber AS LotNumber,
//         OutQty AS Qty,
//         OutRunningTotal AS RunningTotal,
//         OutTransactionBy AS TransactionBy,
// 		TransactionByUsername,
//         OutTimestamp AS Timestamp
//     FROM WideTransactions
//     WHERE OutTransactionID IS NOT NULL

//     UNION ALL

//     SELECT
//         InTransactionID AS StockTransactionID,
//         InAccount AS Account,
//         InStockCode AS StockCode,
//         InLocation AS Location,
//         InBin AS Bin,
//         InLotNumber AS LotNumber,
//         InQty AS Qty,
//         InRunningTotal AS RunningTotal,
//         InTransactionBy AS TransactionBy,
// 		TransactionByUsername,
//         InTimestamp AS Timestamp
//     FROM WideTransactions
//     WHERE InTransactionID IS NOT NULL
// )

// SELECT
//     StockTransactionID,
//     Account,
//     StockCode,
//     Location,
//     Bin,
//     LotNumber,
//     Qty,
//     RunningTotal,
//     TransactionBy,
// 	TransactionByUsername,
//     Timestamp
// FROM SplitTransactions
// ORDER BY StockTransactionID DESC
// LIMIT ? OFFSET ?
//     `

// 	// Prepare params with NULL handling
// 	accountParam := db.NullStringToParam(input.Account)
// 	stockCodeParam := db.NullStringToParam(input.StockCode)
// 	locationParam := db.NullStringToParam(input.Location)
// 	binParam := db.NullStringToParam(input.Bin)
// 	lotNumberParam := db.NullStringToParam(input.LotNumber)
// 	lteTimestampParam := db.NullTimeToParam(input.LTETimestamp)

// 	// Pagination defaults
// 	limit := input.PageSize
// 	if limit == 0 {
// 		limit = 1000
// 	}
// 	offset := (input.Page - 1) * limit

// 	// Execute query
// 	rows, err := exec.Query(ctx, query,
// 		accountParam, accountParam,
// 		stockCodeParam, stockCodeParam,
// 		locationParam, locationParam,
// 		binParam, binParam,
// 		lotNumberParam, lotNumberParam,
// 		lteTimestampParam, lteTimestampParam,

// 		accountParam, accountParam,
// 		stockCodeParam, stockCodeParam,
// 		locationParam, locationParam,
// 		binParam, binParam,
// 		lotNumberParam, lotNumberParam,
// 		lteTimestampParam, lteTimestampParam,

// 		limit, offset,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// Collect results
// 	var transactions []m.StockTransactionEntry
// 	for rows.Next() {
// 		var st m.StockTransactionEntry
// 		err := rows.Scan(
// 			&st.StockTransactionID,
// 			&st.Account,
// 			&st.StockCode,
// 			&st.Location,
// 			&st.Bin,
// 			&st.LotNumber,
// 			&st.Quantity,
// 			&st.RunningTotal,
// 			&st.TransactionBy,
// 			&st.TransactionByUsername,
// 			&st.Timestamp,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		// temp do some rounding of Decimals to correct floating point representation errors
// 		// round to 8 d.p. to align with SyteLine
// 		st.Quantity = st.Quantity.Round(8)
// 		st.RunningTotal = st.RunningTotal.Round(8)

// 		transactions = append(transactions, st)
// 	}

// 	// Check for errors after iteration
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return transactions, nil
// }
