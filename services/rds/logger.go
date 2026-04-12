package rds

import (
	"fmt"

	"github.com/aws/smithy-go/logging"
)

type awsLogger struct {
	loggerFunc func(args ...interface{})
}

func (l awsLogger) Logf(classification logging.Classification, format string, v ...interface{}) {
	l.loggerFunc(fmt.Sprintf(format, v...))
}
