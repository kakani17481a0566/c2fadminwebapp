package Rule

type Root struct {
	Rule []*Rule `json:"Rule,omitempty"`
}

type Rule struct {
	Name      string `json:"Name,omitempty"`
	Entity    string `json:"Entity,omitempty"`
	Type      string `json:"Type,omitempty"`
	Condition string `json:"Condition,omitempty"`
	Success   string `json:"Success,omitempty"`
	Failure   string `json:"Failure,omitempty"`
	IF        []*IF  `json:"IF,omitempty"`
}

type IF struct {
	Id  string `json:"Id,omitempty"`
	LHS string `json:"LHS,omitempty"`
	OP  string `json:"OP,omitempty"`
	RHS string `json:"RHS,omitempty"`
}
