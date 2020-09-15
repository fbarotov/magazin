package services

type Service interface {
	Send(destination, message string) error
}