package nget

import (
	"fmt"
)

// A Request for a score: scoreType = "level" or "episode"
type NScoresRequest struct {
	scoreType string
	scoreId   int
	attempts  int
}

func NewNScoresRequest(scoreType string, scoreId int) *NScoresRequest {
	return &NScoresRequest{
		scoreType: scoreType,
		scoreId:   scoreId,
	}
}

// A sortable list of score responses
type NScoresResponseList []NScoresResponse

func (l NScoresResponseList) Len() int { return len(l) }
func (l NScoresResponseList) Less(i, j int) bool {
	if l[i].EpisodeId != nil {
		if l[j].EpisodeId == nil {
			return true
		}
		return *l[i].EpisodeId < *l[j].EpisodeId
	}
	if l[j].EpisodeId != nil {
		return false
	}
	return *l[i].LevelId < *l[j].LevelId
}
func (l NScoresResponseList) Swap(i, j int){ l[i], l[j] = l[j], l[i] }

type NScoresResponse struct {
	EpisodeId *int            `json:"episode_id,omitempty"`
	LevelId   *int            `json:"level_id,omitempty"`
	Scores    []NScore        `json:"scores"`
	queryType int             `json:"query_type,omitempty"`
	Err       error           `json:"-"`
	Request   *NScoresRequest `json:"-"`
}
func (n NScoresResponse) IsEpisodeScore() bool {
	return n.EpisodeId != nil
}

type NScore struct {
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
	UserName string `json:"user_name"`
	UserId   int    `json:"user_id,omitempty"`
	replayId int    `json:"replay_id,omitempty"`
}

type NUserInfo struct {
	MyDisplayName string `json:"my_display_name"`
	MyRank        int    `json:"my_rank"`
	MyScore       int    `json:"my_score"`
	MyReplayId    int    `json:"my_replay_id"`
}

// Errors
type Error struct {
	Code int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v - %v", e.Code, e.Message)
}

// Possible error responses
var (
	ErrRead        = &Error{1001, "Bad Request"}
	ErrNoResponse  = &Error{2001, "No response from Metanet"}
	ErrBadResponse = &Error{2002, "Could not decode Metanet response"}
	ErrInternal    = &Error{9999, "Internal error"}
)
