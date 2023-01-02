package unsecure_jwt_token

import (
	"fmt"
	"strings"

	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
)

const salt = "fawcoehfoiwefksdnfkbdksf"

var ErrWrongToken = errors.New("wrong token")

type jwtToken struct {
	header    map[string]string
	payload   map[string]interface{}
	signature string
}

func (m jwtToken) GetHeader() map[string]string {
	return m.header
}

func (m jwtToken) GetPayload() map[string]interface{} {
	return m.payload
}

func (m jwtToken) ToString() (string, error) {
	b, err := json.Marshal(m.header)
	if err != nil {
		return "", errors.Wrapf(err, "cannot marshal header: %#v", m.header)
	}
	encodeHeader := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(encodeHeader, b)

	b, err = json.Marshal(m.payload)
	if err != nil {
		return "", errors.Wrapf(err, "cannot marshal header: %#v", m.payload)
	}
	encodePayload := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.StdEncoding.Encode(encodePayload, b)

	signature := getSignature(string(encodeHeader), string(encodePayload))
	encodedSignature := make([]byte, base64.StdEncoding.EncodedLen(len(signature)))
	base64.StdEncoding.Encode(encodedSignature, []byte(signature))

	return fmt.Sprintf("%s.%s.%s", encodeHeader, encodePayload, encodedSignature), nil
}

func Create(header map[string]string, payload map[string]interface{}) (*jwtToken, error) {
	return &jwtToken{
		header:  header,
		payload: payload,
	}, nil
}

func Parse(token string) (*jwtToken, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.Wrapf(ErrWrongToken, "token parts must contains 3 part, given %d", len(parts))
	}

	decodedHeader, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.Wrapf(err, "cannot decode header")
	}

	header := make(map[string]string)
	err = json.Unmarshal(decodedHeader, &header)
	if err != nil {
		return nil, err
	}

	decodedPayload, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.Wrapf(err, "cannot decode payload")
	}

	payload := make(map[string]interface{})
	err = json.Unmarshal(decodedPayload, &payload)
	if err != nil {
		return nil, err
	}

	expectedSignature := getSignature(parts[0], parts[1])
	givenSignature, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, err
	}

	if expectedSignature != string(givenSignature) {
		return nil, errors.Wrapf(ErrWrongToken, "expected signature not equal given signature")
	}

	return &jwtToken{
		header:    header,
		payload:   payload,
		signature: expectedSignature,
	}, nil
}

func getSignature(encodedHeader string, encodedPayload string) (hash string) {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s.%s%s", encodedHeader, encodedPayload, salt)))
	return fmt.Sprintf("%x", sum)[:32]
}
