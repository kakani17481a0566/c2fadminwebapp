package AfricomApi

type BalanceEnquiryResp struct {
	BalanceHeader  *BalanceHeader   `json:"BalanceHeader,omitempty"`
	BalanceRecords []*BalanceRecord `json:"BalanceRecords,omitempty"`
}

//type BalanceRecords struct {
//	BalanceRecord []*BalanceRecords `json:"BalanceRecord,omitempty"`
//}

type BalanceRecord struct {
	AccountType    string `json:"AccountType,omitempty"`
	Description    string `json:"Description,omitempty"`
	ExpirationTime string `json:"ExpirationTime,omitempty"`
	Balance        string `json:"Balance,omitempty"`
}

type BalanceHeader struct {
	ResultCode            string `json:"ResultCode,omitempty"`
	Description           string `json:"Description,omitempty"`
	SubscriberNumber      string `json:"SubscriberNumber,omitempty"`
	DealerReferenceNumber string `json:"DealerReferenceNumber,omitempty"`
}

type RechargeResp struct {
	TransferHeader   *TransferHeader   `json:"TransferHeader,omitempty"`
	TransfereeRecord *TransfereeRecord `json:"TransfereeRecord,omitempty"`
	TranferorRecord  *TranferorRecord  `json:"TranferorRecord,omitempty"`
}
type TransferHeader struct {
	ResultCode     string  `json:"ResultCode,omitempty"`
	Description    string  `json:"Description,omitempty"`
	TransferAmount float32 `json:"TransferAmount,omitempty"`
}
type TransfereeRecord struct {
	TransfereeNumber         string  `json:"TransfereeNumber,omitempty"`
	TransfereeOpenBalance    float32 `json:"TransfereeOpenBalance,omitempty"`
	TransfereeClosingBalance float32 `json:"TransfereeClosingBalance,omitempty"`
}
type TranferorRecord struct {
	TranferorNumber          string  `json:"TranferorNumber,omitempty"`
	TransferorOpenBalance    float32 `json:"TransferorOpenBalance,omitempty"`
	TransferorClosingBalance float32 `json:"TransferorClosingBalance,omitempty"`
}
type ProductSubscribeResp struct {
	SubscribeHeader *SubscribeHeader `json:"SubscribeHeader,omitempty"`
	SubscribeRecord *SubscribeRecord `json:"SubscribeRecord,omitempty"`
}
type SubscribeHeader struct {
	ResultCode  string `json:"ResultCode,omitempty"`
	Description string `json:"Description,omitempty"`
}
type SubscribeRecord struct {
	SubscriberNumber string  `json:"SubscriberNumber,omitempty"`
	ProductName      string  `json:"ProductName,omitempty"`
	ProductId        int     `json:"ProductId,omitempty"`
	Price            float32 `json:"Price,omitempty"`
	EffectiveDate    string  `json:"EffectiveDate,omitempty"`
	ExpireDate       string  `json:"ExpireDateomitempty"`
}

/*type Response struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
	Status  string `json:"Status,omitempty"`
	//Parameter []*Parameter `json:"Parameter,omitempty"`
}
*/
