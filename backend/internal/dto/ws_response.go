package dto

type StateMessage struct {
	Type    string       `json:"type"`
	Room    RoomResponse `json:"room"`
	You     string       `json:"you"`
	Formats []Format     `json:"formats"`
}

type ErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
