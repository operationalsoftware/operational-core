package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"app/pkg/pgconv"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type StockTransactionRepository struct{}

func NewStockTransactionRepository() *StockTransactionRepository {
	return &StockTransactionRepository{}
}

func (r *StockTransactionRepository) GetStockLevels(
	ctx context.Context,
	exec db.PGExecutor,
	input *model.GetStockLevelsInput,
) ([]model.StockLevel, error) {

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
			PARTITION BY ste.account, ste.stock_code, ste.location, ste.bin, ste.lot_number
			ORDER BY st.timestamp DESC, ste.stock_transaction_id DESC
		) AS rn
	FROM
		stock_transaction_entry ste
		JOIN
			stock_transaction st ON ste.stock_transaction_id = st.stock_transaction_id
	WHERE
		($1 = '' OR ste.account = $1)
		AND
		($2 = '' OR ste.stock_code = $2)
		AND
		($3 = '' OR ste.location = $3)
		AND
		($4 = '' OR ste.bin = $4)
		AND
		($5 = '' OR ste.lot_number = $5)
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
		input.Account,
		input.StockCode,
		input.Location,
		input.Bin,
		input.LotNumber,
		pgconv.TimePtrToPGTimestamptz(input.LTETimestamp),

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

		results = append(results, sl)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *StockTransactionRepository) PostStockTransactions(
	ctx context.Context,
	exec pgx.Tx,
	transactions *model.PostStockTransactionsInput,
	userID int,
) error {

	query := `
/*
Parameter mapping:
$1    → transaction_type
$2    → stock_code
$3    → quantity
$4    → transaction_note
$5    → transaction_by
$6    → timestamp
$7    → from_account
$8    → from_location
$9    → from_bin
$10   → from_lot_number
$11   → to_account
$12   → to_location
$13   → to_bin
$14   → to_lot_number
*/

WITH inserted_tx AS (
    INSERT INTO stock_transaction (
        transaction_type, transaction_note, transaction_by, timestamp
    )
    VALUES ($1, $4, $5, COALESCE($6, NOW()))
    RETURNING stock_transaction_id, timestamp
),

inserted_from_entry AS (
    INSERT INTO stock_transaction_entry (
        stock_transaction_id, account, stock_code, location, bin, lot_number, quantity, running_total
    )
    SELECT inserted_tx.stock_transaction_id, $7, $2, $8, $9, $10, -1 * $3,
        COALESCE((
            SELECT running_total
            FROM stock_transaction_entry e
            JOIN stock_transaction t ON t.stock_transaction_id = e.stock_transaction_id
            WHERE e.account = $7
              AND e.stock_code = $2
              AND e.location = $8
              AND e.bin = $9
              AND e.lot_number = $10
              AND t.timestamp <= inserted_tx.timestamp
            ORDER BY t.timestamp DESC, e.stock_transaction_entry_id DESC
            LIMIT 1
        ), 0) - $3
    FROM inserted_tx
    RETURNING stock_transaction_entry_id
),

inserted_to_entry AS (
    INSERT INTO stock_transaction_entry (
        stock_transaction_id, account, stock_code, location, bin, lot_number, quantity, running_total
    )
    SELECT inserted_tx.stock_transaction_id, $11, $2, $12, $13, $14, $3,
        COALESCE((
            SELECT running_total
            FROM stock_transaction_entry e
            JOIN stock_transaction t ON t.stock_transaction_id = e.stock_transaction_id
            WHERE e.account = $11
              AND e.stock_code = $2
              AND e.location = $12
              AND e.bin = $13
              AND e.lot_number = $14
              AND t.timestamp <= inserted_tx.timestamp
            ORDER BY t.timestamp DESC, e.stock_transaction_entry_id DESC
            LIMIT 1
        ), 0) + $3
    FROM inserted_tx
    RETURNING stock_transaction_entry_id
),

updated_future_from AS (
    UPDATE stock_transaction_entry e
    SET running_total = running_total - $3
    FROM stock_transaction t, inserted_tx
    WHERE e.stock_transaction_id = t.stock_transaction_id
      AND e.account = $7
      AND e.stock_code = $2
      AND e.location = $8
      AND e.bin = $9
      AND e.lot_number = $10
      AND t.timestamp > inserted_tx.timestamp
    RETURNING e.stock_transaction_entry_id
),

updated_future_to AS (
    UPDATE stock_transaction_entry e
    SET running_total = running_total + $3
    FROM stock_transaction t, inserted_tx
    WHERE e.stock_transaction_id = t.stock_transaction_id
      AND e.account = $11
      AND e.stock_code = $2
      AND e.location = $12
      AND e.bin = $13
      AND e.lot_number = $14
      AND t.timestamp > inserted_tx.timestamp
    RETURNING e.stock_transaction_entry_id
)

SELECT
    (SELECT count(*) FROM inserted_from_entry) AS inserted_from_count,
    (SELECT count(*) FROM inserted_to_entry) AS inserted_to_count,
    (SELECT count(*) FROM updated_future_from) AS updated_from_count,
    (SELECT count(*) FROM updated_future_to) AS updated_to_count;
	`

	for _, t := range *transactions {
		accounts, _ := model.StockTransacationTypeMap[t.TransactionType]

		// IMPORTANT: posting from and to the same Account, Location, Bin and
		// LotNumber is not allowed as it is not compatible with how running totals
		// are calculated
		if accounts.From == accounts.To &&
			t.FromLocation == t.ToLocation &&
			t.FromBin == t.ToBin &&
			t.FromLotNumber == t.ToLotNumber {
			return fmt.Errorf("account, location, bin and lot number cannot be the same for 'From' and 'To'")
		}

		if t.Qty.Equal(decimal.Zero) {
			fmt.Printf("Skipping qty 0 transaction of %s\n", t.StockCode)
			continue
		}

		var insertedFromCount, insertedToCount, updatedFromCount, updatedToCount int
		err := exec.QueryRow(ctx, query,
			t.TransactionType,
			t.StockCode,
			t.Qty,
			t.TransactionNote,
			userID,
			t.Timestamp,
			accounts.From,
			t.FromLocation,
			t.FromBin,
			t.FromLotNumber,
			accounts.To,
			t.ToLocation,
			t.ToBin,
			t.ToLotNumber,
		).Scan(&insertedFromCount, &insertedToCount, &updatedFromCount, &updatedToCount)
		if err != nil {
			log.Println(err)
			return fmt.Errorf("failed to create stock transaction: %v", err)
		}

	}

	return nil
}

func (r *StockTransactionRepository) GetStockTransactions(ctx context.Context, exec db.PGExecutor, input *model.GetTransactionsInput) ([]model.StockTransactionEntry, error) {

	// Base query
	query := `
WITH matched_tx_ids AS (
	SELECT DISTINCT ste.stock_transaction_id
	FROM stock_transaction_entry ste
	JOIN stock_transaction st ON st.stock_transaction_id = ste.stock_transaction_id
	WHERE
		($1 = '' OR ste.account = $1)
		AND
		($2 = '' OR ste.stock_code = $2)
		AND
		($3 = '' OR ste.location = $3)
		AND
		($4 = '' OR ste.bin = $4)
		AND
		($5 = '' OR ste.lot_number = $5)
		AND
		($6::timestamp IS NULL OR st.timestamp <= $6::timestamp)
)

SELECT
	ste.stock_transaction_entry_id,
	st.transaction_type,
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
ORDER BY st.timestamp DESC, ste.stock_transaction_entry_id DESC
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
			&st.StockTransactionEntryID,
			&st.TransactionType,
			&st.Account,
			&st.StockCode,
			&st.Location,
			&st.Bin,
			&st.Quantity,
			&st.LotNumber,
			&st.RunningTotal,
			&st.TransactionBy,
			&st.TransactionByUsername,
			&st.Timestamp,
			&st.StockTransactionID,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, st)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil

}
