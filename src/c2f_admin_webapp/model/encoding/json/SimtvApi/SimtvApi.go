package SimtvApi

type SimtvResp struct {
	StatusCode int    `json:"status_code,omitempty"`
	TxnStatus  string `json:"status_txt,omitempty"`
	ErrMessage string `json:"error_msg,omitempty"`
	Data       *Data  `json:"data,omitempty"`
}
type Data struct {
	CasID       string  `json:"cas_id,omitempty"`
	CustomerID  string  `json:"customer_id,omitempty"`
	Customer    string  `json:"customer,omitempty"`
	Status      string  `json:"status,omitempty"`
	Message     string  `json:"msg,omitempty"`
	ReferenceID string  `json:"reference_id,omitempty"`
	TxnID       string  `json:"transaction_id,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Date        string  `json:"date,omitempty"`
	OperatorID  string  `json:"operator_id,omitempty"`
	Operator    string  `json:"operator,omitempty"`
	Tpe         string  `json:"type,omitempty"`
}
