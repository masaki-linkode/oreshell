package infra

import "os"

type OSService interface {
	Getenv(name string) string
	Setenv(name string, val string) error
}

type MyOSService struct {
}

func (me MyOSService) Getenv(name string) string {
	return os.Getenv(name)
}

func (me MyOSService) Setenv(name string, val string) error {
	return os.Setenv(name, val)
}
