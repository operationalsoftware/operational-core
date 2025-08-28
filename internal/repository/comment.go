package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
	"encoding/json"
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

func (r *CommentRepository) AddComment(
	ctx context.Context,
	exec db.PGExecutor,
	comment *model.NewComment,
	userID int,
) (int, error) {

	query := `
INSERT INTO comment (
	entity,
	entity_id,
	comment,
	commented_by
)
VALUES (
	$1,
	$2,
	$3,
	$4
)
RETURNING comment_id
`
	var commentID int
	err := exec.QueryRow(
		ctx, query,

		comment.Entity,
		comment.EntityID,
		comment.Comment,
		userID,
	).Scan(&commentID)

	return commentID, err
}

func (r *CommentRepository) GetComments(
	ctx context.Context,
	exec db.PGExecutor,
	swiftConn *swift.Connection,
	entity string,
	entityID int,
) ([]model.Comment, error) {

	query := `
SELECT
	comment_id,
	entity,
	entity_id,
	comment,
	commented_by_username,
	commented_at,
	attachments
FROM comment_view
WHERE entity = $1 AND entity_id = $2
`

	rows, err := exec.Query(ctx, query, entity, entityID)
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
			&comment.Entity,
			&comment.EntityID,
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
