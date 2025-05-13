package euvd

import "github.com/sirupsen/logrus"

const (
	baseURL = "https://euvdservices.enisa.europa.eu/api"
)

type EUVD struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *EUVD {
	if logger == nil {
		logger = logrus.New()
	}

	return &EUVD{
		logger: logger,
	}
}
