%define debug_package %{nil}

%global _dwz_low_mem_die_limit 0

%define _GOPATH %{_builddir}/go
%define googleapis_branch	master

%global provider        github
%global provider_tld	com
%global project         shatteredsilicon
%global repo            ssm-managed
%global provider_prefix	%{provider}.%{provider_tld}/%{project}/%{repo}

Name:		%{repo}
Version:	%{_version}
Release:	%{_release}
Summary:	Shattered Silicon Monitoring and Management management daemon

License:	AGPLv3
URL:		https://%{provider_prefix}
Source0:	%{name}-%{version}.tar.gz
Source1:	https://github.com/googleapis/googleapis/archive/%{googleapis_branch}/googleapis-%{googleapis_branch}.tar.gz

BuildRequires:	golang, protobuf, protobuf-devel, protobuf-compiler

%if 0%{?fedora} || 0%{?rhel} == 7
BuildRequires: systemd
Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd
%endif

%description
ssm-managed manages configuration of SSM server components (Prometheus,
Grafana, etc.) and exposes API for that.  Those APIs are used by ssm-admin tool.
See the SSM docs for more information.


%prep
%setup -q -n %{name}
%setup -q -T -D -a 1 -n %{name}

%build
export GOPATH=%{_GOPATH}
mkdir -p %{_GOPATH}/src
mkdir -p %{_GOPATH}/bin
export PATH=%{getenv:PATH}:%{_GOPATH}/bin
export GO111MODULE=off

mkdir -p vendor/github.com/cespare/xxhash/v2
find vendor/github.com/cespare/xxhash  ! -path '*/v2' -mindepth 1 -maxdepth 1 -exec mv {} vendor/github.com/cespare/xxhash/v2 \;

mkdir -p %{_GOPATH}/src/gopkg.in
cp -r vendor/gopkg.in/reform.v1 %{_GOPATH}/src/gopkg.in/
cp -r vendor %{_GOPATH}/src/gopkg.in/reform.v1/
pushd %{_GOPATH}/src/gopkg.in/reform.v1/
    rm -rf vendor/gopkg.in/reform.v1
    go install ./reform
popd

mkdir -p %{_GOPATH}/src/github.com/vektra
cp -r vendor/github.com/vektra/mockery %{_GOPATH}/src/github.com/vektra/
cp -r vendor %{_GOPATH}/src/github.com/vektra/mockery/v2/
pushd %{_GOPATH}/src/github.com/vektra/mockery/v2/
    rm -rf vendor/github.com/vektra/mockery
    go build -o mockery .
    mv mockery %{_GOPATH}/bin/
popd

mkdir -p %{_GOPATH}/src/github.com/golang
cp -r vendor/github.com/golang/protobuf %{_GOPATH}/src/github.com/golang/
cp -r vendor %{_GOPATH}/src/github.com/golang/protobuf/
pushd %{_GOPATH}/src/github.com/golang/protobuf/
    rm -rf vendor/github.com/golang/protobuf
    go install ./protoc-gen-go
popd

mkdir -p %{_GOPATH}/src/github.com/go-swagger
cp -r vendor/github.com/go-swagger/go-swagger %{_GOPATH}/src/github.com/go-swagger
cp -r vendor %{_GOPATH}/src/github.com/go-swagger/go-swagger/
pushd %{_GOPATH}/src/github.com/go-swagger/go-swagger/
    rm -rf vendor/github.com/go-swagger/go-swagger
    go install ./cmd/swagger
popd

mkdir -p %{_GOPATH}/src/github.com/grpc-ecosystem
cp -r vendor/github.com/grpc-ecosystem/grpc-gateway %{_GOPATH}/src/github.com/grpc-ecosystem
cp -r vendor %{_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/
pushd %{_GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/
    rm -rf vendor/github.com/grpc-ecosystem/grpc-gateway
    go install ./protoc-gen-grpc-gateway
    go install ./protoc-gen-swagger
popd

mkdir -p %{_GOPATH}/src/github.com/percona
cp -r vendor/github.com/percona/kardianos-service %{_GOPATH}/src/github.com/percona/

rm -f models/*_reform.go
go generate ./...
rm -fr api/*.pb.* api/swagger/*.json api/swagger/client api/swagger/models
protoc -Iapi -Igoogleapis-%{googleapis_branch} api/*.proto --go_out=plugins=grpc:api
protoc -Iapi -Igoogleapis-%{googleapis_branch} api/*.proto --grpc-gateway_out=logtostderr=true,request_context=true,allow_delete_body=true:api

mkdir -p %{_GOPATH}/src/%{provider}.%{provider_tld}/%{project}
cp -r $(pwd) %{_GOPATH}/src/%{provider_prefix}
go build -ldflags "${LDFLAGS:-} -s -w -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') -X 'github.com/shatteredsilicon/ssm-managed/utils.Version=%{version}-%{release}'" -a -v -x %{provider_prefix}/cmd/ssm-managed


%install
install -d -p %{buildroot}%{_bindir}
install -d -p %{buildroot}%{_sbindir}
install -p -m 0755 ssm-managed %{buildroot}%{_sbindir}/ssm-managed

install -d %{buildroot}/usr/lib/systemd/system
install -p -m 0644 %{name}.service %{buildroot}/usr/lib/systemd/system/%{name}.service


%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun %{name}.service


%files
%license LICENSE
%doc README.md
%{_sbindir}/ssm-managed
/usr/lib/systemd/system/%{name}.service


%changelog
* Thu Sep 21 2017 Mykola Marzhan <mykola.marzhan@percona.com> - 1.3.0-2
- add consul dependency for pmm-managed

* Tue Sep 12 2017 Mykola Marzhan <mykola.marzhan@percona.com> - 1.3.0-1
- init version