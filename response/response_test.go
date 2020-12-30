package response

import (
	"testing"
)

func Test_AuthError(t *testing.T) {
	statusCode, respBody, headers := AuthErrorResponse()

	if statusCode != 403 {
		t.Errorf("invalid status code")
	}

	if respBody != "Access denied" {
		t.Errorf("invalid response body")
	}

	if headers["Content-Type"][0] != "text/plain" {
		t.Errorf("invalid headers")
	}
}
