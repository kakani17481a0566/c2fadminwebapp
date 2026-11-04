package InfiniApi

type Root struct {
	Status  string `json:"status,omitempty"`
	Data    []Data `json:"data, omitempty"`
	Message string `json:"Message, omitempty"`
}

type Data struct {
	Id        string `json:"id,omitempty"`
	Customid  string `json:"customid,omitempty"`
	Customid1 string `json:"customid1,omitempty"`
	Customid2 string `json:"customid2,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Status    string `json:"status,omitempty"`
}

//////////////////////////////////////////////

//type Root struct {
//	Status  string `json:"status,omitempty"`
//	Data    Data   `json:"data, omitempty"`
//	Message string `json:"Message, omitempty"`
//}

//type Sub struct {
//	Group_id String `json:"group_id,omitempty"`
//	Entity   map[string]Sub
//}

//type Sub struct {
//	Id        string `json:"id,omitempty"`
//	Customid  string `json:"customid,omitempty"`
//	Customid1 string `json:"customid1,omitempty"`
//	Customid2 string `json:"customid2,omitempty"`
//	Mobile    string `json:"mobile,omitempty"`
//	Status    string `json:"status,omitempty"`
//}
