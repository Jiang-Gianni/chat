package config

import (
	"fmt"
)

func PrintListening(serviceName string, port string) {
	fmt.Printf("%s listening on port %s\n", serviceName, port)
	// time.Now().UTC()
	// time.RFC1123
}
