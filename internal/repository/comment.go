package repository

import (
	"app/internal/model"
	"app/pkg/db"
	"context"
)

type CommentRepository struct{}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{}
}

func (r *CommentRepository) AddComment(
	ctx context.Context,
	exec db.PGExecutor,
	comment *model.NewComment,
	userID int,
) error {

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
`
	_, err := exec.Exec(
		ctx, query,

		comment.Entity,
		comment.EntityID,
		comment.Comment,
		userID,
	)

	return err
}

func (r *CommentRepository) GetComments(
	ctx context.Context,
	exec db.PGExecutor,
	entity string,
	entityID int,
) ([]model.Comment, error) {

	query := `
SELECT
	c.comment_id,
	c.entity_id,
	c.comment,
	u.username,
	c.commented_at
FROM comment c
LEFT JOIN
	app_user AS u ON c.commented_by = u.user_id
WHERE c.entity = $1 AND c.entity_id = $2
ORDER BY c.commented_at ASC
`

	rows, err := exec.Query(ctx, query, entity, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(
			&comment.CommentID,
			&comment.EntityID,
			&comment.Comment,
			&comment.CommentedByUsername,
			&comment.CommentedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
