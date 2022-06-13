# Percona Monitoring and Management (PMM) management daemon

[![Build Status](https://travis-ci.com/percona/pmm-managed.svg?branch=master)](https://travis-ci.com/percona/pmm-managed)
[![Go Report Card](https://goreportcard.com/badge/github.com/percona/pmm-managed)](https://goreportcard.com/report/github.com/percona/pmm-managed)
[![CLA assistant](https://cla-assistant.percona.com/readme/badge/percona/pmm-managed)](https://cla-assistant.percona.com/percona/pmm-managed)

pmm-managed manages configuration of [PMM](https://www.percona.com/doc/percona-monitoring-and-management/index.html)
server components (Prometheus, Grafana, etc.) and exposes API for that. Those APIs are used by
[ssm-admin tool](https://github.com/shatteredsilicon/ssm-client).
