package PortalApi

type Response struct {
	Code        string
	Message     string
	Status      string
	TxnNumber   string
	Balance     string
	Name        string
	Email       string
	Mobile      string
	Dob         string
	SessionID   string
	Gender      string
	Key         string
	Salt        string
	ProductInfo string
	Hmac        string
	Data        string
	Skey        string
	ReqID       string
	Rows        []*Row
	Fee         string
	Stax        string
	SurCharge   string
	KYCType     string
	Discount    string
	Commission  string
	BenMobile   string
	QRData      string
	QRId        string
	Amount      string
	Partner     string
	TimeStamp   string
}

type Row struct {
	Cols []*Col
}

type Col struct {
	Key   string
	Value string
}

type Request struct {
	AgentID            string
	MerchantID         string
	CustomerID         string
	TxnNumber          string
	RequestID          string
	Mobile             string
	UserType           string
	Password           string
	NewPassword        string
	NewPasswordC       string
	ChannelName        string
	ChannelKey         string
	OS                 string
	OSVersion          string
	IMEI               string
	IMSI               string
	PIP                string
	ServerDateTime     string
	SessionID          string
	SessionStatus      string
	ServiceName        string
	Amount             string
	SchoolName         string
	Class              string
	Section            string
	RollNumber         string
	FeePeriod          string
	OTP                string
	BenMobile          string
	Operator           string
	ServiceType        string
	ServiceNumber      string
	BenType            string
	BenAccounNumbert   string
	BenAccountIFSC     string
	Package            string
	PackageType        string
	Code               string
	BenMMID            string
	BenName            string
	TransactionRemarks string
	Email              string
	FirstName          string
	MiddleName         string
	LastName           string
	Dob                string
	Gender             string
	Skey               string
	Data               string
	Hmac               string
	PayUid             string `form:"mihpayid" valid:"Required"`
	PaymentMode        string `form:"mode" valid:"Required"`
	PayuStatus         string `form:"status" valid:"Required"`
	PayUKey            string `form:"key" valid:"Required"`
	PayuTxnNo          string `form:"txnid" valid:"Required"`
	PayuAmount         string `form:"amount" valid:"Required"`
	PayuProductInfo    string `form:"productinfo" valid:"Required"`
	PayuFirstName      string `form:"firstname"`
	PayuEmail          string `form:"email"`
	PayuPhone          string `form:"phone"`
	PayuHash           string `form:"hash" valid:"Required"`
	PayuError          string `form:"error" valid:"Required"`
	NIBLAMOUNT         string `form:"AMOUNT" `
	NIBLBankId         string `form:"BankId" `
	NIBLSTFLG          string `form:"STFLG" `
	NIBLPREF           string `form:"PREF" `
	NIBLBID            string `form:"BID" `
	NIBLVAR1           string `form:"VAR1" `
	NIBLVAR2           string `form:"VAR2" `
	NIBLVAR4           string `form:"VAR4" `
	NIBLVAR3           string `form:"VAR3" `
	NIBLAMT1           string `form:"AMT1" `
	NIBLFirstName      string `form:"FirstName" `
	NIBLLastName       string `form:"LastName" `
	NIBLPhone          string `form:"Phone" `
	NIBLMobile         string `form:"Mobile" `
	NIBLCustRefType    string `form:"CustRefType" `
	NIBLEmailID        string `form:"EmailID" `
	NIBLUID            string `form:"UID" `
	NIBLBIDSYS         string `form:"BIDSYS" `
}
