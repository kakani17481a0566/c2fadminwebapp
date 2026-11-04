package BootConfig

type Root struct {
	Tksp Tksp `json:"Tksp,omitempty"`
}

type Tksp struct {
	Config Config `json:"Config,omitempty"`
}

type Config struct {
	Database []*Database `json:"Database,omitempty"`
	Logging  *Logging    `json:"Logging,omitempty"`
	Queues   []*Queues   `json:"Queues,omitempty"`
	Server   *Server     `json:"Server,omitempty"`
	State    []*State    `json:"State,omitempty"`
	Partner  *Partner    `json:"Partner,omitempty"`
	Keys     *Keys       `json:"Keys,omitempty"`
	Misc     *Misc       `json:"Misc,omitempty"`
	Param    []*Param    `json:"Param,omitempty"`
}

type Param struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

type Keys struct {
	Dek string `json:"Dek,omitempty"`
	Mek string `json:"Mek,omitempty"`
	Mdk string `json:"Mdk,omitempty"`
}
type Misc struct {
	EnableDuplicateCheck   bool `json:"EnableDuplicateCheck,omitempty"`
	DuplicateCheckTime     int  `json:"DuplicateCheckTime,omitempty"`
	CashInLimitCheck       bool `json:"CashInLimitCheck,omitempty"`
	CashOutLimitCheck      bool `json:"CashOutLimitCheck,omitempty"`
	FundTransferLimitCheck bool `json:"FundTransferLimitCheck,omitempty"`
}

type Partner struct {
	Rule                string           `json:"Rule,omitempty"`
	RetryTime           int              `json:"RetryTime,omitempty"`
	EableZesaSwitchOver string           `json:"EableZesaSwitchOver,omitempty"`
	PartnerConfig       []*PartnerConfig `json:"PartnerConfig,omitempty"`
}

type PartnerConfig struct {
	ServiceName string `json:"ServiceName,omitempty"`
	TimeOut     string `json:"TimeOut,omitempty"`
	URL         string `json:"URL,omitempty"`
	SSLKey      string `json:"SSLKey,omitempty"`
	CSRKey      string `json:"CSRKey,omitempty"`
	Request     string `json:"Request,omitempty"`
	Username    string `json:"Username,omitempty"`
	Password    string `json:"Password,omitempty"`
	IsHttps     string `json:"IsHttps,omitempty"`
	Method      string `json:"Method,omitempty"`
	Sender      string `json:"Sender,omitempty"`
}

type State struct {
	ServiceName string    `json:"ServiceName,omitempty"`
	States      []*States `json:"States,omitempty"`
}

type States struct {
	Sequence             string
	Name                 string
	Module               string
	ExpectedResponceCode string
	Success              string
	Failure              string
}

type Server struct {
	ApplicationName string      `json:"ApplicationName,omitempty"`
	HTTPConfig      *HTTPConfig `json:"HTTPConfig,omitempty"`
	IsHTTP          string      `json:"IsHTTP,omitempty"`
	IsTCP           string      `json:"IsTCP,omitempty"`
	RunMode         string      `json:"RunMode,omitempty"`
	TCPConfig       *TCPConfig  `json:"TCPConfig,omitempty"`
}

type TCPConfig struct {
	IP         string `json:"IP,omitempty"`
	IsTLS      string `json:"IsTLS,omitempty"`
	Port       string `json:"Port,omitempty"`
	PrivateKey string `json:"PrivateKey,omitempty"`
	PublicKey  string `json:"PublicKey,omitempty"`
}

type HTTPConfig struct {
	HTTPPort          string  `json:"HTTPPort,omitempty"`
	HTTPSPort         string  `json:"HTTPSPort,omitempty"`
	IsHTTPS           string  `json:"IsHTTPS,omitempty"`
	IsSession         string  `json:"IsSession,omitempty"`
	PrivateKey        string  `json:"PrivateKey,omitempty"`
	PublicKey         string  `json:"PublicKey,omitempty"`
	SessionExpiryTime string  `json:"SessionExpiryTime,omitempty"`
	URLs              []*URLs `json:"URLs,omitempty"`
}

type URLs struct {
	ServiceName string `json:"ServiceName,omitempty"`
	URL         string `json:"URL,omitempty"`
}
type Queues struct {
	Instance string `json:"Instance,omitempty"`
	Module   string `json:"Module,omitempty"`
	Queue    *Queue `json:"Queue,omitempty"`
}

type Queue struct {
	Autoack      string `json:"Autoack,omitempty"`
	Binding      string `json:"Binding,omitempty"`
	Ctag         string `json:"Ctag,omitempty"`
	Durable      string `json:"Durable,omitempty"`
	Exchange     string `json:"Exchange,omitempty"`
	ExchangeType string `json:"ExchangeType,omitempty"`
	ExpiryTime   string `json:"ExpiryTime,omitempty"`
	IP           string `json:"Ip,omitempty"`
	Name         string `json:"Name,omitempty"`
	Nowait       string `json:"Nowait,omitempty"`
	Password     string `json:"Password,omitempty"`
	Port         string `json:"Port,omitempty"`
	Username     string `json:"Username,omitempty"`
	Vhost        string `json:"Vhost,omitempty"`
}

type Database struct {
	IP       string `json:"IP,omitempty"`
	Name     string `json:"Name,omitempty"`
	Password string `json:"Password,omitempty"`
	Port     string `json:"Port,omitempty"`
	Type     string `json:"Type,omitempty"`
	Username string `json:"Username,omitempty"`
}

type Logging struct {
	Database *Database `json:"Database,omitempty"`
	File     *File     `json:"File,omitempty"`
	Level    string    `json:"Level,omitempty"`
	Server   *Server   `json:"Server,omitempty"`
	Type     string    `json:"Type,omitempty"`
}

type File struct {
	Name string `json:"Name,omitempty"`
	Size string `json:"Size,omitempty"`
}
