package services

const (
	// TODO: these keys should be either stored in a file outside codebase
	// or be fetched from environment variable
	emailSerViceKey1 = "key1"
	emailServiceKey2 = "key2"
)

var _ Service = &Email{}

// Email implements Service interface.
type Email struct {}

func (e *Email) Send(email, message string) error {
	// make https request to the third-party service
	return nil
}
