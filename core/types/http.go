package types

type ResponseJSON struct {
	Data *DataJSON     `json:"data,omitempty"`
	Action string      `json:"action,omitempty"`
	ResponseID string  `json:"responseID,omitempty"`
	ExecuteTime string `json:"executeTime,omitempty"`
}

type DataJSON struct {
	DeedSeq bool `json:"needSeq,omitempty"`
	Seq int64 `json:"seq,omitempty"`
	Status string `json:"status,omitempty"`
	Result interface{} `json:"result,omitempty"`
}

type RequestJSON struct {
	Action string `json:"action,omitempty"`
	ContractID string `json:"contractID,omitempty"`
	Operation string `json:"operation,omitempty"`
	Arg string `json:"arg,omitempty"`
}