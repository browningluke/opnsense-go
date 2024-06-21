package api

type ReqOpts struct {
	AddEndpoint         string
	GetEndpoint         string
	UpdateEndpoint      string
	DeleteEndpoint      string
	ReconfigureEndpoint string
	SearchEndpoint      string

	Monad string
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

type searchReq struct {
	Current      int                    `json:"current"`
	RowCount     int                    `json:"rowCount"`
	Sort         map[string]interface{} `json:"sort"`
	SearchPhrase string                 `json:"searchPhrase"`
}
