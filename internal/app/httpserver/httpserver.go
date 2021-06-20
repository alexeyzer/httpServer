package httpserver

import "github.com/sirupsen/logrus"

// httpServer
type httpserver struct {
	config *config
	logger *logrus.Logger
}

func NewHttpServer(config *config) *httpserver {
	server := httpserver{
		config: config,
		logger: logrus.New(),
	}
	return &server
}

func (S *httpserver) Start() error{
	err := S.ConfigLogger()
	if err != nil {
		return err
	}
	S.logger.Info("Staring httpserver")
	return nil
}

func (S *httpserver) ConfigLogger() error{
	level, err := logrus.ParseLevel(S.config.LogLevel)
	if err != nil {
		return err
	}
	S.logger.SetLevel(level)
	return nil
}