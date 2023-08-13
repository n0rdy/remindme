package utils

func IsPortValid(port int) bool {
	return port >= 1 && port <= 65535
}
