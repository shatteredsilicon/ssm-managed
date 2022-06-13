# Contributing notes

## Local setup

Run `make init` to install dependencies.

You should also have `mysqld_exporter`, `rds_exporter`, and `postgres_exporter` binaries somewhere in you `$PATH`.
One way to get them is to install them using `go get`:
```sh
go get -u github.com/shatteredsilicon/mysqld_exporter
go get -u github.com/shatteredsilicon/rds_exporter
go get -u github.com/shatteredsilicon/postgres_exporter/...
```

You have to use Docker Compose to run most of the tests.

```sh
make up
```

```sh
make
```

Start pmm-managed with

```sh
make run
```

Swagger UI will be available on http://127.0.0.1:7772/swagger/.

## Vendoring

We use [dep](https://github.com/golang/dep) to vendor dependencies.
