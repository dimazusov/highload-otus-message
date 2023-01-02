package unsecure_jwt_token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJwtToken_ToString(t *testing.T) {
	header := map[string]string{
		"alg": "HS256",
	}
	payload := map[string]interface{}{
		"userId": 2,
	}
	jwtToken, err := Create(header, payload)
	require.Nil(t, err)
	require.Equal(t, jwtToken.header["alg"], "HS256")
	require.Equal(t, jwtToken.payload["userId"].(int), 2)

	givenToken, err := jwtToken.ToString()
	require.Nil(t, err)

	expectedToken := "eyJhbGciOiJIUzI1NiJ9.eyJ1c2VySWQiOjJ9.MDhjZWViMTBiNmNiNjBhYjMxOTdkYTZhYjBjODA0OWY="
	require.Equal(t, expectedToken, givenToken)
}

func TestParser_Parse(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiJ9.eyJ1c2VySWQiOjJ9.MDhjZWViMTBiNmNiNjBhYjMxOTdkYTZhYjBjODA0OWY="

	jwtToken, err := Parse(token)
	require.Nil(t, err)
	require.Equal(t, jwtToken.header["alg"], "HS256")
	require.Equal(t, jwtToken.payload["userId"].(float64), float64(2))
}

func TestParser_getSignature(t *testing.T) {
	header := `{"typ": "JWT","alg": "HS256"}`
	payload := `{"id": "1337","username": "bizone","iat": 1594209600,"role":"user"}`

	givenSignature := getSignature(header, payload)
	expectedSignature := "c9a3593c5767444815bc460e117a90ae"

	require.Equal(t, expectedSignature, givenSignature)
}
