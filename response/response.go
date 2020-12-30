package response

const authErrorStatus = 403
const internalServerError = 503

// AuthErrorResponse send error response via kong
func AuthErrorResponse() (statusCode int, body string, headers map[string][]string) {
	responseHeaders := make(map[string][]string)
	responseHeaders["Content-Type"] = append(responseHeaders["Content-Type"], "text/plain")

	return authErrorStatus, "Access denied", responseHeaders
}
