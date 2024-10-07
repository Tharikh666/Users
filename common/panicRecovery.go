package common

import (
	"log"
)

// PanicRecovery Function for panic error
func PanicRecovery(Module string, Method string) {

	if r := recover(); r != nil {
		log.Println(Module, Method, "PANIC-RECOVERY", r)
	}
}
