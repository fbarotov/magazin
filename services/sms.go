package services

const (
	// TODO: these keys should be either stored in a file outside codebase
	// or be fetched from environment variable
	smsSerViceKey1 = "key1"
	smsServiceKey2 = "key2"
)

var _ Service = &SMS{}

// SMS implements Service interface.
type SMS struct {}

func (s *SMS) Send(phone, message string) error {
	// send http request to the SMS service
	// here, we assume success
	return nil
}