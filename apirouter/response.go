package apirouter

type IndexResponse struct {
	Title  string `json:"title"`
	Header string `json:"header"`
	Block1 string `json:"block1"`
	Block2 string `json:"block2"`
}

type PingResponse struct {
	Message string `json:"message"`
}
