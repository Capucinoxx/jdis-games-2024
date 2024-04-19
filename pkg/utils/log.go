package utils

import (
	"fmt"
	"log"
)

// Log enregistre un message dans les logs sous le format suivant:
// [service   |method    ] um message quelconque.
func Log(service, method, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("[%10s | %10s] %s", service, method, message)
}
