package api

type Controller interface {
	Client() *Client
	Name() string
}

// Response structs
type addResp struct {
	Result      string                 `json:"result"`
	UUID        string                 `json:"uuid"`
	Validations map[string]interface{} `json:"validations,omitempty"`
}

type deleteResp struct {
	Result string `json:"result"`
}
