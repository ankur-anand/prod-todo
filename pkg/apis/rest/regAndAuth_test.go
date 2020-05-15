// +build unit_tests all_tests

package rest

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ankur-anand/prod-todo/pkg/authtoken"
	"github.com/dgrijalva/jwt-go/v4"

	"golang.org/x/crypto/bcrypt"

	"github.com/ankur-anand/prod-todo/pkg/domain"
	"github.com/ankur-anand/prod-todo/pkg/logger"
	"github.com/google/uuid"
)

const (
	rsaPrivateKey = `LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUVvd0lCQUFLQ0FRRUE5b3MvKzkwZG1zZ1U5UFR1NzFuTVVCd3NvcmN0dTE3NHpGQm9BTGxRQjZFYUJGbEQKL3g4alZ5SEZDY0RVbE5ISUpRdlcxNFpUOWNoc0pzLzhMZFF5RmE1cEViVGJOMVpVaFREN0kxRkptOUhhL0Z1WgpaUjdWVVlFSnBkUkNiZEhpcVpUeFRUMmQ0TVBTcXRIaWVWQUdyNHZ4Yyt3aTlqSEoyMVRLanBlR3pKZGRlY1I0Ckdra0t4aExzRC9vUFRiTjBmRUxubEZhdnZ3VDB5QWFqSnk1QzNveWlhQzNVNlVBc1JFWkl2VW0xb0NMUmx6K0EKVGpEMU4rSHZNQ1hDY3BqK2NMQ3NvRFJ6Qm90M3lseC9RcnkwN1FPRjJDc1RkbnlUREVBdFIreWpUMHF2NmNZZwozaGUrZzN0elZaT2ZIZ0VHZHlCVE91aHFKNjBsQmc4V21HMmxwUUlEQVFBQkFvSUJBRit0TUhwMGw5Mk9VaHV4CnhkdmJGRi91WHlBU1NFd1RraWZ2K0R4M3JlZ1lDL211RHFZK0ZqL2xHZ3NyNnhPSnljc2VxaFJmeTh0eEtROXkKM1dHSG5Kd3ZZQlVBQTZhWStSbnJKVHJTZStkZGJFZE00TjJPTnFoM2xCL25uSlB6eEt1YzRudmdNcG1jUlBBSApuWVVJbWwrYnhtci9NNTRwT2pYRTFRcTdJUlBaMHUwTUFXN2ZFb3A0dnBKaVpDUnB4cE9lcXl6dFd6R0JhZnhmCkNmaWg2SUUzV29reDBXMHJtVXkrUG5QYnB4dFYvbElwNHo5ZTI3S3pCZ09wQm95aVhUbFdTdjBHcC9PVnozZGcKenRtRGc3SFRJdmVsWjhRQ094N0ZsT0xSbDBPQnMwcUtEdVVTRTlTY3ZXODRGTlJ2dGROZEVidXhpWENXakJ6QgpUSlkxVnMwQ2dZRUErYVgzNFFIUWhUUWtFNjRqSmNtOWZnUG5zaEp2czFTV0FBSUl5NzVKSnBoaXQ4NENmdTBHClNpTUtzd0IyZnJWU0lNaERIMUtvYXZaU0RPQ1R5T2tCc002enVEdG5ET1IzTHdQZmtLZlcyZkluUi9Qb2ZDODAKc3hnZnEzaWNOR2RvTHV3anFSRGl4TEhVZml0K3BsY01VaTE1VDNsTEtxZUlSOEE2cUxzeUJQTUNnWUVBL05FUAozOVRRZXo5TkZGeTNhVkkzbExmd1NKWC9nVGFqdEpFVGZUNnVTQW5QdXdiQkM5b05qS1lWZnBYbWE0K1BFM2g0CktZYzVPMlZtQzBCRkh3NDh1VW04dTRKdDNiRE5lMTBNSk4ybVVaN1hlNWlweWo4dGpVT2c2dzFuaXd4a3IyUWMKRVI3VytUK1MwMnk0c3dsOEJqM1lzQkFOb1R5a3o5OVJJNmV2TVFjQ2dZRUF2VklqM1RzbkN4MHpqc2tzVm1mYgprRWtkMkdrcTFIQjlJSnhxVnppQytRWHZOenkvbjhuWVR6aXIwSHUxWVBuWXdvdWNlNUNQc0M0RW8wZGNTNnlJClg3RWhrY0Zhc09oQmlpSUIxUTJ6WGF6S0pVTFBOLzRFbFJ6aVI0TTcwbkhwREV4LzdxS2psazdWdXFqNWJ1UHMKc0JWVVBmVGFFQXJreXFUNDF5Uy9GZ2tDZ1lCTTJ2VTdjME9wby9XM3NmUGo1YVdWNVZEN2ovWHJmd1BIT2E4MgpETjhJY1VzZ0xRNTBudVl0a3JQSUZxUEVvUkM2dDQyMytpNncyc05wdWpFSkh0Zmc4QVNhOEN5Y0QwcDRMVElxCjV1TFB2eno4aXMxYStWZk1zUGx6VzFEVjJYK21QZ1cyUXF6UmVyMFUzdUZMTkIvcStkUXN1Y1NhOW9lWDFlaWgKc1RFMTh3S0JnRW4vd0ZTVDZwK01OOUVGcUNUT2VIQVJQeVA3V05JWkxoZUtJZC9QRDB2M0xQdFcvdEMwQlhCMwpHSm5UUGMxMklCamNLbUNPcTg3NWlKM2t0TWIzUllVdmJCRXJscy82OWRYOG1VanVEdXNmL0E4RUFtN2lIWHRGCk4xVHJVbHpybUtZRUkrbnIrejhUbk5HYnAwRmYxckJjUU8rQTM1REFNSHMycW5BVTBSM1MKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQo=`
	rsaPublicKey  = `LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUE5b3MvKzkwZG1zZ1U5UFR1NzFuTQpVQndzb3JjdHUxNzR6RkJvQUxsUUI2RWFCRmxEL3g4alZ5SEZDY0RVbE5ISUpRdlcxNFpUOWNoc0pzLzhMZFF5CkZhNXBFYlRiTjFaVWhURDdJMUZKbTlIYS9GdVpaUjdWVVlFSnBkUkNiZEhpcVpUeFRUMmQ0TVBTcXRIaWVWQUcKcjR2eGMrd2k5akhKMjFUS2pwZUd6SmRkZWNSNEdra0t4aExzRC9vUFRiTjBmRUxubEZhdnZ3VDB5QWFqSnk1Qwozb3lpYUMzVTZVQXNSRVpJdlVtMW9DTFJseitBVGpEMU4rSHZNQ1hDY3BqK2NMQ3NvRFJ6Qm90M3lseC9RcnkwCjdRT0YyQ3NUZG55VERFQXRSK3lqVDBxdjZjWWczaGUrZzN0elZaT2ZIZ0VHZHlCVE91aHFKNjBsQmc4V21HMmwKcFFJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==`
)

var (
	rsaPrK *rsa.PrivateKey
	rsaPuK *rsa.PublicKey
)

func TestMain(m *testing.M) {
	// decode the private certificate
	cert, err := base64.StdEncoding.DecodeString(rsaPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	rsaPrK, err = jwt.ParseRSAPrivateKeyFromPEM(cert)
	if err != nil {
		log.Fatal(err)
	}

	// decode the private certificate
	cert, err = base64.StdEncoding.DecodeString(rsaPublicKey)
	if err != nil {
		log.Fatal(err)
	}

	rsaPuK, err = jwt.ParseRSAPublicKeyFromPEM(cert)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestSignUpHandler(t *testing.T) {
	t.Parallel()
	l, _ := logger.NewTesting(nil)
	defer l.Sync()
	mockRep := &_mockUserRepoStorage{}
	a := auth{logger: l, svc: domain.NewRegAndAuthService(mockRep)}

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
		returnFunc  func() domain.UserModel
		returnStore func(model domain.UserModel) (uuid.UUID, error)
	}{
		{
			name: "empty body json unmarshalling err",
			want: 400,
			body: nil,
			returnFunc: func() domain.UserModel {
				t.Fatal("this should not have been called")
				return domain.UserModel{
					Email: "ankur@example.com",
				}
			},
			returnStore: func(model domain.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "duplicate registration",
			want: 409,
			body: body,
			returnFunc: func() domain.UserModel {
				return domain.UserModel{
					Email: "ankur@example.com",
				}
			},
			returnStore: func(model domain.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "success registration",
			want: 201,
			body: body,
			returnFunc: func() domain.UserModel {
				return domain.UserModel{}
			},
			returnStore: func(model domain.UserModel) (uuid.UUID, error) {
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
	l, _ := logger.NewTesting(nil)
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
	// we are not stubbing
	tokenizer, err := authtoken.NewJWT(rsaPrK, rsaPuK, "test", "ankur", 5)
	if err != nil {
		t.Fatal(err)
	}
	a := auth{logger: l, svc: domain.NewRegAndAuthService(mockRep), tokenizer: tokenizer}

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
		returnFunc  func() domain.UserModel
		returnStore func(model domain.UserModel) (uuid.UUID, error)
	}{
		{
			name: "empty body json unmarshalling err",
			want: 400,
			body: nil,
			returnFunc: func() domain.UserModel {
				t.Fatal("this should not have been called")
				return domain.UserModel{
					Email: "ankur@example.com",
				}
			},
			returnStore: func(model domain.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "wrong password",
			want: 422,
			body: body,
			returnFunc: func() domain.UserModel {
				return domain.UserModel{
					Email:    "ankur@example.com",
					Password: string(wrongPassEncrypt),
				}
			},
			returnStore: func(model domain.UserModel) (uuid.UUID, error) {
				t.Fatal("this should not have been called")
				return model.ID, nil
			},
		},
		{
			name: "success login",
			want: 201,
			body: body,
			returnFunc: func() domain.UserModel {
				return domain.UserModel{
					Email:    "ankur@example.com",
					Password: string(encryptedPass),
				}
			},
			returnStore: func(model domain.UserModel) (uuid.UUID, error) {
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
