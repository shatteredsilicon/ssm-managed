// pmm-managed
// Copyright (C) 2017 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package tests

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/shatteredsilicon/ssm-managed/models"
)

func OpenTestDB(t testing.TB) *sql.DB {
	t.Helper()

	db, err := models.OpenDB("unix", "/var/lib/mysql/mysql.sock", "", "ssm-managed", "ssm-managed", t.Logf)
	require.NoError(t, err)

	const testDatabase = "ssm-managed-dev"
	_, err = db.Exec("DROP DATABASE `" + testDatabase + "`")
	require.NoError(t, err)
	_, err = db.Exec("CREATE DATABASE `" + testDatabase + "`")
	require.NoError(t, err)

	db.Close()

	db, err = models.OpenDB("unix", "/var/lib/mysql/mysql.sock", testDatabase, "ssm-managed", "ssm-managed", t.Logf)
	require.NoError(t, err)
	return db
}
