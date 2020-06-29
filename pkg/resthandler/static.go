package resthandler

import (
	"net/http"

	"github.com/ankur-anand/prod-todo/pkg/logger"
)

var (
	homeRouteStaticResponse = []byte(`{"message": "hello world from resthandler svc"}`)
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
	checkResponseWriteErr(err, sh.logger)
	sh.logger.Info("homepage", httpReqField(http.StatusOK, r, nil)...)
}

func (sh staticHandler) healthLive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(healthiest)
	checkResponseWriteErr(err, sh.logger)
	sh.logger.Info("healthlive", httpReqField(http.StatusOK, r, nil)...)
}
