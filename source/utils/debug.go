package utils

import (
	"fmt"
	"os"
)

func DebugPrint(args ...interface{}) {
	if os.Getenv("GIN_MODE") == "debug" {
		fmt.Println("[DEBUG] ::", args)
	}
}
