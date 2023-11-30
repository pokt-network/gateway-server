package common

func IsHttpOk(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}
