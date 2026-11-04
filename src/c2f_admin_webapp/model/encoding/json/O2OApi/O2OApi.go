package O2OApi

type Response struct {
	Code            string
	Message         string
	Stores          []*StoreRes
	DeliveryAddress []*DeliveryAddress
	//Products        []*ProductSubCategory
	ProductsArr  []*ProductsArrayByCategory
	Categories   []*CategoryRes
	Feedback     []*Feedback
	Orders       []*Orders
	Products     []*Products
	Clusters     []*Clusters
	ViewProfile  []*ViewProfile
	LatLngValues []*LatLngValues

	ProductRes        []*ProductRes
	Status            string
	Suggestion        []string
	AssignOrder       string
	AcceptedOrder     string
	IssuedOrder       string
	ReceivedOrder     string
	CompletedOrder    string
	DeclinedOrder     string
	CollectedAmount   string
	DirectOrders      string
	TotalSupplyAmount string
	Name              string
	Email             string
	Mobile            string
	Userpic           string
	Dob               string
	SessionID         string
	Gender            string
	Post              string
	PIN               string
	District          string
	State             string
	Nationlaity       string
	Locality          string
	HouseNo           string
	Landmark          string
	OrderNumber       string
	Latitude          float64
	Longitude         float64
	StoreImageUrl     string
	VersionNumber     string
	ForceUpgrade      string
	Veg               bool
	NonVeg            bool
	Egg               bool
	Eggless           bool
	Key               string
	ReqID             string
	Currency          string
	ClassProducts     string
	Threshold         string
	UserType          string
	ClusterName       string
	Address1          string
	Address2          string
	City              string
	Clustername       string
	Createdate        string
	Customerid        string
	MapmyindiaToken   string
	Delhi             string
	AndraPradesh      string
	TotalAmount       string
	TotalBillAmount   string
	SupplyCharge      string
	PckgCharge        string
	BankCharge        string
	Discount          string
	CGST              string
	SGST              string
	CESS              string
	OtherTax          string
	ProductName       string
	ProductQty        string
	Units             string
	SKU               string
	Timings           string
	Location          string
	PayUTax           string
	PaymentType       string
	DeliveryCharges   string
}

type StoreRes struct {
	Id                       string    `json:"id"`
	Address                  string    `json:"address"`
	Locality                 string    `json:"locality"`
	Landmark                 string    `json:"landmark"`
	StoreName                string    `json:"store_name"`
	City                     string    `json:"city"`
	Url                      string    `json:"url"`
	Rating                   string    `json:"rating"`
	MinimumOrder             string    `json:"minimumOrder"`
	DeliveryTime             string    `json:"deliveryTime"`
	Location                 *Location `json:"location"`
	Description              string    `json:"description"`
	Mobile                   string    `json:"mobile"`
	Status                   string    `json:"status"`
	SuperCategory            string    `json:"superCategory"`
	FlowType                 string    `json:"flowType"`
	StoreDescription         string    `json:"storeDescription"`
	StoreTimings             string    `json:"storeTimings"`
	StoreState               string    `json:"storeState"`
	Taxes                    string    `json:"taxes"`
	DeliveryType             string    `json:"deliveryType"`
	PickupTime               string    `json:"pickupTime"`
	OnlinePayment            bool      `json:"onlinePayment"`
	FixedSlots               string    `json:"fixedSlots"`
	AvailableSlots           string    `json:"availableSlots"`
	SlotDays                 string    `json:"slotDays"`
	BusinessDays             string    `json:"businessDays"`
	StockManagement          string    `json:"stockmanagement"`
	StockManagementThreshold string    `json:"stockmanagementthreshold"`
}

type DeliveryAddress struct {
	Id        string
	Address   string
	Locality  string
	City      string
	Landmark  string
	Allow     int
	Latitude  float64
	Longitude float64
}

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type Clusters struct {
	Id          string
	ClusterName string
}

type LatLngValues struct {
	Clustername  string
	Customername string
	Values       string
}

type ViewProfile struct {
	Customerid     string
	Agentid        string
	Deliveryuserid string
	Name           string
	Mobile         string
	Email          string
	Address1       string
	Address2       string
	City           string
	Clustername    string
	Createdate     string
	UserPic        string
	State          string
	Pincode        string
}

type Products struct {
	ClusterProductId           string
	P2P2Cid                    string
	ClusterId                  string
	ProductId                  string
	Status                     string
	CategoryName               string
	SubcategoryName            string
	ProductName                string
	ProductQty                 string
	ProductPackageQty          string
	ProductQualityType         string
	Product_QualityMeasure     string
	ShelLife                   string
	PackageType                string
	PackageQty                 string
	Krishivala_BuyingPrice     string
	Cluster_BuyingPrice        string
	Cluster_Name               string
	ProductSKU                 string
	Cluster_SellingPrice       string
	Krishivala_CGST_charge     string
	Krishivala_SGST_charge     string
	Krishivala_CESS_charge     string
	Krishivala_OtherTAX_charge string
	ClusterSellingPrice        string
	ClusterSupplyCharge        string
	ClusterPackagingCharge     string
	ClusterBankCharge          string
	ClusterDiscount            string
	DeliveryDiscount           string
	ClusterBuyingPrice         string
	ClusterMarketPrice         string
	ClusterName                string
	Nounits                    string
	ProductTotalPrice          string
	ProductTotalSellPrice      string
	FromDate                   string
	ToDate                     string
	Time                       string
	Repeat                     string
	Payment_Status             string
	Payment_Type               string
	Order_Type                 string
	TimeStamp                  string
	Order_ID                   string
	PdctImg1                   string
	PdctImg2                   string
	PdctImg3                   string
	PdctImg4                   string
	PdctVideo1                 string
	CategoryImage              string
	SubCategoryImage           string
	ProducerId                 string
	ProducerName               string
	ProducerMobile             string
	ProducerSellingPrice       string
	ProducerCGST               string
	ProducerSGST               string
	ProducerCESS               string
	ProducerOtherTax           string
	PurchaseType               string
	Add1                       string
	Add2                       string
	Town                       string
	City                       string
	State                      string
	Pincode                    string
	SupplyForecastTxnId        string
	ClusterType                string
	Remarks                    string
	PaymentRemarks             string
	PayUTxnId                  string
}

type Orders struct {
	Txnid               string
	Id                  string
	OrderNumber         string
	Order_id            string
	SupplyOrdertxn_id   string
	Order_date          string
	DeliveryUserMob     string
	DeliveryassignDate  string
	Status              string
	Name                string
	Add1                string
	Add2                string
	Town                string
	City                string
	State               string
	Mobile              string
	Pincode             string
	Latlong             string
	LocationURL         string
	ClusterName         string
	ProductName         string
	ProductPackageQty   string
	ProductQty          string
	ProductSupplyUnits  string
	DeliveryTime        string
	DeliveryFrom        string
	DeliveryTo          string
	ProductPrice        string
	Productunits        string
	Cart                *Cart
	OrderTimeStamp      string
	Payment             string
	CurrentStatus       string
	Deliverydate        string
	Remarks             string
	UnitPrice           string
	TotalPrice          string
	ProductUnit         string
	CustomerMobile      string
	DeliveryCharges     string
	DeliveryInformation string
	ReviewedSellPrice   string
	Amount              string
	ClusterBuyingPrice  string
	TotalBuyPrice       string
	PaymentType         string
	PaymentStatus       string
	SKU                 string
	PaymentRemarks      string
	Timestamp           string
	TotalCollectAmount  string
	Category            string
	SubCategory         string
	ReceivedUnits       string
	UpdatedPrice        string
	Supplycharge        string
	Pckgingcharge       string
	Bankcharge          string
	Discount            string
	CGST                string
	SGST                string
	CESS                string
	Othertax            string
	SellingPrice        string
	Startdate           string
	Enddate             string
	PurchaseType        string
	RepeatOrder         string
	OrderType           string
	Days                string
	CreateDate          string
	DeliveryDiscount    string
	PayUTotalTax        string
	PayUTransactionId   string

	// LatLngValues       string
	//	Status             *Status
}

type Status struct {
	Placed                  string
	PlacedLable             string
	PlacedDesc              string
	Received                string
	ReceivedLable           string
	ReceivedDesc            string
	Packed                  string
	PackedLable             string
	PackedDesc              string
	OutForDelivery          string
	OutForDeliveryLable     string
	OutForDeliveryDesc      string
	Delivered               string
	DeliveredLable          string
	DeliveredDesc           string
	PackedTimeStamp         string
	PlacedTimeStamp         string
	ReceivedTimeStamp       string
	OutForDeliveryTimeStamp string
	DeliveredTimeStamp      string
	Cancel                  string
	CancelDesc              string
	CancelLable             string
	CancelBy                string
	CancelTimeStamp         string
	ReadyForPickup          string
	ReadyForPickupDesc      string
	ReadyForPickupLable     string
	ReadyForPickupTimeStamp string
}

type Cart struct {
	ItemsCount      int
	TotalAmt        string
	Category        []*Category
	StoreInfo       *StoreInfo
	DeliveryAddress *DeliveryAddress
	ContactMobile   string
	ContactEmail    string
	TotalGST        string
	TotalSGST       string
	TotalCGST       string
	Taxes           []*Taxes
	GrandTotal      string
	PaymentType     string
	SelectedSlot    string
	DeliverType     string
	ClassName       string
}

type Taxes struct {
	Name   string
	Value  string
	Type   string
	Amount string
}

type StoreInfo struct {
	Address           string  `json:"Address"`
	ShopName          string  `json:"ShopName"`
	City              string  `json:"City"`
	ImageUrl          string  `json:"ImageUrl"`
	Rating            string  `json:"Rating"`
	MinimumOrder      string  `json:"MinimumOrder"`
	DeliveryTime      string  `json:"DeliveryTime"`
	Latitude          float64 `json:"Latitude"`
	Longitude         float64 `json:"Longitude"`
	StoreDeliverType  string  `json:"StoreDeliverType"`
	StorePickupTime   string  `json:"StorePickupTime"`
	StoreMobileNumber string  `json:"StoreMobileNumber"`
}

type Category struct {
	Name                  string
	ProductSubCategoryRes *[]ProductSubCategory `json:"SubCategories"`
}

// type ProductSubCategoryRes struct {
// 	Id       string `json:"Id"`
// 	Name     string `json:"Name"`
// 	Title    string `json:"Title"`
// 	Products *ProductRes
// }

type ProductSubCategory struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Title    string `json:"Title"`
	Products []*ProductRes
}

type ProductsArrayByCategory struct {
	Name        string `json:"Name"`
	SubCategory []*SubCategory
}

type SubCategory struct {
	Name     string `json:"Name"`
	Products []*ProductCatRes
}

type Feedback struct {
	Retailer string `json:"retailer"`
	Rating   []*Rating
}

type Rating struct {
	Value float64 `json:"value"`
}

type ProductCatRes struct {
	Id                   string    `json:"Id"`
	Name                 string    `json:"Name"`
	ImageURL             string    `json:"ImageUrl"`
	TotalPrice           string    `json:"TotalPrice"`
	Desc                 string    `json:"Desc"`
	Unit                 string    `json:"Unit"`
	UnitArray            []UnitArr `json:"UnitArray"`
	Price                string    `json:"Price"`
	Added                string    `json:"Added"`
	Qty                  int       `json:"Qty"`
	FoodType             string    `json:"FoodType"`
	OpeningStock         string    `json:"openingstock"`
	PCGST                string
	PSGST                string
	SCGST                string
	SSGST                string
	OfferPercentagePrice string
	OfferFlatPrice       string
	Offer                string
	MRP                  string `json:"mrp"`
}

type ProductRes struct {
	Id                   string    `json:"Id"`
	Name                 string    `json:"Name"`
	ImageURL             string    `json:"ImageUrl"`
	TotalPrice           string    `json:"TotalPrice"`
	Desc                 string    `json:"Desc"`
	Unit                 string    `json:"Unit"`
	UnitArray            []UnitArr `json:"UnitArray"`
	Price                string    `json:"Price"`
	Added                string    `json:"Added"`
	Qty                  int       `json:"Qty"`
	OpeningStock         string    `json:"openingstock"`
	PCGST                string
	PSGST                string
	SCGST                string
	SSGST                string
	OfferPercentagePrice string
	OfferFlatPrice       string
	Offer                string
	MRP                  string `json:"mrp"`

	ClassesType              string                         `json:"ClassesType"`
	ClassName                string                         `json:"ClassName"`
	ClassDesc                string                         `json:"ClassDesc"`
	ClassStartDate           string                         `json:"ClassStartDate"`
	ClassEndDate             string                         `json:"ClassEndDate"`
	ClassClosingDate         string                         `json:"ClassClosingDate"`
	ClassVideoLink           string                         `json:"ClassVideoLink"`
	ClassTrainerName         string                         `json:"ClassTrainerName"`
	ClassTrainerDesc         string                         `json:"ClassTrainerDesc"`
	ClassAddressType         string                         `json:"ClassAddressType"`
	ClassFullAddress         string                         `json:"ClassFullAddress"`
	ClassSlots               string                         `json:"ClassSlots"`
	ClassImages              string                         `json:"ClassImages"`
	ClassTrainerDetails      string                         `json:"ClassTrainerDetails"`
	ClassSlotsArray          []SlotArr                      `json:"SlotArr"`
	ClassImagesArray         []ClassImageArr                `json:"ClassImageArr"`
	ClassTrainerDetailsArray []ClassTClassTrainerDetailsArr `json:"ClassTClassTrainerDetailsArr"`
	ClassesProductsRes       string                         `json:"ClassesProductsRes"`
}

type SlotArr struct {
	StartTime    string `json:"StartTime"`
	EndTime      string `json:"EndTime"`
	Fees         string `json:Fees`
	FeesPerClass string `json:FeesPerClass`
	Duration     string `json:Duration`
	MaxStudents  string `json:MaxStudents`
}

type ClassImageArr struct {
	Photo       string `json:"Photo"`
	base64Image string `json:"base64Image"`
}

type ClassTClassTrainerDetailsArr struct {
	TrainerName string `json:"TrainerName"`
	TrainerDesc string `json:"TrainerDesc"`
}

type UnitArr struct {
	Unit  string `json:"unit"`
	Price string `json:"price"`
}

type ProductSubCategoryRes struct {
	Id                   string `json:"Id"`
	SubCategory          string `json:"SubCategory"`
	Category             string `json:"Category"`
	Name                 string `json:"Name"`
	ImageURL             string `json:"ImageUrl"`
	TotalPrice           string `json:"TotalPrice"`
	UnitArray            string `json:"UnitArray"`
	Desc                 string `json:"Desc"`
	Unit                 string `json:"Unit"`
	Price                string `json:"Price"`
	Added                string `json:"Added"`
	Qty                  int    `json:"Qty"`
	OpeningStock         string `json:"openingstock"`
	FoodType             string `json:"FoodType"`
	PCGST                string
	PSGST                string
	SCGST                string
	SSGST                string
	OfferPercentagePrice string
	OfferFlatPrice       string
	Offer                string
	MRP                  string `json:"mrp"`

	ClassesType    string `json:"ClassesType"`
	ClassName      string `json:"ClassName"`
	ClassDesc      string `json:"ClassDesc"`
	ClassStartDate string `json:"ClassStartDate"`
	ClassEndDate   string `json:"ClassEndDate"`
	SlotStartTime  string `json:"SlotStartTime"`
	SlotEndTime    string `json:"SlotEndTime"`

	ClassClosingDate string `json:"ClassClosingDate"`
	ClassVideoLink   string `json:"ClassVideoLink"`
	ClassTrainerName string `json:"ClassTrainerName"`
	ClassTrainerDesc string `json:"ClassTrainerDesc"`
	ClassAddressType string `json:"ClassAddressType"`

	StudentName       string `json:"StudentName"`
	StudentFatherName string `json:"StudentFatherName"`
	StudentDOB        string `json:"StudentDOB"`
	StudnetGender     string `json:"StudnetGender"`
	StudentAge        string `json:"StudentAge"`
	StudentDesc       string `json:"StudentDesc"`
	StudentAddress    string `json:"StudentAddress"`
	StudentMobile     string `json:"StudentMobile"`
	StudnetEmail      string `json:"StudnetEmail"`

	ClassFullAddress         string `json:"ClassFullAddress"`
	ClassSlots               string `json:"ClassSlots"`
	ClassTrainerDetails      string `json:"ClassTrainerDetails"`
	ClassImages              string `json:"ClassImages"`
	ClassSlotsArray          string `json:"ClassSlotsArray"`
	ClassImagesArray         string `json:"ClassImagesArray"`
	ClassTrainerDetailsArray string `json:"ClassTClassTrainerDetailsArr"`
	ClassesProductsRes       string `json:"ClassesProductsRes"`
}
type CategoryRes struct {
	Id       string `json:"retailer_id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	CoverUrl string `json:"cover_url"`
}

type Request struct {
	Id           string
	Status       string
	Status2      string
	Remarks      string
	ClusterName  string
	CategoryName string
	SubCategory  string
	Terms        string

	Name              string
	StoreName         string
	Txnid             string
	Category          string
	StoreId           string
	AddressId         string
	FirstName         string
	LastName          string
	Email             string
	Mobile            string `json:Mobile",omitempty"`
	DeviceID          string
	SupplyOrdertxn_id string
	ProductPackageQty string
	ClusterId         string
	ProductId         string
	P2P2Cid           string
	ProfilePicture    string
	DOB               string
	Password          string
	NewPassword       string
	Post              string
	District          string
	PIN               string
	State             string
	Province          string
	City              string
	Nationality       string
	ChannelName       string
	ChannelKey        string
	OS                string
	OSVersion         string
	APPVersion        string
	IMEI              string
	IMSI              string
	PIP               string
	ServerDateTime    string
	SessionID         string
	SessionStatus     string
	Locality          string
	HouseNo           string
	Keyword           string
	Address           string
	Landmark          string
	Latitude          float64
	Longitude         float64
	From              string
	To                string
	Cart              string
	OrderNumber       string
	ContactMobile     string
	ContactEmail      string
	DeliveryAddress   string
	OTP               string
	StatusJson        string
	Description       string
	FCMToken          string

	ConfirmPassword string
	UserType        string

	Amount                 string
	DeliveryCharge         string
	ErrorMsg               string
	PaymentId              string
	Signature              string
	Nounce                 string
	RequestID              string
	TermsAndCondition      string
	ProductID              string
	No_packets             string
	Payment_type           string
	Order_type             string
	ProductName            string
	ProductQty             string
	ProductQualityType     string
	ShelLife               string
	Product_QualityMeasure string
	PackageType            string
	PackageQty             string
	Krishivala_BuyingPrice string
	ClusterBuyingPrice     string

	ProductSKU                 string
	ClusterSellingPrice        string
	Krishivala_CGST_charge     string
	Krishivala_SGST_charge     string
	Krishivala_CESS_charge     string
	Krishivala_OtherTAX_charge string
	ClusterSupplyCharge        string
	ClusterPackagingCharge     string
	ClusterBankCharge          string
	ClusterDiscount            string
	DeliveryTimngs             string
	DeliveryDate               string

	Address1             string
	Address2             string
	Fromdate             string
	Todate               string
	CustomerMobile       string
	Repeatorder          string
	Locationinfo         string
	ProducerSellingPrice string
	ProducerUnits        string
	Startdate            string
	Enddate              string
	ProducerId           string
	ProducerName         string
	ProducerMobile       string
	DeliveryDiscount     string

	ProducerCGST     string
	ProducerSGST     string
	ProducerCESS     string
	ProducerOtherTax string
	// Address1 string
	// Address1 string
	// Address1 string
	// Address1 string

	PayUTxnId string
	PayUTax   string

	PayU_Key     string
	PayU_Txnid   string
	PayU_Amount  string
	PayU_Product string
	PayU_Name    string
	PayU_Email   string
	PayU_Mobile  string
	PayU_Surl    string
	PayU_Furl    string
	PayU_Hash    string
	PayU_Salt    string
}
