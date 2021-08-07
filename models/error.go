package models

// Error codes
const (
	InvalidConnectionDetailsError = 1000
	UnKnownAddressError           = 1001
	NotFoundError                 = 1002
	UsersNotFound                 = 1003
)

var (
	errorTypes = map[int]string{
		InvalidConnectionDetailsError: "InvalidConnectionDetails",
		UnKnownAddressError:           "UnKnownAddressError",
		NotFoundError:                 "NotFoundError",
		UsersNotFound:                 "UsersNotFound",
	}

	errorMessages = map[int]string{
		InvalidConnectionDetailsError: "failed to get peer from ctx",
		UnKnownAddressError:           "failed to get peer address",
		NotFoundError:                 "unknown identity",
		UsersNotFound:                 "no connected users",
	}
)

// Type ...
func Type(code int) string {
	if value, ok := errorTypes[code]; ok {
		return value
	}
	return "UnKnownError"
}

// Message ...
func Message(code int) string {
	if value, ok := errorMessages[code]; ok {
		return value
	}
	return "unknown"
}
