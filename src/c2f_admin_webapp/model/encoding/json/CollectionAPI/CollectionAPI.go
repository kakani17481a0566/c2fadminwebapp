package CollectionAPI

type Root struct {
	ServiceName     string           `json:"ServiceName,omitempty"`
	ServiceType     string           `json:"ServiceType,omitempty"`
	Version         string           `json:"Version,omitempty"`
	Channel         *Channel         `json:"Channel,omitempty"`
	Device          *Device          `json:"Device,omitempty"`
	Agent           *Users           `json:"Agent,omitempty"`
	TransactionData *TransactionData `json:"TransactionData,omitempty"`
	Security        *Security        `json:"Security,omitempty"`
	Response        *Response        `json:"Response,omitempty"`
	Parameters      []*Parameter     `json:"Parameters,omitempty"`
	Row             []*Row           `json:"Row,omitempty"`
}

type TransactionData struct {
	CountryCode            string           `json:"CountryCode,omitempty"`
	Currency               string           `json:"Currency,omitempty"`
	CustomerDetails        *CustomerDetails `json:"CustomerDetails,omitempty"`
	PaymentMode            string           `json:"PaymentMode,omitempty"`
	Qty                    string           `json:"Qty,omitempty"`
	RequestID              string           `json:"RequestID,omitempty"`
	TimeStamp              string           `json:"TimeStamp,omitempty"`
	TransactionDescription string           `json:"TransactionDescription,omitempty"`
	TransactionNumber      string           `json:"TransactionNumber,omitempty"`
	Amount                 string           `json:"Amount,omitempty"`
}

type CustomerDetails struct {
	ServiceNo string `json:"ServiceNo,omitempty"`
	Name      string `json:"Name,omitempty"`
}

type Security struct {
	Data string `json:"Data,omitempty"`
	Hmac string `json:"Hmac,omitempty"`
	Skey string `json:"Skey,omitempty"`
}

type Response struct {
	Code        string       `json:"Code,omitempty"`
	Message     string       `json:"Message,omitempty"`
	Status      string       `json:"Status,omitempty"`
	HostCode    string       `json:"HostCode,omitempty"`
	HostMessage string       `json:"HostMessage,omitempty"`
	Parameter   []*Parameter `json:"Parameter,omitempty"`
}

type Parameter struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Device struct {
	IP      IP     `json:"IP,omitempty"`
	ID      string `json:"Id,omitempty"`
	OS      string `json:"OS,omitempty"`
	Version string `json:"Version,omitempty"`
}

type IP struct {
	Type  string `json:"Type,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Channel struct {
	HostIP HostIP `json:"HostIP,omitempty"`
	Name   string `json:"Name,omitempty"`
	Type   string `json:"Type,omitempty"`
}

type HostIP struct {
	Type  string `json:"Type,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Users struct {
	Id            string
	Address       []*Address       `json:"Address,omitempty"`
	Communication []*Communication `json:"Communication,omitempty"`
	CreateDate    string           `json:"CreateDate,omitempty"`
	LastUpdate    string           `json:"LastUpdate,omitempty"`
	ParentID      string           `json:"ParentId,omitempty"`
	Profile       *Profile         `json:"Profile,omitempty"`
	Status        string           `json:"Status,omitempty"`
	Type          string           `json:"Type,omitempty"`
	Username      string           `json:"Username,omitempty"`
}

type Profile struct {
	Id         string
	CreateDate string `json:"CreateDate,omitempty"`
	DOB        string `json:"DOB,omitempty"`
	FirstName  string `json:"FirstName,omitempty"`
	Gender     string `json:"Gender,omitempty"`
	LastName   string `json:"LastName,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	MiddleName string `json:"MiddleName,omitempty"`
	NationalID string `json:"NationalID,omitempty"`
}

type Address struct {
	Id         string
	Address1   string `json:"Address1,omitempty"`
	Address2   string `json:"Address2,omitempty"`
	City       string `json:"City,omitempty"`
	Country    string `json:"Country,omitempty"`
	CreateDate string `json:"CreateDate,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	Province   string `json:"Province,omitempty"`
	Type       string `json:"Type,omitempty"`
	Zip        string `json:"Zip,omitempty"`
}

type Communication struct {
	Id         string
	CreateDate string `json:"CreateDate,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	Type       string `json:"Type,omitempty"`
	Value      string `json:"Value,omitempty"`
}

type Column struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Row struct {
	Columns []*Column `json:"Columns,omitempty"`
}
