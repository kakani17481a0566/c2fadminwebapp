package PID

type Root struct {
	DeviceValidationOTP    string `json:"DeviceValidationOTP,omitempty"`
	LoginPassword          string `json:"LoginPassword,omitempty"`
	NewLoginPassword       string `json:"NewLoginPassword,omitempty"`
	RegistrtionOTP         string `json:"RegistrtionOTP,omitempty"`
	TransactionOTP         string `json:"TransactionOTP,omitempty"`
	TransactionPassowrd    string `json:"TransactionPassowrd,omitempty"`
	NewTransactionPassowrd string `json:"NewTransactionPassowrd,omitempty"`
	BrowserID              string `json:"BrowserID,omitempty"`
	SessionID              string `json:"SessionID,omitempty"`

	Voucher *[]Voucher `json:"Voucher,omitempty"`
}

type Voucher struct {
	Code         string `json:"Code,omitempty"`
	Denomination string `json:"Denomination,omitempty"`
	ExpiryDate   string `json:"ExpiryDate,omitempty"`
	Partner      string `json:"Partner,omitempty"`
	SerialNumber string `json:"SerialNumber,omitempty"`
}
