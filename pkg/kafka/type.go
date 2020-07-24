package kafka

//WorkerMessage structure of message in topic worker
type WorkerMessage struct {
	AccessKey string `json:"access_key"`
	XML       string `json:"xml"`
}

//DBMessage structure of message in topic db
type DBMessage struct {
	AccessKey string `json:"access_key"`
	Amount    string `json:"amount"`
}
