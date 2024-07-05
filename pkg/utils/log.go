package utils

import (
	"fmt"
	"log"
)

// Log records a message in the logs in the following format:
// [service   |method    ] some message.
func Log(service, method, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("[%10s | %10s] %s", service, method, message)
}
