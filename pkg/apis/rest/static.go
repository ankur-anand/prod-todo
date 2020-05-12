package rest

import (
	"net/http"

	"github.com/ankur-anand/prod-app/pkg/logger"
)

var (
	homeRouteStaticResponse = []byte(`{"message": "hello world from rest svc"}`)
	healthiest              = []byte(`{"alive": true}`)
)

type staticHandler struct {
	logger *logger.Logger
}

func newStaticHandler(logger *logger.Logger) staticHandler {
	sh := staticHandler{logger: logger}
	return sh
}

func (sh staticHandler) home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(homeRouteStaticResponse)
	if err != nil {
		sh.logger.Error("response writer err", logger.Error(err))
	}
	sh.logger.Info("homepage", httpReqField(http.StatusOK, r, nil)...)
}

func (sh staticHandler) healthLive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(healthiest)
	if err != nil {
		sh.logger.Error("response writer err", logger.Error(err))
	}
	sh.logger.Info("healthlive", httpReqField(http.StatusOK, r, nil)...)
}
