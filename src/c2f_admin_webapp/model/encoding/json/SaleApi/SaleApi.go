package SaleApi

type Root struct {
	Tksp Tksp `json:"Tksp,omitempty"`
}

type Tksp struct {
	ServiceInterface *ServiceInterface `json:"ServiceInterface,omitempty"`
}

type ServiceInterface struct {
	Channel         *Channel         `json:"Channel,omitempty"`
	Device          *Device          `json:"Device,omitempty"`
	Name            string           `json:"Name,omitempty"`
	Type            string           `json:"Type,omitempty"`
	Operator        string           `json:"Operator,omitempty"`
	Response        *Response        `json:"Response,omitempty"`
	Security        *Security        `json:"Security,omitempty"`
	TransactionData *TransactionData `json:"TransactionData,omitempty"`
	Users           *User            `json:"Users,omitempty"`
	//Partner         *Partner         `json:"Partner,omitempty"`
	SysUsers []*SysUsers `json:"SysUsers,omitempty"`
	Version  string      `json:"Version,omitempty"`
}

type TransactionData struct {
	CountryCode              string            `json:"CountryCode,omitempty"`
	Currency                 string            `json:"Currency,omitempty"`
	ProductDetails           []*ProductDetails `json:"ProductDetails,omitempty"`
	RequestID                string            `json:"RequestID,omitempty"`
	TimeStamp                string            `json:"TimeStamp,omitempty"`
	TransactionDescription   string            `json:"TransactionDescription,omitempty"`
	TransactionNumber        string            `json:"TransactionNumber,omitempty"`
	HistoryTransactionNumber string            `json:"HistoryTransactionNumber,omitempty"`
	FileName                 string            `json:"FileName,omitempty"`
	FilePath                 string            `json:"FilePath,omitempty"`
	FTPSite                  string            `json:"FTPSite,omitempty"`
	FTPUsername              string            `json:"FTPUsername,omitempty"`
	FTPPassword              string            `json:"FTPPassword,omitempty"`
	RequestTime              string            `json:"RequestTime,omitempty"`
	ResponseTime             string            `json:"ResponseTime,omitempty"`
	TotalAmount              string            `json:"TotalAmount,omitempty"`
	PaidAmount               string            `json:"PaidAmount,omitempty"`
	RemainingAmount          string            `json:"RemainingAmount,omitempty"`
	TotRemainingAmount       string            `json:"TotRemainingAmount,omitempty"`
}

type ProductDetails struct {
	Name          string `json:"Name,omitempty"`
	BarCode       string `json:"BarCode,omitempty"`
	SerialNumber  string `json:"SerialNumber,omitempty"`
	BatchNumber   string `json:"BatchNumber,omitempty"`
	ActualQty     string `json:"ActualQty,omitempty"`
	ActualPrice   string `json:"ActualPrice,omitempty"`
	LogisticPrice string `json:"LogisticPrice,omitempty"`
	SellingPrice  string `json:"SellingPrice,omitempty"`
	ExpiryDate    string `json:"ExpiryDate,omitempty"`
	Partner       string `json:"Partner,omitempty"`
	Qty           string `json:"Qty,omitempty"`
	RemainingQty  string `json:"RemainingQty,omitempty"`
	GLProduct     *GL    `json:"GLProduct,omitempty"`
	Status        string `json:"Status,omitempty"`
	GlRule        string `json:"GlRule,omitempty"`
	Remarks       string `json:"Remarks,omitempty"`
	Description   string `json:"Description,omitempty"`
	BatchNo       string `json:"BatchNo,omitempty"`
	Id            string `json:"Id,omitempty"`
	PaymentStatus string `json:"PaymentStatus,omitempty"`
	ManfDate      string `json:"ManfDate,omitempty"`
}

type ProductDetailsResp struct {
	Name          string `json:"Name,omitempty"`
	BarCode       string `json:"BarCode,omitempty"`
	SerialNumber  string `json:"SerialNumber,omitempty"`
	BatchNumber   string `json:"BatchNumber,omitempty"`
	Qty           string `json:"Qty,omitempty"`
	ExpiryDate    string `json:"ExpiryDate,omitempty"`
	Partner       string `json:"Partner,omitempty"`
	Price         string `json:"Price,omitempty"`
	RemainingQty  string `json:"RemainingQty,omitempty"`
	GLProduct     *GL    `json:"GLProduct,omitempty"`
	Remarks       string `json:"Remarks,omitempty"`
	Description   string `json:"Description,omitempty"`
	BatchNo       string `json:"BatchNo,omitempty"`
	Id            string `json:"Id,omitempty"`
	CreateDate    string `json:"CreateDate,omitempty"`
	LastUpdate    string `json:"LastUpdate,omitempty"`
	CreatedBy     string `json:"CreatedBy,omitempty"`
	UpdatedBy     string `json:"UpdatedBy,omitempty"`
	PaymentStatus string `json:"PaymentStatus,omitempty"`
	ManfDate      string `json:"ManfDate,omitempty"`
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

type Security struct {
	CertificateExpiry string `json:"CertificateExpiry,omitempty"`
	Data              string `json:"Data,omitempty"`
	Hmac              string `json:"Hmac,omitempty"`
	Skey              string `json:"Skey,omitempty"`
}

type Response struct {
	Code           string                `json:"Code,omitempty"`
	Message        string                `json:"Message,omitempty"`
	Status         string                `json:"Status,omitempty"`
	HostCode       string                `json:"HostCode,omitempty"`
	HostMessage    string                `json:"HostMessage,omitempty"`
	Parameter      []*Parameter          `json:"Parameter,omitempty"`
	ProductDetails []*ProductDetailsResp `json:"ProductDetails,omitempty"`
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

type User struct {
	Id              string
	Accounts        []*Accounts      `json:"Accounts,omitempty"`
	Address         []*Address       `json:"Address,omitempty"`
	Communication   []*Communication `json:"Communication,omitempty"`
	CreateDate      string           `json:"CreateDate,omitempty"`
	LastUpdate      string           `json:"LastUpdate,omitempty"`
	Profile         *Profile         `json:"Profile,omitempty"`
	Status          string           `json:"Status,omitempty"`
	Type            string           `json:"Type,omitempty"`
	Username        string           `json:"Username,omitempty"`
	UsernameType    string           `json:"UsernameType,omitempty"`
	KycType         string           `json:"KycType,omitempty"`
	LoginCount      string           `json:"LoginCount,omitempty"`
	RemainingAmount string           `json:"RemainingAmount,omitempty"`
	CreatedBy       string           `json:"CreatedBy,omitempty"`
	UpdatedBy       string           `json:"UpdatedBy,omitempty"`
}

type Profile struct {
	Id         string `json:"Id,omitempty"`
	CreateDate string `json:"CreateDate,omitempty"`
	DOB        string `json:"DOB,omitempty"`
	FirstName  string `json:"FirstName,omitempty"`
	Gender     string `json:"Gender,omitempty"`
	LastName   string `json:"LastName,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	MiddleName string `json:"MiddleName,omitempty"`
	NationalID string `json:"NationalID,omitempty"`
	FullName   string `json:"FullName,omitempty"`
}

type Accounts struct {
	Id            string `json:"Id,omitempty"`
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
	Id         string `json:"Id,omitempty"`
	Address1   string `json:"Address1,omitempty"`
	Address2   string `json:"Address2,omitempty"`
	City       string `json:"City,omitempty"`
	Country    string `json:"Country,omitempty"`
	CreateDate string `json:"CreateDate,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	Province   string `json:"Province,omitempty"`
	Type       string `json:"Type,omitempty"`
	Zip        string `json:"Zip,omitempty"`
	Area       string `json:"Area,omitempty"`
}

type Communication struct {
	Id         string `json:"Id,omitempty"`
	CreateDate string `json:"CreateDate,omitempty"`
	LastUpdate string `json:"LastUpdate,omitempty"`
	Type       string `json:"Type,omitempty"`
	Value      string `json:"Value,omitempty"`
}

type SysUsers struct {
	Id            string
	Communication []*Communication `json:"Communication,omitempty"`
	CreateDate    string           `json:"CreateDate,omitempty"`
	LastUpdate    string           `json:"LastUpdate,omitempty"`
	FromDate      string           `json:"FromDate,omitempty"`
	ToDate        string           `json:"ToDate,omitempty"`
	CreatedBy     string           `json:"CreatedBy,omitempty"`
	UpdatedBy     string           `json:"UpdatedBy,omitempty"`
	Status        string           `json:"Status,omitempty"`
	Type          string           `json:"Type,omitempty"`
	Username      string           `json:"Username,omitempty"`
	UsernameType  string           `json:"UsernameType,omitempty"`
	RoleName      string           `json:"RoleName,omitempty"`
	MainUser      string           `json:"MainUser,omitempty"`
	Remarks       string           `json:"Remarks,omitempty"`
	DOB           string           `json:"DOB,omitempty"`
	FirstName     string           `json:"FirstName,omitempty"`
	Gender        string           `json:"Gender,omitempty"`
	LastName      string           `json:"LastName,omitempty"`
	MiddleName    string           `json:"MiddleName,omitempty"`
	Email         string           `json:"Email,omitempty"`
	MobileNo      string           `json:"MobileNo,omitempty"`
	Accounts      *Accounts        `json:"Accounts,omitempty"`
	Organization  string           `json:"Organization,omitempty"`
	CountryCode   string           `json:"CountryCode,omitempty"`
}

type Partner struct {
	Id         string    `json:"Id,omitempty"`
	Name       string    `json:"Name,omitempty"`
	Remarks    string    `json:Remarks,omitempty`
	CreateDate string    `json:"CreateDate,omitempty"`
	LastUpdate string    `json:"LastUpdate,omitempty"`
	CreatedBy  string    `json:"CreatedBy,omitempty"`
	UpdatedBy  string    `json:"UpdatedBy,omitempty"`
	Accounts   *Accounts `json:"Accounts,omitempty"`
}
