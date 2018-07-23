package models

type OperationCore struct {
	BaseModel
	SystemID      string `json:"systemId"`
	ProcessID     string `json:"processId"`
	Name          string `json:"name"`
	EventIn       string `json:"event_in"`
	EventOut      string `json:"event_out"`
	Image         string `json:"image"`
	Commit        bool   `json:"commit"`
	Reprocessable bool   `json:"reprocessable"`
	Version       string `json:"version"`
}

func NewOperationCore() *OperationCore {
	op := new(OperationCore)
	op.Metadata = Metadata{
		Type: "operation",
	}
	return op
}
