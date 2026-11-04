package SwitchApi

type Root struct {
	Payer              *Payer              `json:"Payer,omitempty"`
	Payee              *Payee              `json:"Payee,omitempty"`
	ServiceDetails     *ServiceDetails     `json:"ServiceDetails,omitempty"`
	CustomerDetails    *CustomerDetails    `json:"CustomerDetails,omitempty"`
	TransactionDetails *TransactionDetails `json:"TransactionDetails,omitempty"`
	Response           *Response           `json:"Response,omitempty"`
}

type Payer struct {
	Name string `json:"Name,omitempty"`
}

type Payee struct {
	Name string `json:"Name,omitempty"`
}

type ServiceDetails struct {
	BusinessCategory      string `json:"BusinessCategory,omitempty"`
	Category              string `json:"Category,omitempty"`
	Denomination          string `json:"Denomination,omitempty"`
	Partner               string `json:"Partner,omitempty"`
	SubCategory           string `json:"SubCategory,omitempty"`
	UserType              string `json:"UserType,omitempty"`
	Name                  string `json:"Name,omitempty"`
	Operator              string `json:"Operator,omitempty"`
	MODE                  string `json:"MODE,omitempty"`
	Number                string `json:"Number,omitempty"`
	NumberType            string `json:"NumberType,omitempty"`
	BeneficiaryNumber     string `json:"BeneficiaryNumber,omitempty"`
	BeneficiaryIFSC       string `json:"BeneficiaryIFSC,omitempty"`
	BeneficiaryMMID       string `json:"BeneficiaryMMID,omitempty"`
	BeneficiaryNumberType string `json:"BeneficiaryNumberType,omitempty"`
	Code                  string `json:"Code,omitempty"`
	Amount                string `json:"Amount,omitempty"`
	Qty                   string `json:"Qty,omitempty"`
	CustomerID            string `json:"CustomerID,omitempty"`
	ServiceType           string `json:"ServiceType,omitempty"`
	PaymentType           string `json:"PaymentType,omitempty"`
}

type CustomerDetails struct {
	Id            string           `json:"Id,omitempty"`
	Name          string           `json:"Name,omitempty"`
	Address       string           `json:"Address,omitempty"`
	MSISDN        string           `json:"MSISDN,omitempty"`
	Email         string           `json:"Email,omitempty"`
	Accounts      []*Accounts      `json:"Accounts,omitempty"`
	UserAddress   []*Address       `json:"UserAddress,omitempty"`
	Communication []*Communication `json:"Communication,omitempty"`
	Documents     []*Document      `json:"Documents,omitempty"`
	BankDetails   []*BankDetail    `json:"BankDetails,omitempty"`
	CreateDate    string           `json:"CreateDate,omitempty"`
	LastUpdate    string           `json:"LastUpdate,omitempty"`
	Profile       *Profile         `json:"Profile,omitempty"`
	Status        string           `json:"Status,omitempty"`
	Type          string           `json:"Type,omitempty"`
	Username      string           `json:"Username,omitempty"`
	UsernameType  string           `json:"UsernameType,omitempty"`
	ClassName     string           `json:"ClassName,omitempty"`
	KYCType       string           `json:"KYCType,omitempty"`
}

type TransactionDetails struct {
	RequestID                string `json:"RequestID,omitempty"`
	TransactionNumber        string `json:"TransactionNumber,omitempty"`
	HistoryTransactionNumber string `json:"HistoryTransactionNumber,omitempty"`
	RequestTime              string `json:"RequestTime,omitempty"`
	ResponseTime             string `json:"ResponseTime,omitempty"`
	OTP                      string `json:"OTP,omitempty"`
	OTPType                  string `json:"OTPType,omitempty"`
	TrasactionDescription    string `json:"TrasactionDescription,omitempty"`
	GL                       *GL    `json:"GLCustomer,omitempty"`
	Qty                      string `json:"Qty,omitempty"`
	Amount                   string `json:"Amount,omitempty"`
}
type Response struct {
	Code        string       `json:"Code,omitempty"`
	Message     string       `json:"Message,omitempty"`
	Status      string       `json:"Status,omitempty"`
	HostCode    string       `json:"HostCode,omitempty"`
	HostMessage string       `json:"HostMessage,omitempty"`
	Parameter   []*Parameter `json:"Parameter,omitempty"`
	Token       string       `json:"Token,omitempty"`
}

type Parameter struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
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

type Accounts struct {
	Id            string
	Balance       string `json:"Balance,omitempty"`
	CreateDate    string `json:"CreateDate,omitempty"`
	LastUpdate    string `json:"LastUpdate,omitempty"`
	LedgerBalance string `json:"LedgerBalance,omitempty"`
	Number        string `json:"Number,omitempty"`
	Status        string `json:"Status,omitempty"`
	Type          string `json:"Type,omitempty"`
	AccountType   string `json:"AccountType,omitempty"`
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

type Document struct {
	Id         string `json:"Id,omitempty"`
	CreateDate string `json:"CreateDate,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	Type       string `json:"Type,omitempty"`
	Value      string `json:"Value,omitempty"`
}

type BankDetail struct {
	Id         string
	CreateDate string `json:"CreateDate,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	Number     string `json:"Number,omitempty"`
	Name       string `json:"Name,omitempty"`
	Branch     string `json:"Branch,omitempty"`
	Code       string `json:"Code,omitempty"`
	Type       string `json:"Type,omitempty"`
	Address    string `json:"Address,omitempty"`
}

type GL struct {
	Commission      string `json:"Commission,omitempty"`
	Discount        string `json:"Discount,omitempty"`
	Fee             string `json:"Fee,omitempty"`
	STax            string `json:"STax,omitempty"`
	SurCharge       string `json:"SurCharge,omitempty"`
	TotalAmount     string `json:"TotalAmount,omitempty"`
	VAT             string `json:"VAT,omitempty"`
	Formulae        string `json:"Formulae,omitempty"`
	LedgerOperation string `json:"LedgerOperation,omitempty"`
}
