package aws

import (
	"errors"
	"log"

	"github.com/aws/smithy-go"
)

func LogAwsError(err error) {
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Printf("failed to call service: %s, operation %s, error %v", oe.Service(), oe.Operation(), oe.Error())
		}
	}
}

func ToAPIError[T smithy.APIError](err error) *T {
	var value T
	errors.As(err, &value)
	return &value
}
