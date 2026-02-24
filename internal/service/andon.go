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

type AndonService struct {
	db                  *pgxpool.Pool
	swiftConn           *swift.Connection
	andonRepository     *repository.AndonRepository
	commentRepository   *repository.CommentRepository
	galleryRepository   *repository.GalleryRepository
	teamRepository      *repository.TeamRepository
	notificationService *NotificationService
}

func NewAndonService(
	db *pgxpool.Pool,
	swiftConn *swift.Connection,
	andonRepo *repository.AndonRepository,
	commentRepository *repository.CommentRepository,
	galleryRepository *repository.GalleryRepository,
	teamRepository *repository.TeamRepository,
	notificationService *NotificationService,
) *AndonService {
	return &AndonService{
		db:                  db,
		swiftConn:           swiftConn,
		andonRepository:     andonRepo,
		commentRepository:   commentRepository,
		galleryRepository:   galleryRepository,
		teamRepository:      teamRepository,
		notificationService: notificationService,
	}
}

func (s *AndonService) CreateAndon(
	ctx context.Context,
	andon model.NewAndon,
	userID int,
) error {
	var err error

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	galleryId, err := s.galleryRepository.CreateGallery(
		ctx,
		tx,
		userID,
	)
	if err != nil {
		return err
	}

	andon.GalleryID = galleryId

	// Create a dedicated comment thread for this andon (post migration thread model)
	threadID, err := s.commentRepository.CreateCommentThread(ctx, tx)
	if err != nil {
		return err
	}
	andon.CommentThreadID = threadID

	andonID, err := s.andonRepository.CreateAndonEvent(
		ctx,
		tx,
		andon,
		userID,
	)
	if err != nil {
		return err
	}

	if err := s.commentRepository.SetCommentThreadTargetURL(
		ctx,
		tx,
		threadID,
		fmt.Sprintf("/andons/%d", andonID),
	); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	if err := s.notifyAndonCreated(ctx, andonID, userID); err != nil {
		log.Println("error sending andon notifications:", err)
	}

	return nil
}

func (s *AndonService) notifyAndonCreated(ctx context.Context, andonID int, userID int) error {
	andon, err := s.andonRepository.GetAndonByID(ctx, s.db, andonID, userID)
	if err != nil {
		return err
	}
	if andon == nil || andon.AssignedTeam == 0 {
		return nil
	}

	userIDs, err := s.teamRepository.ListTeamUserIDs(ctx, s.db, andon.AssignedTeam)
	if err != nil {
		return err
	}
	if len(userIDs) == 0 {
		return nil
	}

	title := "New Andon"
	summary := strings.TrimSpace(andon.Description)
	if summary == "" {
		parts := make([]string, 0, 2)
		if strings.TrimSpace(andon.IssueName) != "" {
			parts = append(parts, strings.TrimSpace(andon.IssueName))
		}
		if strings.TrimSpace(andon.Location) != "" {
			parts = append(parts, strings.TrimSpace(andon.Location))
		}
		if strings.TrimSpace(andon.Source) != "" {
			parts = append(parts, strings.TrimSpace(andon.Source))
		}
		summary = strings.Join(parts, " · ")
	}
	if summary == "" {
		summary = "New andon raised."
	}

	targetURL := fmt.Sprintf("/andons/%d", andon.AndonID)

	reasonType := mapAndonSeverityToNotificationReason(andon.Severity)

	for _, recipientID := range userIDs {
		if recipientID == userID {
			continue
		}
		notificationID, err := s.notificationService.CreateNotification(ctx, model.NewNotification{
			UserID:      recipientID,
			ActorUserID: &userID,
			Category:    "andon",
			Title:       title,
			Summary:     summary,
			URL:         targetURL,
			Reason:      andon.IssueName,
			ReasonType:  reasonType,
		})
		if err != nil {
			log.Println("error creating andon notification:", err)
		}

		payload := model.PushNotificationPayload{
			Title:          title,
			Body:           summary,
			URL:            targetURL,
			NotificationID: notificationID,
		}
		if notificationID > 0 {
			query := url.Values{}
			query.Set("Redirect", targetURL)
			payload.URL = fmt.Sprintf("/notifications/%d?%s", notificationID, query.Encode())
		}

		if err := s.notificationService.SendPushNotification(ctx, recipientID, payload, ""); err != nil {
			log.Println("error sending andon push notification:", err)
		}
	}

	return nil
}

func mapAndonSeverityToNotificationReason(severity model.AndonSeverity) string {
	switch severity {
	case model.AndonSeverityRequiresIntervention:
		return model.NotificationReasonDanger
	case model.AndonSeveritySelfResolvable:
		return model.NotificationReasonWarning
	default:
		return model.NotificationReasonInfo
	}
}

func (s *AndonService) GetAndonByID(
	ctx context.Context,
	andonEventID int,
	userID int,
) (*model.Andon, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	andonEvent, err := s.andonRepository.GetAndonByID(
		ctx,
		s.db,
		andonEventID,
		userID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return andonEvent, nil
}

func (s *AndonService) UpdateAndon(
	ctx context.Context,
	andonEventID int,
	action string,
	userID int,
) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	switch action {
	case "acknowledge":
		err = s.andonRepository.AcknowledgeAndon(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "resolve":
		err = s.andonRepository.ResolveAndon(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "cancel":
		err = s.andonRepository.CancelAndon(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	case "reopen":
		err = s.andonRepository.ReopenAndon(
			ctx,
			tx,
			andonEventID,
			userID,
		)
	}

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *AndonService) ListAndons(
	ctx context.Context,
	q model.ListAndonQuery,
	userID int,
) ([]model.Andon, int, model.AndonAvailableFilters, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}
	defer tx.Rollback(ctx)

	var andons []model.Andon

	andons, err = s.andonRepository.ListAndons(
		ctx,
		tx,
		q,
		userID,
	)
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	count, err := s.andonRepository.Count(ctx, tx, q)
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	availableFilters, err := s.andonRepository.GetAvailableFilters(ctx, tx, model.AndonFilters{
		StartDate:                q.StartDate,
		EndDate:                  q.EndDate,
		LocationIn:               q.LocationIn,
		IssueIn:                  q.IssueIn,
		TeamIn:                   q.TeamIn,
		SeverityIn:               q.SeverityIn,
		StatusIn:                 q.StatusIn,
		RaisedByUsernameIn:       q.RaisedByUsernameIn,
		AcknowledgedByUsernameIn: q.AcknowledgedByUsernameIn,
		ResolvedByUsernameIn:     q.ResolvedByUsernameIn,
	})
	if err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return []model.Andon{}, 0, model.AndonAvailableFilters{}, err
	}

	return andons, count, availableFilters, nil
}

func (s *AndonService) GetAndonChangelog(
	ctx context.Context,
	andonEventID int,
	userID int,
) ([]model.AndonChange, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	changelog, err := s.andonRepository.GetAndonChangelog(
		ctx,
		tx,
		andonEventID,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return changelog, nil
}

func (s *AndonService) CreateAndonComment(
	ctx context.Context,
	comment *model.NewComment,
	userID int,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = s.commentRepository.AddComment(
		ctx,
		tx,
		comment,
		userID,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
