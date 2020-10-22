package types

type ResponseJSON struct {
	Data        DataJSON `json:"data,omitempty"`
	Action      string   `json:"action,omitempty"`
	ResponseID  string   `json:"responseID,omitempty"`
	ExecuteTime string   `json:"executeTime,omitempty"`
}

type TempResponseJSON struct {
	*ResponseJSON
	Data string `json:"data,omitempty"`
}

type DataJSON struct {
	DeedSeq      bool     `json:"needSeq,omitempty"`
	Seq          int64    `json:"seq,omitempty"`
	Status       string   `json:"status,omitempty"`
	Result       string   `json:"result,omitempty"`
	IsInsnLimit  bool     `json:"isInsnLimit,omitempty"`
	TotalGas     int64    `json:"totalGas,omitempty"`
	ExecutionGas int64    `json:"executionGas,omitempty"`
	ExtraGas     int64    `json:"extraGas,omitempty"`
	NodeIDs      []string `json:"nodeIDs,omitempty"`
	ExecuteTime  []string `json:"executeTime,omitempty"`
}

type RequestJSON struct {
	Action     string `json:"action,omitempty"`
	ContractID string `json:"contractID,omitempty"`
	Operation  string `json:"operation,omitempty"`
	Arg        string `json:"arg,omitempty"`
}
