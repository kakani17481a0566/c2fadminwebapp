package ReportApi

type Root struct {
	Version           string       `json:"Version,omitempty"`
	ReportName        string       `json:"ReportName,omitempty"`
	User              string       `json:"User,omitempty"`
	RequestID         string       `json:"RequestID,omitempty"`
	TransactionNumber string       `json:"TransactionNumber,omitempty"`
	TimeStamp         string       `json:"RequestID,omitempty"`
	Parameters        []*Parameter `json:"Parameters,omitempty"`
	Response          *Response    `json:"Response,omitempty"`
	Row               []*Row       `json:"Row,omitempty"`
}

type Response struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

type Parameter struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Column struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Row struct {
	Columns []*Column `json:"Columns,omitempty"`
}
