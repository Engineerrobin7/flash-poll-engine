package service

import (
	"context"
	"database/sql"
	"errors"
	"flashpoll/internal/domain"
	"flashpoll/internal/repository"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

var (
	ErrPollNotFound   = errors.New("poll not found")
	ErrOptionNotFound = errors.New("option not found")
	ErrInvalidInput   = errors.New("invalid input")
)

type PollService struct {
	repo   *repository.PollRepository
	Broker *Broker
}

func NewPollService(repo *repository.PollRepository, broker *Broker) *PollService {
	return &PollService{repo: repo, Broker: broker}
}

func (s *PollService) GetPolls(ctx context.Context) ([]domain.Poll, error) {
	return s.repo.GetAll(ctx)
}

func (s *PollService) CreatePoll(ctx context.Context, question string, category domain.Category, optionTexts []string) (*domain.Poll, error) {
	// Validation
	question = strings.TrimSpace(question)
	if len(question) < 5 || len(question) > 280 {
		return nil, ErrInvalidInput
	}

	if category != domain.CategoryTech && category != domain.CategoryBusiness && category != domain.CategoryDesign {
		return nil, ErrInvalidInput
	}

	if len(optionTexts) < 2 || len(optionTexts) > 10 {
		return nil, ErrInvalidInput
	}

	normalizedOptions := make(map[string]bool)
	var options []domain.Option
	for _, text := range optionTexts {
		trimmed := strings.TrimSpace(text)
		// Collapse extra internal spaces
		collapsed := strings.Join(strings.Fields(trimmed), " ")
		if collapsed == "" || len(collapsed) > 80 {
			return nil, ErrInvalidInput
		}
		lower := strings.ToLower(collapsed)
		if normalizedOptions[lower] {
			return nil, ErrInvalidInput
		}
		normalizedOptions[lower] = true
		options = append(options, domain.Option{OptionText: collapsed})
	}

	poll := &domain.Poll{
		ID:        uuid.New().String(),
		Question:  question,
		Category:  category,
		CreatedAt: time.Now().UTC(),
		Options:   options,
	}

	if err := s.repo.Create(ctx, poll); err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *PollService) Vote(ctx context.Context, pollID string, optionID int64) (*domain.Poll, error) {
	log.Printf("Processing vote for Poll: %s, Option: %d", pollID, optionID)
	err := s.repo.IncrementVote(ctx, pollID, optionID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Vote failed: poll or option not found")
			return nil, ErrPollNotFound
		}
		log.Printf("Vote failed in repository: %v", err)
		return nil, err
	}

	poll, err := s.repo.GetByID(ctx, pollID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPollNotFound
		}
		log.Printf("Failed to fetch updated poll after vote: %v", err)
		return nil, err
	}

	// Broadcast the update to all connected clients
	s.Broker.BroadcastUpdate(poll)

	return poll, nil
}

func (s *PollService) DeletePoll(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
