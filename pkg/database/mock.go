package database

//Mock - struct to test in the application
//this mock implements the StoreInterface
type Mock struct {
	Err    error
	Amount string
}

func (m *Mock) InsertNfeAmount(accessKey, amount string) error {
	return m.Err
}

func (m *Mock) GetNfeAmount(accessKey string) (string, error) {
	return m.Amount, m.Err
}
