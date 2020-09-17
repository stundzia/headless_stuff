package errors

import (
	"fmt"
	"github.com/prometheus/common/log"
)

func HandleJSONError(err error) {
	if err != nil {
		fmt.Println("JSON error: ", err)
		log.Warnf("Got JSON Error: %s", err)
	}
}

func HandleGenericError(err error) {
	if err != nil {
		fmt.Println("Got error: ", err)
		log.Warn("Error: %s", err)
	}
}

func HandleResponseWriteError(err error) {
	if err != nil {
		fmt.Println("Response write error: ", err)
		log.Warn("Response write error: %s", err)
	}
}

func HandleCDPError(err error, msg string) {
	if err != nil {
		fmt.Println(fmt.Sprintf("CDP engine job execution error: %s (%s)", err, msg))
		log.Warnf("CDP engine job execution error: %s (%s)", err, msg)
	}
}