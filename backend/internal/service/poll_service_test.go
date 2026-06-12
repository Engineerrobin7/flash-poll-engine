package service

import (
	"context"
	"flashpoll/internal/domain"
	"testing"
)

func TestPollValidation(t *testing.T) {
	svc := &PollService{} // Mock or partial init for logic testing

	t.Run("Reject short question", func(t *testing.T) {
		_, err := svc.CreatePoll(context.Background(), "No", domain.CategoryTech, []string{"A", "B"})
		if err == nil {
			t.Errorf("Expected error for short question, got nil")
		}
	})

	t.Run("Reject duplicate options", func(t *testing.T) {
		_, err := svc.CreatePoll(context.Background(), "Valid Question", domain.CategoryTech, []string{"Same", "same"})
		if err == nil {
			t.Errorf("Expected error for duplicate options (case-insensitive), got nil")
		}
	})
}
