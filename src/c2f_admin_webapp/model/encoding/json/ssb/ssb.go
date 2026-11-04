package ssb

type Root struct {
	Tksp Tksp `json:"Tksp,omitempty"`
}

type Tksp struct {
	SSBInterface SSBInterface `json:"SSBInterface,omitempty"`
}

type SSBInterface struct {
	Channel         *Channel         `json:"Channel,omitempty"`
	Device          *Device          `json:"Device,omitempty"`
	Name            string           `json:"Name,omitempty"`
	Response        *Response        `json:"Response,omitempty"`
	Security        *Security        `json:"Security,omitempty"`
	TransactionData *TransactionData `json:"TransactionData,omitempty"`
	Users           []*Users         `json:"Users,omitempty"`
	Records         []*Records       `json:"Records,omitempty"`
	Notification    *Notification    `json:"Notification,omitempty"`
	Version         string           `json:"Version,omitempty"`
	Files           *Files           `json:"Files,omitempty"`
}

type Notification struct {
	Type    string `json:"Type,omitempty"`
	Dest    string `json:"Dest,omitempty"`
	Message string `json:"Message,omitempty"`
}

type TransactionData struct {
	RequestID              string `json:"RequestID,omitempty"`
	TimeStamp              string `json:"TimeStamp,omitempty"`
	TransactionDescription string `json:"TransactionDescription,omitempty"`
	TransactionNumber      string `json:"TransactionNumber,omitempty"`
}

type Security struct {
	CertificateExpiry string `json:"CertificateExpiry,omitempty"`
	Data              string `json:"Data,omitempty"`
	Hmac              string `json:"Hmac,omitempty"`
	Skey              string `json:"Skey,omitempty"`
}

type Response struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
	Param   string `json:"Param,omitempty"`
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
	Id                  string
	OrganizationName    string           `json:"OrganizationName,omitempty"`
	OrganizationDetails string           `json:"OrganizationDetails,omitempty"`
	Communication       []*Communication `json:"Communication,omitempty"`
	CreateDate          string           `json:"CreateDate,omitempty"`
	LastUpdate          string           `json:"LastUpdate,omitempty"`
	PayeeCode           string           `json:"PayeeCode,omitempty"`
	ADCode              string           `json:"ADCode,omitempty"`
	Profile             *Profile         `json:"Profile,omitempty"`
	Status              string           `json:"Status,omitempty"`
	Type                string           `json:"Type,omitempty"`
	Username            string           `json:"Username,omitempty"`
	IsSMS               string           `json:"IsSMS,omitempty"`
	IsSSB               string           `json:"IsSSB,omitempty"`
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

type Communication struct {
	Id         string
	CreateDate string `json:"CreateDate,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	Type       string `json:"Type,omitempty"`
	Value      string `json:"Value,omitempty"`
}

type Records struct {
	Id           string
	EC_Code      string `json:"EC_Code,omitempty"`
	Type         string `json:"Type,omitempty"`
	AD_Code      string `json:"AD_Code,omitempty"`
	PayeeCode    string `json:"PayeeCode,omitempty"`
	PolicyNumber string `json:"PolicyNumber,omitempty"`
	IDNumber     string `json:"IDNumber,omitempty"`
	FromDate     string `json:"FromDate,omitempty"`
	ToDate       string `json:"ToDate,omitempty"`
	Amount       string `json:"Amount,omitempty"`
	CreateDate   string `json:"CreateDate,omitempty"`
	Name         string `json:"Name,omitempty"`
	Branch       string `json:"Branch,omitempty"`
	Account      string `json:"Account,omitempty"`
	CreatedBy    string `json:"CreatedBy,omitempty"`
	UploadDate   string `json:"UploadDate,omitempty"`
	UploadBy     string `json:"UploadBy,omitempty"`
	Status       string `json:"Status,omitempty"`
	RejectReason string `json:"RejectReason,omitempty"`
}

type Files struct {
	Id                               string
	Path                             string `json:"Path,omitempty"`
	InstructionFileName              string `json:"InstructionFileName,omitempty"`
	ProvisionSuccessFileName         string `json:"ProvisionSuccessFileName,omitempty"`
	ProvisionErrorFileName           string `json:"ProvisionErrorFileName,omitempty"`
	UploadDate                       string `json:"UploadDate,omitempty"`
	UploadBy                         string `json:"UploadBy,omitempty"`
	ProvisionSuccessFileUploadDate   string `json:"ProvisionSuccessFileUploadDate,omitempty"`
	ProvisionSuccessFileUploadDateBy string `json:"ProvisionSuccessFileUploadDateBy,omitempty"`
	ProvisionErrorFileUploadDate     string `json:"ProvisionErrorFileUploadDate,omitempty"`
	ProvisionErrorFileUploadDateBy   string `json:"ProvisionErrorFileUploadDateBy,omitempty"`
	Status                           string `json:"Status,omitempty"`
	RecordCount                      string `json:"RecordCount,omitempty"`
	DownloadDate                     string `json:"DownloadDate,omitempty"`
}

type InvoiceDetails struct {
	Id               string
	PayeeCode        string
	OrgName          string
	FilePath         string //complete file path
	FileName         string
	Path             string //path from DB
	TxnCost          string
	TxnCostType      string
	TotalTxn         string
	ForMonth         string
	BillAmount       string
	Paid             string
	PaidAmount       string
	PaidDate         string
	DueAmount        string
	PaymentReceipt   string
	CreatedBy        string
	ModifiedBy       string
	CreateDate       string
	ModifiedDate     string
	InvoiceStatus    string
	InvoiceAuthority string
	VAT              string
}
