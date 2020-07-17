// +build unit_tests all_tests

package resthandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankur-anand/prod-todo/pkg"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"golang.org/x/crypto/bcrypt"
)

type testTokenizer struct {
}

func (t testTokenizer) Validate(token string) (string, error) {
	return "", nil
}

func (t testTokenizer) Generate(id string) (string, error) {
	return "token", nil
}

func TestSignUpHandler(t *testing.T) {
	t.Parallel()
	l := zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel))
	defer l.Sync()
	mockRep := &_mockUserRepoStorage{}
	a := auth{logger: l, svc: pkg.NewRegAndAuthService(mockRep)}

	user := signUpForm{
		EmailID:   "ankur@example.com",
		Password:  "ajsjjssj",
		FirstName: "Ankur",
		LastName:  "Anand",
		Username:  "ankur-anand",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	tc := []struct {
		name        string
		want        int
		body        []byte
		returnFunc  func() pkg.UserModel
		returnStore func(model pkg.UserModel) (uuid.UUID, error)
	}{
		{
			name: "empty body json unmarshalling err",
			want: 400,
			body: nil,
			returnFunc: func() pkg.UserModel {
				t.Fatal("this should not have been called")
				return pkg.UserModel{}
			},
			returnStore: func(model pkg.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "duplicate registration",
			want: 409,
			body: body,
			returnFunc: func() pkg.UserModel {
				return pkg.UserModel{
					Email: "ankur@example.com",
				}
			},
			returnStore: func(model pkg.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "success registration",
			want: 201,
			body: body,
			returnFunc: func() pkg.UserModel {
				return pkg.UserModel{}
			},
			returnStore: func(model pkg.UserModel) (uuid.UUID, error) {
				return model.ID, nil
			},
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			mockRep.returnFunc = c.returnFunc
			mockRep.returnStore = c.returnStore
			req := httptest.NewRequest(http.MethodPost, "/v1/users/signup", bytes.NewBuffer(c.body))
			rr := httptest.NewRecorder()
			a.signUp(rr, req)

			if rr.Code != c.want {
				t.Errorf("Expected Status Code %d Got %d", c.want, rr.Code)
			}
		})
	}

}

func TestLoginHandler(t *testing.T) {
	t.Parallel()
	l := zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel))
	defer l.Sync()
	password := "ankuranand"
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	wrongPassEncrypt, err := bcrypt.GenerateFromPassword([]byte("wrong password"),
		bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	mockRep := &_mockUserRepoStorage{}
	tokenizer := testTokenizer{}
	a := auth{logger: l, svc: pkg.NewRegAndAuthService(mockRep), tokenizer: tokenizer}

	user := loginForm{
		EmailID:  "ankur@example.com",
		Password: password,
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	tc := []struct {
		name        string
		want        int
		body        []byte
		returnFunc  func() pkg.UserModel
		returnStore func(model pkg.UserModel) (uuid.UUID, error)
	}{
		{
			name: "empty body json unmarshalling err",
			want: 400,
			body: nil,
			returnFunc: func() pkg.UserModel {
				t.Fatal("this should not have been called")
				return pkg.UserModel{}
			},
			returnStore: func(model pkg.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "wrong password",
			want: 422,
			body: body,
			returnFunc: func() pkg.UserModel {
				return pkg.UserModel{
					Email:    "ankur@example.com",
					Password: string(wrongPassEncrypt),
				}
			},
			returnStore: func(model pkg.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "success login",
			want: 201,
			body: body,
			returnFunc: func() pkg.UserModel {
				return pkg.UserModel{
					Email:    "ankur@example.com",
					Password: string(encryptedPass),
				}
			},
			returnStore: func(model pkg.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
	}

	for _, c := range tc {
		t.Run(c.name, func(t *testing.T) {
			mockRep.returnFunc = c.returnFunc
			mockRep.returnStore = c.returnStore
			req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(c.body))
			rr := httptest.NewRecorder()
			a.login(rr, req)

			if rr.Code != c.want {
				t.Errorf("Expected Status Code %d Got %d", c.want, rr.Code)
			}
		})
	}

}
