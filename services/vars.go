package services

import (
	"os"
	"strconv"
	"time"
)

var (
	// sqlCheckTimeout maximum time for connecting to the database and running all queries
	sqlCheckTimeout = 5 * time.Second

	// disableTablestatsLimit the table amount limit for deciding whether to collect tablestats metrics
	disableTablestatsLimit uint16 = 1000
)

func init() {
	timeout, _ := time.ParseDuration(os.Getenv("SSM_SQL_CHECK_TIMEOUT"))
	if timeout > 0 {
		sqlCheckTimeout = timeout
	}

	tablestatsLimit, err := strconv.ParseUint(os.Getenv("SSM_DISABLE_TABLESTATS_LIMIT"), 10, 16)
	if err == nil {
		disableTablestatsLimit = uint16(tablestatsLimit)
	}
}

func SQLCheckTimeout() time.Duration {
	return sqlCheckTimeout
}

func DisableTablestatsLimit() uint16 {
	return disableTablestatsLimit
}
