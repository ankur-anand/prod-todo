package resthandler

import (
	"net/http"

	"go.uber.org/zap"
)

var (
	homeRouteStaticResponse = getRespMsg("hello world from todo rest svc")
	healthiest              = getJSONResp(`{"alive": true}`)
)

type staticHandler struct {
	logger *zap.Logger
}

func newStaticHandler(logger *zap.Logger) staticHandler {
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
