package common

import (
	"log"
)

// Err error logger
func Err(err interface{}, message ...string) {
	if err != nil {
		log.Fatal(message[0], err)
	}
}
