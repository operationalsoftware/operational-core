package service

import (
	"app/internal/model"
	"app/internal/repository"
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ncw/swift/v2"
)

type CommentService struct {
	db                  *pgxpool.Pool
	swiftConn           *swift.Connection
	commentRepository   *repository.CommentRepository
	userRepository      *repository.UserRepository
	notificationService *NotificationService
}

func NewCommentService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	commentRepository *repository.CommentRepository,
	userRepository *repository.UserRepository,
	notificationService *NotificationService,
) *CommentService {
	return &CommentService{
		db:                  db,
		swiftConn:           swiftConn,
		commentRepository:   commentRepository,
		userRepository:      userRepository,
		notificationService: notificationService,
	}
}

func (s *CommentService) GetComments(
	ctx context.Context,
	commentThreadID int,
	userID int,
) ([]model.Comment, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.Comment{}, err
	}
	defer tx.Rollback(ctx)

	comments, err := s.commentRepository.GetComments(
		ctx,
		tx,
		s.swiftConn,
		commentThreadID,
	)

	if err != nil {
		return []model.Comment{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return []model.Comment{}, err
	}

	return comments, nil
}

func (s *CommentService) CreateComment(
	ctx context.Context,
	comment *model.NewComment,
	userID int,
) (int, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	commentID, err := s.commentRepository.AddComment(
		ctx,
		tx,
		comment,
		userID,
	)
	if err != nil {
		return 0, err
	}

	// Parse @mentions now so mention notifications can be added without
	// changing this write path later.
	mentionedUsernames := extractMentionUsernames(comment.Comment)

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	if err := s.sendMentionNotifications(
		ctx,
		mentionedUsernames,
		userID,
		commentID,
		comment.CommentThreadID,
		comment.Comment,
	); err != nil {
		log.Println("error sending mention notifications:", err)
	}

	return commentID, nil
}

func (s *CommentService) sendMentionNotifications(
	ctx context.Context,
	mentionedUsernames []string,
	actorUserID int,
	commentID int,
	commentThreadID int,
	commentText string,
) error {
	if len(mentionedUsernames) == 0 || s.notificationService == nil {
		return nil
	}

	threadURL, found, err := s.commentRepository.ResolveCommentThreadURL(ctx, s.db, commentThreadID)
	if err != nil {
		return fmt.Errorf("error resolving comment thread url: %w", err)
	}
	if !found {
		return fmt.Errorf("comment thread url not found for thread %d", commentThreadID)
	}

	recipientsByUsername, err := s.userRepository.ResolveUserIDsByUsernames(ctx, s.db, mentionedUsernames)
	if err != nil {
		return fmt.Errorf("error resolving mentioned users: %w", err)
	}
	if len(recipientsByUsername) == 0 {
		return nil
	}

	targetURL := fmt.Sprintf("%s#comment-%d", strings.TrimSpace(threadURL), commentID)
	notificationTitle := "Mention in comment"
	notificationSummary := mentionNotificationSummary(commentText)
	actor := actorUserID

	var firstErr error
	for _, username := range mentionedUsernames {
		recipientID, ok := recipientsByUsername[username]
		if !ok || recipientID == actorUserID {
			continue
		}

		notificationID, err := s.notificationService.CreateNotification(ctx, model.NewNotification{
			UserID:      recipientID,
			ActorUserID: &actor,
			Category:    "comment",
			Title:       notificationTitle,
			Summary:     notificationSummary,
			URL:         targetURL,
			Reason:      "Mention",
			ReasonType:  model.NotificationReasonInfo,
		})
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			log.Println("error creating mention notification:", err)
			continue
		}

		payloadURL := targetURL
		if notificationID > 0 {
			query := url.Values{}
			query.Set("Redirect", targetURL)
			payloadURL = fmt.Sprintf("/notifications/%d?%s", notificationID, query.Encode())
		}

		payload := model.PushNotificationPayload{
			Title:          notificationTitle,
			Body:           notificationSummary,
			URL:            payloadURL,
			NotificationID: notificationID,
		}
		if err := s.notificationService.SendPushNotification(ctx, recipientID, payload, ""); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			log.Println("error sending mention push notification:", err)
		}
	}

	return firstErr
}

func extractMentionUsernames(commentText string) []string {
	if strings.TrimSpace(commentText) == "" {
		return nil
	}

	mentions := []string{}
	seen := map[string]struct{}{}
	runes := []rune(commentText)

	for i, r := range runes {
		if r != '@' {
			continue
		}

		// Ignore emails and words containing '@' in the middle.
		if i > 0 && isMentionUsernameChar(runes[i-1]) {
			continue
		}

		start := i + 1
		end := start
		for end < len(runes) && isMentionUsernameChar(runes[end]) {
			end++
		}
		if end == start {
			continue
		}

		username := strings.ToLower(string(runes[start:end]))
		if len(username) < 3 || len(username) > 20 {
			continue
		}

		if _, exists := seen[username]; exists {
			continue
		}
		seen[username] = struct{}{}
		mentions = append(mentions, username)
	}

	return mentions
}

func isMentionUsernameChar(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '_'
}

func mentionNotificationSummary(commentText string) string {
	summary := strings.TrimSpace(commentText)
	summary = strings.Join(strings.Fields(summary), " ")
	if summary == "" {
		return "You were mentioned in a comment."
	}

	const maxLen = 160
	runes := []rune(summary)
	if len(runes) <= maxLen {
		return summary
	}

	return string(runes[:maxLen-3]) + "..."
}
