package PortalApi

type Response struct {
	Code             string
	Message          string
	Upic             string
	TPinImage        string
	TPin             string
	TPinNumber       string
	Status           string
	Balance          string
	CashBalance      string
	Name             string
	FirstName        string
	MiddleName       string
	LastName         string
	Email            string
	Mobile           string
	Dob              string
	SessionID        string
	Gender           string
	AgentID          string
	CustomerID       string
	WalletType       string
	WalletStatus     string
	Key              string
	Salt             string
	ProductInfo      string
	BenAccounNumbert string
	BenBankName      string
	BenName          string
	BenAccountIFSC   string
	TxnNo            string
	Hmac             string
	Data             string
	Skey             string
	ReqID            string
	Rows             []*Row
	Plan             []*Plan
	BillerDetails    []*BillerDetails
	Fee              string
	Stax             string
	SurCharge        string
	KYCType          string
	Discount         string
	Commission       string
	BenMobile        string
	QRData           string
	QRId             string
	TxnNumber        string
	Amount           string
	TotalAmount      string
	Partner          string
	PartnerRefNo     string
	PartnerStatus    string
	TimeStamp        string
	UserType         string
	AgentName        string

	PasswordState string
	UserId        string
	DocType       string
	DocNo         string
	AccNo         string
	Count         string
	KYCRefNo      string
	NationalID    string
	CreateDate    string

	ServiceNo string

	BillDate  string
	KYCMethod string
	KYCNumber string
	KYCCode   string

	Address1            string
	Address2            string
	Post                string
	PIN                 string
	District            string
	State               string
	Nationlaity         string
	Mode                string
	ServiceName         string
	ServiceNumber       string
	CustomerSegment     string
	CUstomerSegmentDesc string
}

type Plan struct {
	PlanType            string
	RechargeAmount      string `json:"rechargeAmount"`
	TalkTime            string `json:"talkTime"`
	Validity            string `json:"validity"`
	RechargePlanID      string `json:"rechargePlanId"`
	RechargeDescription string `json:"rechargeDescription"`
}

type Row struct {
	Cols []*Col
}

type Col struct {
	Key   string
	Value string
}

type BillerDetails struct {
	BillerId           string
	BillerFullName     string
	BillerNickName     string
	BillerCategory     string
	BillerPresence     string
	KeyParameter       string
	NoOfAuthenticators string
	Authenticator1     *Authenticator1
	Authenticator2     *Authenticator2
	Authenticator3     *Authenticator3
}

type Authenticator1 struct {
	LableName         string
	RegExValidation   string
	ValueType         string
	ValidationMessage string
}
type Authenticator2 struct {
	LableName         string
	RegExValidation   string
	ValueType         string
	ValidationMessage string
}

type Authenticator3 struct {
	LableName         string
	RegExValidation   string
	ValueType         string
	ValidationMessage string
}

type Request struct {
	AgentID      string
	AdminID      string
	CustomerID   string
	Mobile       string
	DOB          string
	Partner      string
	Password     string
	TxnPassword  string
	UserType     string
	Wallet       string
	WalletType   string
	WalletStatus string
	TxnPass      string
	Class        string
	BenBankName  string
	ValidDocType string
	ValidDocNo   string
	DocumentType string
	DocumentNo   string
	KYCMethod    string
	KYCType      string
	KYCNumber    string
	KYCCode      string
	KYCRefNo     string
	Photo        string
	ToDate       string
	FromDate     string
	Remarks      string
	Upic         string
	AccountType  string
	TPinNumber   string
	Address1     string
	Address2     string
	Post         string
	District     string
	PIN          string
	State        string
	Province     string
	City         string
	Nationality  string
	TPinImage    string
	UserId       string
	QRCode       string

	TransactionNumber string
	BillDetails       string

	NewPassword      string
	NewPasswordC     string
	PassType         string
	NewTxnPassword   string
	NewTxnPasswordC  string
	DepositAmount    string
	WithdrawalAmount string
	CBSAccount       string

	RRN string

	ChannelName         string
	ChannelKey          string
	OS                  string
	OSVersion           string
	APPVersion          string
	IMEI                string
	IMSI                string
	PIP                 string
	ServerDateTime      string
	SessionID           string
	SessionStatus       string
	ServiceName         string
	Amount              string
	OTP                 string
	BenMobile           string
	BenWallet           string
	AuthDetails         string
	BillNumber          string
	Operator            string
	ServiceType         string
	CircleId            string
	RequestID           string
	ServiceNumber       string
	ServiceDetails      string
	AccountNumber       string
	BenType             string
	BenAccounNumbert    string
	BenAccountIFSC      string
	Package             string
	PackageType         string
	Code                string
	BenMMID             string
	BenName             string
	TransactionRemarks  string
	Email               string
	FirstName           string
	MiddleName          string
	LastName            string
	Dob                 string
	Gender              string
	Skey                string
	Data                string
	Hmac                string
	CustomerSegment     string
	CUstomerSegmentDesc string
}
