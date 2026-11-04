package ZETDC

type Root struct {
	Code         string          `json:"Code,omitempty"`
	Message      string          `json:"Message,omitempty"`
	Token        string          `json:"Token,omitempty"`
	Amount       string          `json:"Amount,omitempty"`
	Status       string          `json:"Status,omitempty"`
	Meter        string          `json:"Meter,omitempty"`
	Name         string          `json:"Name,omitempty"`
	Address      string          `json:"Address,omitempty"`
	Mobile       string          `json:"Mobile,omitempty"`
	TxnNumber    string          `json:"TxnNumber,omitempty"`
	ReadingDone  string          `json:"ReadingDone,omitempty"`
	UsernameType string          `json:"UsernameType,omitempty"`
	FriendMobile string          `json:"FriendMobile,omitempty"`
	FriendMeter  string          `json:"FriendMeter,omitempty"`
	Emergency    string          `json:"Emergency,omitempty"`
	OTP          string          `json:"OTP,omitempty"`
	Energy       string          `json:"Energy,omitempty"`
	TS           string          `json:"TS,omitempty"`
	Usage        []*Usage        `json:"Usage,omitempty"`
	History      []*History      `json:"History,omitempty"`
	Metertype    string          `json:"MeterType,omitempty"`
	MessageType  string          `json:"MessageType,omitempty"`
	Notification []*Notification `json:"Notification,omitempty"`
	Password     string          `json:"Password,omitempty"`
	Reading      string          `json:"MeterReading,omitempty"`
	BillDate     string          `json:"BillDate,omitempty"`
	BillAmount   string          `json:"BillAmount,omitempty"`
	DueDate      string          `json:"DueDate,omitempty"`
	DueAmount    string          `json:"DueAmount,omitempty"`
	Comment      string          `json:"Comment,omitempty"`
}

type Notification struct {
	Message string `json:"Message,omitempty"`
}

type Usage struct {
	Month string `json:"Month,omitempty"`
	Value string `json:"Value,omitempty"`
}

type History struct {
	Type          string `json:"Type,omitempty"`
	Token         string `json:"Token,omitempty"`
	DateTime      string `json:"DateTime,omitempty"`
	TokenAmount   string `json:"TokenAmount,omitempty"`
	Energy        string `json:"Energy,omitempty"`
	BillDate      string `json:"BillDate,omitempty"`
	BillAmount    string `json:"BillAmount,omitempty"`
	DueDate       string `json:"DueDate,omitempty"`
	DueAmount     string `json:"DueAmount,omitempty"`
	PaymentDate   string `json:"PaymentDate,omitempty"`
	PaymentAmount string `json:"PaymentAmount,omitempty"`
}
