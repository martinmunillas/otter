package server

import "fmt"

func PortString(port int64) string {
	return fmt.Sprintf(":%d", port)
}
