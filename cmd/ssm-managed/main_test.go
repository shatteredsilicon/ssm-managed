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

package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackages(t *testing.T) {
	cmd := exec.Command("ssm-managed", "-h")
	b, err := cmd.CombinedOutput()
	require.EqualError(t, err, "exit status 2")

	out := string(b)
	assert.False(t, strings.Contains(out, "-httptest.serve"), `ssm-managed should not import package "net/http/httptest"`)
	assert.False(t, strings.Contains(out, "-test.run"), `ssm-managed should not import package "testing"`)
}
