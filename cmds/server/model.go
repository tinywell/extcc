package server

// ReqChaincode ...
type ReqChaincode struct {
	Name string `json:"name,omitempty"`
	CCID string `json:"ccid,omitempty"`
	Addr string `json:"addr,omitempty"`
	TLS  string `json:"tls,omitempty"` //TODO:
}

// ReqTransaction ...
type ReqTransaction struct {
	Chaincode    string            `json:"chaincode,omitempty"`
	Func         string            `json:"func,omitempty"`
	Args         []string          `json:"args,omitempty"`
	MSPID        string            `json:"mspid,omitempty"`
	MSPCert      string            `json:"msp_cert,omitempty"`
	TransientMap map[string][]byte `json:"transient_map,omitempty"`
}

// Response ...
type Response struct {
	Code    int         `json:"code" `
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message" `
}

// RspTx ...
type RspTx struct {
	Status  int    `json:"status" `
	Payload []byte `json:"payload" `
	Message string `json:"message" `
}
