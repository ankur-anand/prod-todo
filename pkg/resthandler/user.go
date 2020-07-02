package resthandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ankur-anand/prod-todo/pkg"
	"github.com/ankur-anand/prod-todo/pkg/logger"
)

var (
	invalidJSON = []byte(`{"status": "ERROR", 
	"error": {"message": "invalid request.", "kind": "APIException"}`)

	invalidEmailAddress = []byte(`{"status": "ERROR", 
	"error": {"message": "invalid email address.", 
    "kind": "APIException"}`)

	invalidPassword = []byte(`{"status": "ERROR", 
	"error": {"message": "invalid Password Length should be > 8 and < 254.", 
    "kind": "APIException"}`)

	userAlreadyRegistered = []byte(`{"status": "ERROR", 
	"error": {"message": "email id already registered.", 
    "kind": "APIException"}`)

	userCreated = []byte(`{"status": "SUCCESS", 
	"success": {"message": "email registered.", 
    "kind": "APISuccess"}`)

	invalidCredential = []byte(`{"status": "SUCCESS", 
	"success": {"message": "invalid credentials.", 
    "kind": "APIException"}`)

	tokenString = `{"status": "SUCCESS", 
	"success": {"message": "%s", 
    "kind": "APISuccess"}`
)

// Tokenizer provide an abstraction to work with
// Validation and Generation of an Auth Token
type Tokenizer interface {
	Validate(token string) (string, error)
	Generate(id string) (string, error)
}

// auth encapsulates various types of handlerFunc
// that responds to various model api request
type auth struct {
	svc       pkg.RegAndAuthService
	logger    *logger.Logger
	tokenizer Tokenizer
}

// signUpForm type Decode the submitted json body.
type signUpForm struct {
	EmailID   string `json:"email_id"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (ar auth) signUp(w http.ResponseWriter, r *http.Request) {
	var err error
	var code int
	var body []byte

	body, err = ioutil.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			ar.logger.Error("err closing underlying stream", logger.Error(err))
		}
	}()

	if err != nil {
		code = http.StatusInternalServerError
		writeInternalServerError(w, ar.logger)

		ar.logger.Error("err reading body", httpReqField(code, r, err)...)
		return
	}
	// decode the json body.
	var signForm signUpForm
	err = json.Unmarshal(body, &signForm)

	if err != nil {
		code = http.StatusBadRequest
		w.WriteHeader(code)
		_, err = w.Write(invalidJSON)
		checkResponseWriteErr(err, ar.logger)

		ar.logger.Error("err unmarshalling json", httpReqField(code, r, err)...)
		return
	}

	// precondition
	code, err = ar.precondition(w, signForm.EmailID, signForm.Password)
	if err != nil {
		ar.logger.Error("precondition check failed", httpReqField(code, r, err)...)
	}

	ok, err := ar.svc.IsDuplicateRegistration(r.Context(),
		signForm.EmailID)
	if err != nil {
		code = http.StatusInternalServerError
		writeInternalServerError(w, ar.logger)

		ar.logger.Error("err IsDuplicateRegistration", httpReqField(code, r, err)...)
		return
	}

	if ok {
		code = http.StatusConflict
		w.WriteHeader(code)
		_, err = w.Write(userAlreadyRegistered)
		checkResponseWriteErr(err, ar.logger)

		ar.logger.Error("email already registered", httpReqField(code, r, err)...)
		return
	}

	user := pkg.UserModel{
		Email:     signForm.EmailID,
		Password:  signForm.Password,
		FirstName: signForm.FirstName,
		LastName:  signForm.LastName,
		Username:  signForm.Username,
	}
	_, err = ar.svc.StoreUser(r.Context(), user)
	if err != nil {
		writeInternalServerError(w, ar.logger)
		ar.logger.Error("err StoreUser", httpReqField(code, r, err)...)
		return
	}

	code = http.StatusCreated
	w.WriteHeader(code)
	_, err = w.Write(userCreated)
	checkResponseWriteErr(err, ar.logger)

	ar.logger.Info("user created", httpReqField(code, r, err)...)
}

type loginForm struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

func (ar auth) login(w http.ResponseWriter, r *http.Request) {
	var err error
	var code int
	var body []byte

	body, err = ioutil.ReadAll(r.Body)

	defer func() {
		err := r.Body.Close()
		if err != nil {
			ar.logger.Error("err closing underlying stream", logger.Error(err))
		}
	}()

	if err != nil {
		code = http.StatusInternalServerError
		writeInternalServerError(w, ar.logger)

		ar.logger.Error("err reading body", httpReqField(code, r, err)...)
		return
	}

	// decode the json body.
	var logForm loginForm
	err = json.Unmarshal(body, &logForm)
	if err != nil {
		code = http.StatusBadRequest
		w.WriteHeader(code)
		_, err = w.Write(invalidJSON)
		checkResponseWriteErr(err, ar.logger)

		ar.logger.Error("err unmarshalling json", httpReqField(code, r, err)...)
		return
	}

	// precondition
	code, err = ar.precondition(w, logForm.EmailID, logForm.Password)
	if err != nil {
		ar.logger.Error("precondition check failed", httpReqField(code, r, err)...)
	}

	ok, user, err := ar.svc.IsCredentialValid(r.Context(), logForm.EmailID, logForm.Password)
	if err != nil {
		code = http.StatusInternalServerError
		writeInternalServerError(w, ar.logger)
		ar.logger.Error("err IsCredentialValid", httpReqField(code, r, err)...)
		return
	}

	if !ok {
		code = http.StatusUnprocessableEntity
		w.WriteHeader(code)
		_, err = w.Write(invalidCredential)
		checkResponseWriteErr(err, ar.logger)
		ar.logger.Error("invalid Credential", httpReqField(code, r, err)...)
		return
	}

	token, err := ar.tokenizer.Generate(user.ID.String())
	if err != nil {
		code = http.StatusInternalServerError
		writeInternalServerError(w, ar.logger)
		ar.logger.Error("err generating token", httpReqField(code, r, err)...)
		return
	}

	resJSON := fmt.Sprintf(tokenString, token)

	code = http.StatusCreated
	w.WriteHeader(code)
	_, err = w.Write([]byte(resJSON))
	checkResponseWriteErr(err, ar.logger)

	ar.logger.Info("user logged in", httpReqField(code, r, err)...)
}

func (ar auth) precondition(w http.ResponseWriter, email, password string) (code int, err error) {

	if !ar.svc.IsValidEmail(email) {
		code = http.StatusPreconditionFailed
		w.WriteHeader(code)
		_, err = w.Write(invalidEmailAddress)
		checkResponseWriteErr(err, ar.logger)
		return
	}

	if !ar.svc.IsValidPassword(password) {
		code = http.StatusPreconditionFailed
		w.WriteHeader(code)
		_, err = w.Write(invalidPassword)
		checkResponseWriteErr(err, ar.logger)

		return
	}
	return
}
