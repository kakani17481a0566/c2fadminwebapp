package ServiceApi

type Root struct {
	Tksp Tksp `json:"Tksp,omitempty"`
}

type Tksp struct {
	ServiceInterface ServiceInterface `json:"ServiceInterface,omitempty"`
}

type ServiceInterface struct {
	Channel         *Channel         `json:"Channel,omitempty"`
	Device          *Device          `json:"Device,omitempty"`
	Name            string           `json:"Name,omitempty"`
	Type            string           `json:"Type,omitempty"`
	Response        *Response        `json:"Response,omitempty"`
	Security        *Security        `json:"Security,omitempty"`
	TransactionData *TransactionData `json:"TransactionData,omitempty"`
	Agent           []*Users         `json:"Agent,omitempty"`
	AgentCash       []*Users         `json:"AgentCash,omitempty"`
	Distributor     []*Users         `json:"Distributor,omitempty"`
	Service         []*Users         `json:"Service,omitempty"`
	PartnerService  []*Users         `json:"PartnerService,omitempty"`
	System          []*Users         `json:"System,omitempty"`
	Commission      []*Users         `json:"Commission,omitempty"`
	Fee             []*Users         `json:"Fee,omitempty"`
	Discount        []*Users         `json:"Discount,omitempty"`
	VAT             []*Users         `json:"VAT,omitempty"`
	STax            []*Users         `json:"STax,omitempty"`
	SurCharge       []*Users         `json:"SurCharge,omitempty"`
	Marketing       []*Users         `json:"Marketing,omitempty"`
	CashPool        []*Users         `json:"CashPool,omitempty"`
	ECashPool       []*Users         `json:"ECashPool,omitempty"`
	Customer        []*Users         `json:"Customer,omitempty"`
	Merchant        []*Users         `json:"Merchant,omitempty"`
	Beneficiary     []*Users         `json:"Beneficiary,omitempty"`
	SysUsers        []*SysUsers      `json:"SysUsers,omitempty"`
	Notification    *Notification    `json:"Notification,omitempty"`

	Version string `json:"Version,omitempty"`
}

type Notification struct {
	Type    string `json:"Type,omitempty"`
	Dest    string `json:"Dest,omitempty"`
	Message string `json:"Message,omitempty"`
}

type TransactionData struct {
	CountryCode              string           `json:"CountryCode,omitempty"`
	Currency                 string           `json:"Currency,omitempty"`
	CustomerDetails          *CustomerDetails `json:"CustomerDetails,omitempty"`
	GLAgent                  *GL              `json:"GLAgent,omitempty"`
	GLDistributor            *GL              `json:"GLDistributor,omitempty"`
	GLService                *GL              `json:"GLService,omitempty"`
	GLPartnerService         *GL              `json:"GLPartnerService,omitempty"`
	GLSystem                 *GL              `json:"GLSystem,omitempty"`
	GLCommission             *GL              `json:"GLCommission,omitempty"`
	GLFee                    *GL              `json:"GLFee,omitempty"`
	GLDiscount               *GL              `json:"GLDiscount,omitempty"`
	GLSTax                   *GL              `json:"GLSTax,omitempty"`
	GLSurCharge              *GL              `json:"GLSurCharge,omitempty"`
	GLVAT                    *GL              `json:"GLVAT,omitempty"`
	GLCustomer               *GL              `json:"GLCustomer,omitempty"`
	GLMerchant               *GL              `json:"GLMerchant,omitempty"`
	GLBeneficiary            *GL              `json:"GLBeneficiary,omitempty"`
	GLMarketing              *GL              `json:"GLMarketing,omitempty"`
	GLCashPool               *GL              `json:"GLCashPool,omitempty"`
	GLAgentCash              *GL              `json:"GLAgentCash,omitempty"`
	GLECashPool              *GL              `json:"GLECashPool,omitempty"`
	PaymentMode              string           `json:"PaymentMode,omitempty"`
	Qty                      string           `json:"Qty,omitempty"`
	RequestID                string           `json:"RequestID,omitempty"`
	ServiceDetail            *ServiceDetail   `json:"ServiceDetail,omitempty"`
	TimeStamp                string           `json:"TimeStamp,omitempty"`
	TransactionDescription   string           `json:"TransactionDescription,omitempty"`
	TransactionNumber        string           `json:"TransactionNumber,omitempty"`
	HistoryTransactionNumber string           `json:"HistoryTransactionNumber,omitempty"`
	Amount                   string           `json:"Amount,omitempty"`
	OTP                      string           `json:"OTP,omitempty"`
	OTPType                  string           `json:"OTPType,omitempty"`
}

type ServiceDetail struct {
	BusinessCategory string `json:"BusinessCategory,omitempty"`
	Category         string `json:"Category,omitempty"`
	Denomination     string `json:"Denomination,omitempty"`
	Partner          string `json:"Partner,omitempty"`
	SubCategory      string `json:"SubCategory,omitempty"`
	Role             string `json:"Role,omitempty"`
	Group            string `json:"Group,omitempty"`
	Service          string `json:"Service,omitempty"`
	Code             string `json:"Code,omitempty"`
	UserType         string `json:"UserType,omitempty"`
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

type CustomerDetails struct {
	Email       string `json:"Email,omitempty"`
	MSISDN      string `json:"MSISDN,omitempty"`
	ServiceNo   string `json:"ServiceNo,omitempty"`
	Operator    string `json:"Operator,omitempty"`
	ServiceType string `json:"ServiceType,omitempty"`
	PaymentType string `json:"PaymentType,omitempty"`
	Name        string `json:"Name,omitempty"`
	Address     string `json:"Address,omitempty"`
	IFSC        string `json:"IFSC,omitempty"`
	Class       string `json:"Class,omitempty"`
	Section     string `json:"Section,omitempty"`
	SchoolName  string `json:"SchoolName,omitempty"`
	RollNumber  string `json:"RollNumber,omitempty"`
	Date        string `json:"Date,omitempty"`
}

type Security struct {
	CertificateExpiry string `json:"CertificateExpiry,omitempty"`
	Data              string `json:"Data,omitempty"`
	Hmac              string `json:"Hmac,omitempty"`
	Skey              string `json:"Skey,omitempty"`
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
	IP       IP       `json:"IP,omitempty"`
	ID       string   `json:"Id,omitempty"`
	OS       string   `json:"OS,omitempty"`
	Version  string   `json:"Version,omitempty"`
	Location Location `json:"Location,omitempty"`
}

type Location struct {
	Latitude         string `json:"Latitude,omitempty"`
	Longitude        string `json:"Longitude,omitempty"`
	City             string `json:"City,omitempty"`
	Country          string `json:"Country,omitempty"`
	Countrycode      string `json:"Countrycode,omitempty"`
	Internetprovider string `json:"Internetprovider,omitempty"`
	Ip               string `json:"Ip,omitempty"`
	Region           string `json:"Region,omitempty"`
	Regionname       string `json:"Regionname,omitempty"`
	Timezone         string `json:"Timezone,omitempty"`
	Zip              string `json:"Zip,omitempty"`
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
	Id                    string
	Accounts              []*Accounts          `json:"Accounts,omitempty"`
	Address               []*Address           `json:"Address,omitempty"`
	Communication         []*Communication     `json:"Communication,omitempty"`
	Documents             []*Document          `json:"Documents,omitempty"`
	BankDetails           []*BankDetail        `json:"BankDetails,omitempty"`
	BenficiaryDetails     []*BenficiaryDetails `json:"BenficiaryDetails,omitempty"`
	CreateDate            string               `json:"CreateDate,omitempty"`
	LastUpdate            string               `json:"LastUpdate,omitempty"`
	ParentID              string               `json:"ParentId,omitempty"`
	Profile               *Profile             `json:"Profile,omitempty"`
	Status                string               `json:"Status,omitempty"`
	Type                  string               `json:"Type,omitempty"`
	Username              string               `json:"Username,omitempty"`
	UsernameType          string               `json:"UsernameType,omitempty"`
	ClassName             string               `json:"ClassName,omitempty"`
	KYCType               string               `json:"KYCType,omitempty"`
	KYCStatus             string               `json:"KYCStatus,omitempty"`
	ManagerName           string               `json:"ManagerName,omitempty"`
	ManagerDesignation    string               `json:"ManagerDesignation,omitempty"`
	NatureOfBusiness      string               `json:"NatureOfBusiness,omitempty"`
	VATRegistrationNumber string               `json:"VATRegistrationNumber,omitempty"`
	OperatingChannel      string               `json:"OperatingChannel,omitempty"`
	OperatingDays         string               `json:"OperatingDays,omitempty"`
	LoginCount            string               `json:"LoginCount,omitempty"`
	QRCode                string               `json:"QRCode,omitempty"`
	SuspenseAccount       string               `json:"SuspenseAccount,omitempty"`
	CBSAccountNumber      string               `json:"CBSAccountNumber,omitempty"`
	AccountType           string               `json:"AccountType,omitempty"`
	AgentType             string               `json:"AgentType,omitempty"`
	CustomerSegment       string               `json:"CustomerSegment,omitempty"`
}
type BenficiaryDetails struct {
	AccNumber string `json:"AccNumber,omitempty"`
	Name      string `json:"Name,omitempty"`
	NickName  string `json:"NickName,omitempty"`
	BankName  string `json:"BankName,omitempty"`
	Code      string `json:"Code,omitempty"`
}

type FingerData struct {
	RightThumb         string `json:"RightThumb,omitempty"`
	Right_IndexFinger  string `json:"Right_IndexFinger,omitempty"`
	Right_MiddleFinger string `json:"Right_MiddleFinger,omitempty"`
	Right_RingFinger   string `json:"Right_RingFinger,omitempty"`
	Right_LittleFinger string `json:"Right_LittleFinger,omitempty"`
	LeftThumb          string `json:"Dest,omitempty"`
	Left_IndexFinger   string `json:"Left_IndexFinger,omitempty"`
	Left_MiddleFinger  string `json:"Left_MiddleFinger,omitempty"`
	Left_RingFinger    string `json:"Left_RingFinger,omitempty"`
	Left_LittleFinger  string `json:"Left_LittleFinger,omitempty"`
}

type Profile struct {
	Id              string
	CreateDate      string `json:"CreateDate,omitempty"`
	DOB             string `json:"DOB,omitempty"`
	FirstName       string `json:"FirstName,omitempty"`
	Gender          string `json:"Gender,omitempty"`
	LastName        string `json:"LastName,omitempty"`
	LastUpdate      string `json:"LastUpdate,omitempty"`
	MiddleName      string `json:"MiddleName,omitempty"`
	NationalID      string `json:"NationalID,omitempty"`
	TPin            string `json:"TPin,omitempty"`
	CustomerSegment string `json:"Customersegment,omitempty"`
}

type Accounts struct {
	Id             string
	Balance        string `json:"Balance,omitempty"`
	CreateDate     string `json:"CreateDate,omitempty"`
	LastUpdate     string `json:"LastUpdate,omitempty"`
	LedgerBalance  string `json:"LedgerBalance,omitempty"`
	Number         string `json:"Number,omitempty"`
	Status         string `json:"Status,omitempty"`
	Type           string `json:"Type,omitempty"`
	AccountType    string `json:"AccountType,omitempty"`
	CBSAccountType string `json:"AccountType,omitempty"`
	Min            string `json:"Min,omitempty"`
	Max            string `json:"Max,omitempty"`
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
	BeneName   string `json:"Name,omitempty"`
	Branch     string `json:"Branch,omitempty"`
	Code       string `json:"Code,omitempty"`
	Type       string `json:"Type,omitempty"`
	Address    string `json:"Address,omitempty"`
}

type SysUsers struct {
	Id            string
	Communication []*Communication `json:"Communication,omitempty"`
	CreateDate    string           `json:"CreateDate,omitempty"`
	LastUpdate    string           `json:"LastUpdate,omitempty"`
	Profile       *Profile         `json:"Profile,omitempty"`
	Status        string           `json:"Status,omitempty"`
	Type          string           `json:"Type,omitempty"`
	Username      string           `json:"Username,omitempty"`
	UsernameType  string           `json:"UsernameType,omitempty"`
	RoleName      string           `json:"ClassName,omitempty"`
}
