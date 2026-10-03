package dto

import "time"

type Format struct {
	Teams int `json:"teams"`
	Size  int `json:"size"`
}

type ParticipantResponse struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Online   bool      `json:"online"`
	JoinedAt time.Time `json:"joinedAt"`
}

type SpinResponse struct {
	WinnerID     string    `json:"winnerId"`
	PoolSnapshot []string  `json:"poolSnapshot"`
	TargetOffset float64   `json:"targetOffset"`
	Rotations    int       `json:"rotations"`
	StartedAt    time.Time `json:"startedAt"`
	DurationMs   int       `json:"durationMs"`
}

type RoomResponse struct {
	Participants []ParticipantResponse   `json:"participants"`
	AdminID      string                  `json:"adminId"`
	Spinners     map[string]bool         `json:"spinners"`
	Ready        map[string]bool         `json:"ready"`
	AutoSpin     bool                    `json:"autoSpin"`
	Format       Format                  `json:"format"`
	Pool         []string                `json:"pool"`
	Teams        [][]ParticipantResponse `json:"teams"`
	Phase        string                  `json:"phase"`
	Spin         *SpinResponse           `json:"spin"`
	Picks        int                     `json:"picks"`
}
