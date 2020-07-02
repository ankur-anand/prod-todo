package authstrategy

import (
	"crypto/rsa"
	"encoding/base64"
	"log"
	"os"
	"testing"

	"github.com/ankur-anand/prod-todo/pkg/authtoken"

	"github.com/dgrijalva/jwt-go/v4"
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

func TestJWT_Generate(t *testing.T) {
	t.Parallel()
	nJwt, err := authtoken.NewJWT(rsaPrK, rsaPuK, "test", "ankur", 5)
	if err != nil {
		t.Error(err)
	}
	_, err = nJwt.Generate("randomID")
	if err != nil {
		t.Error(err)
	}
}

func TestJWT_Validate(t *testing.T) {
	t.Parallel()
	userID := "randomUUUID"
	nJwt, err := authtoken.NewJWT(rsaPrK, rsaPuK, "test", "ankur", 5)
	if err != nil {
		t.Error(err)
	}
	token, err := nJwt.Generate(userID)
	if err != nil {
		t.Error(err)
	}
	id, err := nJwt.Validate(token)
	if err != nil {
		t.Error(err)
	}
	if id != userID {
		t.Errorf("expected jwt to return valid user id %s, got %s", userID, id)
	}
}

func BenchmarkJWT_Validate(b *testing.B) {
	userID := "randomUUUID"
	nJwt, err := authtoken.NewJWT(rsaPrK, rsaPuK, "test", "ankur", 5)
	if err != nil {
		b.Error(err)
	}
	token, err := nJwt.Generate(userID)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		_, err := nJwt.Validate(token)
		if err != nil {
			b.Error(err)
		}
	}
}
