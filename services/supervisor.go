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

package services

import (
	"context"
	"errors"

	servicelib "github.com/percona/kardianos-service"
)

var (
	// ErrNoSuchFileOrDir no such file or directory error
	ErrNoSuchFileOrDir = errors.New("no such file or directory")
)

//go:generate mockery -name=Supervisor -case=snake

// Supervisor is an interface for supervisor.Supervisor for mock generation.
type Supervisor interface {
	// Start installs and starts a service.
	Start(ctx context.Context, config *servicelib.Config) error

	// Start stops and uninstalls a service.
	Stop(ctx context.Context, name string) error

	// Status returns nil if service is installed and running.
	// It returns error otherwise or if service status can't be determined.
	Status(ctx context.Context, name string) error
}
