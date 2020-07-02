package resthandler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type contextKey string

var (
	contextKeyDuration = contextKey("duration")

	someThingWentWrong = []byte(`{"status": "ERROR", 
	"error": {"message": "Something went wrong.", "kind": "ServerException"}`)
)

// MuxHandler is a Handler that responds to an HTTP request.
type MuxHandler struct {
	log           *zap.Logger
	regAndAuth    auth
	staticHandler staticHandler
	router        *mux.Router
}

// NewMuxHandler returns an initialized http.Handler
func NewMuxHandler(logger *zap.Logger) *MuxHandler {
	mh := MuxHandler{
		staticHandler: newStaticHandler(logger),
		log:           logger,
		router:        mux.NewRouter(),
	}
	mh.initializeRoutes()
	return &mh
}

// ServeHTTP responds to an HTTP request
func (mh *MuxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// if the content-type is not application json reject the request upfront
	h := r.Header.Get("Content-Type")
	if !strings.Contains(h, "application/json") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// put start time in the context
	r = r.WithContext(context.WithValue(r.Context(), contextKeyDuration, time.Now()))
	// all are json response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// prevent mime sniff
	w.Header().Set("X-Content-Type-Options", "nosniff")

	mh.router.ServeHTTP(w, r)
}

func (mh *MuxHandler) initializeRoutes() {
	// home route
	mh.router.HandleFunc("/", mh.staticHandler.home)

	// health check for liveness Probe
	mh.router.HandleFunc("/health/live", mh.staticHandler.healthLive)

	// health check for readinessProbe
	mh.router.HandleFunc("/health/ready", func(writer http.ResponseWriter, request *http.Request) {

	})

	// login and registration
	mh.router.HandleFunc("/v1/users/signup", mh.regAndAuth.signUp)
	mh.router.HandleFunc("/v1/users/login", mh.regAndAuth.login)
}

// httpReqField is an helper method to build logger filed from an HTTPRequest
func httpReqField(statusCode int, r *http.Request, err error) []zap.Field {
	field := []zap.Field{
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.Int("status", statusCode),
		zap.Duration("duration", durationFromReqCtx(r)),
	}
	if err == nil {
		return field
	}
	field = append(field, zap.Error(err))
	return field
}

func durationFromReqCtx(r *http.Request) time.Duration {
	startTime, _ := r.Context().Value(contextKeyDuration).(time.Time)
	return time.Since(startTime)
}

func writeInternalServerError(w http.ResponseWriter, l *zap.Logger) {
	code := http.StatusInternalServerError
	w.WriteHeader(code)
	_, err := w.Write(someThingWentWrong)
	if err != nil {
		l.Error("writing to the response writer failed", zap.Error(err))
	}
}

func checkResponseWriteErr(err error, l *zap.Logger) {
	if err != nil {
		l.Error("response writer err", zap.Error(err))
	}
}
