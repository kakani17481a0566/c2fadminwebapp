package ESolutions

type Root struct {
	Originator       string `json:"originator,omitempty"`
	Destination      string `json:"destination,omitempty"`
	MessageText      string `json:"messageText,omitempty"`
	MessageReference string `json:"messageReference,omitempty"`
	MessageDate      string `json:"messageDate,omitempty"`
	MessageId        string `json:"messageId,omitempty"`
	Status           string `json:"status,omitempty"`
	Charge           int    `json:"charge,omitempty"`
	Scheduled        bool   `json:"scheduled,omitempty"`
}
