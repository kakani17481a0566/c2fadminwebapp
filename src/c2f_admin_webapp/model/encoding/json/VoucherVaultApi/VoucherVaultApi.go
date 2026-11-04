package VoucherVaultApi

type Root struct {
	Source             *Source             `json:"Source,omitempty"`
	ServiceDetails     *ServiceDetails     `json:"ServiceDetails,omitempty"`
	TransactionDetails *TransactionDetails `json:"TransactionDetails,omitempty"`
	Response           *Response           `json:"Response,omitempty"`
}

type Source struct {
	Name string `json:"Name,omitempty"`
}

type ServiceDetails struct {
	Name         string `json:"Name,omitempty"`
	Operator     string `json:"Operator,omitempty"`
	Denomination string `json:"Denomination,omitempty"`
	Qty          string `json:"Qty,omitempty"`
	FileName     string `json:"FileName,omitempty"`
	FilePath     string `json:"FilePath,omitempty"`
	FTPSite      string `json:"FTPSite,omitempty"`
	FTPUsername  string `json:"FTPUsername,omitempty"`
	FTPPassword  string `json:"FTPPassword,omitempty"`
}

type Response struct {
	Code      string       `json:"Code,omitempty"`
	Message   string       `json:"Message,omitempty"`
	Status    string       `json:"Status,omitempty"`
	Parameter []*Parameter `json:"Parameter,omitempty"`
}

type Parameter struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

type TransactionDetails struct {
	RequestID         string `json:"RequestID,omitempty"`
	TransactionNumber string `json:"TransactionNumber,omitempty"`
	RequestTime       string `json:"RequestTime,omitempty"`
	ResponseTime      string `json:"ResponseTime,omitempty"`
}
