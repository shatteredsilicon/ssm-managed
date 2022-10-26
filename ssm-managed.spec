%define debug_package %{nil}

%global _dwz_low_mem_die_limit 0

%define _GOPATH %{_builddir}/go

%global provider        github
%global provider_tld	com
%global project         shatteredsilicon
%global repo            ssm-managed
%global provider_prefix	%{provider}.%{provider_tld}/%{project}/%{repo}

Name:		%{repo}
Version:	%{_version}
Release:	1%{?dist}
Summary:	Shattered Silicon Monitoring and Management management daemon

License:	AGPLv3
URL:		https://%{provider_prefix}
Source0:	%{name}-%{version}.tar.gz

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
%setup -q -n %{repo}

%build
export GOPATH=%{_GOPATH}
mkdir -p %{_GOPATH}/src
mkdir -p %{_GOPATH}/bin
export PATH=%{getenv:PATH}:%{_GOPATH}/bin
export GO111MODULE=off

mkdir -p vendor/github.com/cespare/xxhash/v2
find vendor/github.com/cespare/xxhash  ! -path '*/v2' -mindepth 1 -maxdepth 1 -exec mv {} vendor/github.com/cespare/xxhash/v2 \;

cp -r $(pwd)/vendor/* %{_GOPATH}/src/
go install gopkg.in/reform.v1/reform
go install github.com/golang/protobuf/protoc-gen-go
go install github.com/go-swagger/go-swagger/cmd/swagger
go install github.com/vektra/mockery/cmd/mockery
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
rm -f models/*_reform.go
go generate ./...
rm -fr api/*.pb.* api/swagger/*.json api/swagger/client api/swagger/models
protoc -Iapi -Ivendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis api/*.proto --go_out=plugins=grpc:api
protoc -Iapi -Ivendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis api/*.proto --grpc-gateway_out=logtostderr=true,request_context=true,allow_delete_body=true:api

mkdir -p %{_GOPATH}/src/%{provider}.%{provider_tld}/%{project}
cp -r $(pwd) %{_GOPATH}/src/%{provider_prefix}
go build -ldflags "${LDFLAGS:-} -s -w -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') -X 'github.com/shatteredsilicon/ssm-managed/utils.Version=%{version}'" -a -v -x %{provider_prefix}/cmd/ssm-managed


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