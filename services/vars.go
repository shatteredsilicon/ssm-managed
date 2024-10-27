package services

import (
	"os"
	"time"
)

// sqlCheckTimeout maximum time for connecting to the database and running all queries
var sqlCheckTimeout = 5 * time.Second

func init() {
	timeout, _ := time.ParseDuration(os.Getenv("SSM_SQL_CHECK_TIMEOUT"))
	if timeout > 0 {
		sqlCheckTimeout = timeout
	}
}

func SQLCheckTimeout() time.Duration {
	return sqlCheckTimeout
}
