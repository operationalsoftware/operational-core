package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ncw/swift/v2"
)

type CommentRepository struct {
	fileRepo *FileRepository
}

func NewCommentRepository(fileRepo *FileRepository) *CommentRepository {
	return &CommentRepository{
		fileRepo: fileRepo,
	}
}

// CreateCommentThread creates an empty comment thread and returns its id.
// After migration 00000200 comment threads are decoupled from entities, so we
// just insert default values.
func (r *CommentRepository) CreateCommentThread(
	ctx context.Context,
	exec db.PGExecutor,
) (int, error) {
	var id int
	err := exec.QueryRow(ctx, `INSERT INTO comment_thread DEFAULT VALUES RETURNING comment_thread_id`).Scan(&id)
	return id, err
}

func (r *CommentRepository) SetCommentThreadTargetURL(
	ctx context.Context,
	exec db.PGExecutor,
	commentThreadID int,
	targetURL string,
) error {
	query := `
UPDATE comment_thread
SET
	target_url = $2
WHERE
	comment_thread_id = $1
`
	result, err := exec.Exec(ctx, query, commentThreadID, targetURL)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("comment thread %d not found", commentThreadID)
	}
	return nil
}

func (r *CommentRepository) DeleteCommentThread(
	ctx context.Context,
	exec db.PGExecutor,
	commentThreadID int,
) error {
	query := `
DELETE FROM
	comment_thread
WHERE
	comment_thread_id = $1
`

	_, err := exec.Exec(ctx, query, commentThreadID)
	if err != nil {
		return err
	}

	return nil
}

func (r *CommentRepository) AddComment(
	ctx context.Context,
	exec db.PGExecutor,
	comment *model.NewComment,
	userID int,
) (int, error) {

	query := `
INSERT INTO comment (
	comment_thread_id,
	comment,
	commented_by
)
VALUES (
	$1,
	$2,
	$3
)
RETURNING comment_id
`
	var commentID int
	err := exec.QueryRow(
		ctx, query,

		comment.CommentThreadID,
		comment.Comment,
		userID,
	).Scan(&commentID)

	return commentID, err
}

func (r *CommentRepository) GetComments(
	ctx context.Context,
	exec db.PGExecutor,
	swiftConn *swift.Connection,
	commentThreadID int,
) ([]model.Comment, error) {

	query := `
SELECT
	comment_id,
	comment_thread_id,
	comment,
	commented_by_username,
	commented_at,
	attachments
FROM comment_view
WHERE comment_thread_id = $1
`

	rows, err := exec.Query(ctx, query, commentThreadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		var attachments []byte
		if err := rows.Scan(
			&comment.CommentID,
			&comment.CommentThreadID,
			&comment.Comment,
			&comment.CommentedByUsername,
			&comment.CommentedAt,
			&attachments,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(attachments, &comment.Attachments); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for ci := range comments {
		for fi := range comments[ci].Attachments {
			url, err := r.fileRepo.GetSignedDownloadURL(ctx, swiftConn, exec, comments[ci].Attachments[fi].FileID, 15*time.Minute)

			if err != nil {
				return nil, err
			}
			comments[ci].Attachments[fi].DownloadURL = url
		}
	}

	return comments, nil
}

func (r *CommentRepository) ResolveCommentThreadURL(
	ctx context.Context,
	exec db.PGExecutor,
	commentThreadID int,
) (string, bool, error) {
	query := `
SELECT
	target_url
FROM
	comment_thread
WHERE
	comment_thread_id = $1
	`

	var url *string
	if err := exec.QueryRow(ctx, query, commentThreadID).Scan(&url); err != nil {
		return "", false, err
	}
	if url == nil || *url == "" {
		return "", false, nil
	}

	return *url, true, nil
}
