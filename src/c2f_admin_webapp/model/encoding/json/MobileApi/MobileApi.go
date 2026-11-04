package MobileApi

type Response struct {
	Code                     string
	Message                  string
	Status                   string
	StoreStatus              string
	Name                     string
	Email                    string
	Mobile                   string
	Products                 []*ProductArr
	ActiveProducts           []*ProductArr
	FoodProducts             []*FoodProducts
	Categories               []*CategoriesRes
	SubCategories            []*CategoriesRes
	ClassesTraingTypesRes    []*ClassesTraingTypesRes
	ClassesProductsRes       string
	ProductArr               *ProductArr
	CityList                 []*City
	StateList                []*State
	Category                 string
	Suggestion               []string
	Orders                   []*Orders
	Dob                      string
	SessionID                string
	Gender                   string
	Post                     string
	PIN                      string
	District                 string
	State                    string
	Nationlaity              string
	Locality                 string
	HouseNo                  string
	Landmark                 string
	StoreImageUrl            string
	Cart                     string
	OrderNumber              string
	Address                  string
	RetailerName             string
	Radius                   string
	MinOrder                 string
	StoreImage               string
	DeliveryTime             string
	StoreName                string
	TotalAmount              string
	SubUsers                 []SubUser
	Role                     string
	Roles                    []string
	VersionNumber            string
	ForceUpgrade             string
	ProductImageUrl          string
	OrderDate                string
	Redirect                 string
	DeliveryOnline           string
	Lon                      string
	Lat                      string
	OpenTime                 string
	CloseTime                string
	AccNumber                string
	IFSC                     string
	AccName                  string
	BankName                 string
	PAN                      string
	GSTIN                    string
	CIN                      string
	FSSAI                    string
	AcceptOnline             string
	City                     string
	Description              string
	Taxes                    string
	DeliveryType             string
	PickupTime               string
	CustomerCareMobile1      string
	CustomerCareEmail1       string
	NotificationMsg          string
	Slot                     string
	SlotDays                 string
	BusinessDays             string
	StockManagement          string
	StockManagementThreshold string
}

type CategoriesRes struct {
	Name        string
	ImageUrl    string
	CoverUrl    string
	Active      string
	SubCategory []*SubCategory
	IsSelect    bool
}

type ClassesTraingTypesRes struct {
	Name        string
	Description string
}

type SubCategory struct {
	Name string
}
type City struct {
	Name string
}

type State struct {
	Name string
}

type Orders struct {
	OrderNumber    string
	Cart           *Cart
	OrderTimeStamp string
	Payment        string
	CurrentStatus  string
	Status         *Status
	TotalAmount    string
	ConsumerMobile string
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
	DeliveryType    string

	//classes
	ClassesType    string
	ClassName      string
	ClassDesc      string
	ClassStartDate string
	ClassEndDate   string

	SlotStartTime string
	SlotEndTime   string

	ClassClosingDate  string
	ClassVideoLink    string
	ClassTrainerName  string
	ClassTrainerDesc  string
	ClassAddressType  string
	ClassFullAddress  string
	StudentName       string
	StudentFatherName string
	StudentDOB        string
	StudentGender     string
	StudentAge        string
	StudentDesc       string
	StudentAddress    string
	StudentMobile     string
	StudentEmail      string

	//SelectedClassSlots          string
	//SelectedClassTrainerDetails *ClassTClassTrainerDetailsArr

	//	ClassImages                 string
	//	ClassTrainerDetailsArray    []*ClassTClassTrainerDetailsArr
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

type SubUser struct {
	Username string
	Name     string
	Role     string
	Email    string
}

type ProductSubCategory struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Title    string `json:"Title"`
	Products []*ProductRes
}
type FoodProducts struct {
	Name        string `json:"Name"`
	SubCategory []*SubCategoryRes
}

type SubCategoryRes struct {
	Name     string `json:"Name"`
	Products []*ProductArr
}

type ClassProducts struct {
	Name     string `json:"Name"`
	Products []*ProductArr
}

type Taxes struct {
	Name   string
	Value  string
	Type   string
	Amount string
}

type ProductArr struct {
	Id                   string    `json:"Id"`
	Category             string    `json:"Category"`
	SubCategory          string    `json:"SubCategory"`
	MRP                  string    `json:"MRP"`
	SKU                  string    `json:"SKU"`
	Brand                string    `json:"Brand"`
	BarCode              string    `json:"BarCode"`
	BatchNo              string    `json:"BatchNo"`
	ExpiryDate           string    `json:"ExpiryDate"`
	ManufactureDate      string    `json:"ManufactureDate"`
	OpeningStock         string    `json:"OpeningStock"`
	Name                 string    `json:"Name"`
	ImageURL             string    `json:"ImageUrl"`
	PurchasePrice        string    `json:"PurchasePrice"`
	Desc                 string    `json:"Desc"`
	Unit                 string    `json:"Unit"`
	Price                string    `json:"Price"`
	SelectedUnit         string    `json:"SelectedUnit"`
	Qty                  int       `json:"Qty"`
	GST                  string    `json:"GST"`
	TotalPrice           string    `json:"TotalPrice"`
	Active               string    `json:"Active"`
	Filter               string    `json:"Filter"`
	HSNCode              string    `json:"HsnCode"`
	PurchaseGST          string    `json:"PurchaseGST"`
	SellingGST           string    `json:"SellingGST"`
	Added                string    `json:"Added"`
	UnitArray            []UnitArr `json:"UnitArr"`
	FoodType             string    `json:"FoodType"`
	PCGST                string
	PSGST                string
	SCGST                string
	SSGST                string
	OfferPercentagePrice string
	OfferFlatPrice       string
	Offer                string
	NonEAN               string

	//Classes
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
	ClassTrainerDetails      string                         `json:"ClassTrainerDetails"`
	ClassImages              string                         `json:"ClassImages"`
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

type ProductRes struct {
	Id                   string
	MRP                  string `json:"MRP"`
	SKU                  string `json:"SKU"`
	Brand                string `json:"Brand"`
	BarCode              string `json:"BarCode"`
	BatchNo              string `json:"BatchNo"`
	ExpiryDate           string `json:"ExpiryDate"`
	ManufactureDate      string `json:"ManufactureDate"`
	OpeningStock         string `json:"OpeningStock"`
	Name                 string `json:"Name"`
	ImageURL             string `json:"ImageUrl"`
	PurchasePrice        string `json:"PurchasePrice"`
	Desc                 string `json:"Desc"`
	Unit                 string `json:"Unit"`
	Price                string `json:"Price"`
	Qty                  int    `json:"Qty"`
	GST                  string `json:"GST"`
	TotalPrice           string `json:"TotalPrice"`
	Active               string `json:"Active"`
	Filter               string `json:"Filter"`
	HSNCode              string `json:"HsnCode"`
	PurchaseGST          string `json:"PurchaseGST"`
	SellingGST           string `json:"SellingGST"`
	Added                string `json:"Added"`
	PCGST                string
	PSGST                string
	SCGST                string
	SSGST                string
	SGSTAmt              string
	CGSTAmt              string
	BasicPrice           string
	OfferPercentagePrice string
	OfferFlatPrice       string
	Offer                string

	//Classes
	ClassesType         string                         `json:"ClassesType"`
	ClassName           string                         `json:"ClassName"`
	ClassDesc           string                         `json:"ClassDesc"`
	ClassStartDate      string                         `json:"ClassStartDate"`
	ClassEndDate        string                         `json:"ClassStartDate"`
	SlotStartTime       string                         `json:"SlotStartTime"`
	SlotEndTime         string                         `json:"SlotEndTime"`
	ClassClosingDate    string                         `json:"ClassClosingDate"`
	ClassVideoLink      string                         `json:"ClassVideoLink"`
	ClassTrainerName    string                         `json:"ClassTrainerName"`
	ClassTrainerDesc    string                         `json:"ClassTrainerDesc"`
	ClassAddressType    string                         `json:"ClassAddressType"`
	ClassFullAddress    string                         `json:"ClassFullAddress"`
	ClassSlots          []SlotArr                      `json:"SlotArr"`
	ClassImages         []ClassImageArr                `json:"ClassImageArr"`
	ClassTrainerDetails []ClassTClassTrainerDetailsArr `json:"ClassTrainerDetails"`
}
type Request struct {
	Name                     string
	StoreName                string
	StoreImage               string
	Radius                   string
	DeliveryTime             string
	MinOrder                 string
	Category                 string
	StoreId                  string
	FirstName                string
	Email                    string
	Mobile                   string
	DOB                      string
	Password                 string
	NewPassword              string
	Post                     string
	District                 string
	PIN                      string
	State                    string
	Province                 string
	City                     string
	ProductsData             string
	Nationality              string
	ChannelName              string
	ChannelKey               string
	ToDate                   string
	FromDate                 string
	OS                       string
	OSVersion                string
	APPVersion               string
	IMEI                     string
	IMSI                     string
	PIP                      string
	ServerDateTime           string
	SessionID                string
	SessionStatus            string
	Locality                 string
	Address                  string
	Landmark                 string
	Latitude                 float64
	Longitude                float64
	OrderNumber              string
	PaymentStatus            string
	OrderStatus              string
	Subscribe                string
	ConsumerId               string
	Cart                     string
	Order                    string
	Keyword                  string
	ProductStatus            string
	SuperCategory            string
	ProductImage             string
	Role                     string
	NewCategory              bool
	NewSubCategory           bool
	SubUserMobile            string
	Description              string
	ZIP                      string
	CommAddress1             string
	CommAddress2             string
	PAN                      string
	GSTIN                    string
	CIN                      string
	FSSAI                    string
	DeliveryOnline           string
	AcceptOnline             string
	AccountNumber            string
	BankName                 string
	IFSCCode                 string
	BranchName               string
	AccountNo                string
	AccountName              string
	Photo                    string
	Status                   string
	OpenTime                 string
	CloseTime                string
	Location                 string
	Taxes                    string
	DeliveryType             string
	PickupTime               string
	CustomerCareMobile1      string
	CustomerCareEmail1       string
	SubCategory              string
	OTP                      string
	Slot                     string
	SlotDays                 string
	BusinessDays             string
	StoreUsername            string
	Productid                string
	StockManagement          string
	StockManagementThreshold string
	Products                 []Products
	ClassImages              []ClassImages
}

type Products struct {
	MRP                  string
	BarCode              string
	PurchasePrice        string
	Price                string
	Qty                  string
	HSNCode              string
	SCGST                string
	SSGST                string
	OfferPercentagePrice string
	OfferFlatPrice       string
	Offer                string
}

type ClassImages struct {
	Photo string
}

type RetailerApiResponse struct {
	Code          string
	Message       string
	Status        string
	SessionID     string
	StoreUsername string
}
