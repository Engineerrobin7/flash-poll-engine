package domain

import "time"

type Category string

const (
	CategoryTech     Category = "Tech"
	CategoryBusiness Category = "Business"
	CategoryDesign   Category = "Design"
)

type Poll struct {
	ID         string    `json:"id"`
	Category   Category  `json:"category"`
	Question   string    `json:"question"`
	CreatedAt  time.Time `json:"created_at"`
	Options    []Option  `json:"options"`
	TotalVotes int       `json:"total_votes"`
}

type Option struct {
	ID         int64   `json:"id"`
	PollID     string  `json:"poll_id"`
	OptionText string  `json:"option_text"`
	VoteCount  int     `json:"vote_count"`
	Percentage float64 `json:"percentage"`
}
