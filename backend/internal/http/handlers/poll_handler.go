package handlers

import (
	"encoding/json"
	"flashpoll/internal/domain"
	"flashpoll/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type PollHandler struct {
	service *service.PollService
}

func NewPollHandler(service *service.PollService) *PollHandler {
	return &PollHandler{service: service}
}

type Response struct {
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

func (h *PollHandler) GetPolls(w http.ResponseWriter, r *http.Request) {
	polls, err := h.service.GetPolls(r.Context())
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Unable to retrieve polls")
		return
	}
	h.sendJSON(w, http.StatusOK, Response{Data: polls})
}

type CreatePollRequest struct {
	Question string          `json:"question"`
	Category domain.Category `json:"category"`
	Options  []string        `json:"options"`
}

func (h *PollHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	var req CreatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload")
		return
	}

	poll, err := h.service.CreatePoll(r.Context(), req.Question, req.Category, req.Options)
	if err != nil {
		if err == service.ErrInvalidInput {
			h.sendError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid poll payload")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Unable to create poll")
		return
	}

	h.sendJSON(w, http.StatusCreated, Response{Data: poll})
}

type VoteRequest struct {
	OptionID int64 `json:"option_id"`
}

func (h *PollHandler) Vote(w http.ResponseWriter, r *http.Request) {
	pollID := chi.URLParam(r, "id")
	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON payload")
		return
	}

	if req.OptionID <= 0 {
		h.sendError(w, http.StatusBadRequest, "VALIDATION_ERROR", "option_id is required and must be a positive integer")
		return
	}

	poll, err := h.service.Vote(r.Context(), pollID, req.OptionID)
	if err != nil {
		if err == service.ErrPollNotFound || err == service.ErrOptionNotFound {
			h.sendError(w, http.StatusNotFound, "POLL_OR_OPTION_NOT_FOUND", "The specified poll or option does not exist")
			return
		}
		h.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Unable to record vote")
		return
	}

	h.sendJSON(w, http.StatusOK, Response{Data: poll})
}

func (h *PollHandler) DeletePoll(w http.ResponseWriter, r *http.Request) {
	pollID := chi.URLParam(r, "id")
	err := h.service.DeletePoll(r.Context(), pollID)
	if err != nil {
		h.sendError(w, http.StatusNotFound, "POLL_NOT_FOUND", "The requested poll does not exist")
		return
	}

	h.sendJSON(w, http.StatusOK, Response{Data: map[string]interface{}{"id": pollID, "deleted": true}})
}

func (h *PollHandler) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *PollHandler) sendError(w http.ResponseWriter, status int, code, message string) {
	log.Printf("API Error [%d] %s: %s", status, code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	json.NewEncoder(w).Encode(resp)
}

func (h *PollHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	polls, err := h.service.GetPolls(r.Context())
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Unable to compute stats")
		return
	}

	totalVotes := 0
	for _, p := range polls {
		totalVotes += p.TotalVotes
	}

	h.sendJSON(w, http.StatusOK, Response{Data: map[string]int{
		"total_polls": len(polls),
		"total_votes": totalVotes,
	}})
}
