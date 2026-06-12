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
Source0:	%{name}-%{version}-%{release}.tar.gz
Source1:	https://github.com/googleapis/googleapis/archive/%{googleapis_branch}/googleapis-%{googleapis_branch}.tar.gz

BuildRequires:	golang >= 1.24, protobuf, protobuf-devel, protobuf-compiler

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
export PATH=%{getenv:PATH}:%{_GOPATH}/bin

GOTOOLCHAIN=local go install ./vendor/gopkg.in/reform.v1/reform
GOTOOLCHAIN=local go install ./vendor/github.com/vektra/mockery/v2
GOTOOLCHAIN=local go install ./vendor/github.com/golang/protobuf/protoc-gen-go
GOTOOLCHAIN=local go install ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
GOTOOLCHAIN=local go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
GOTOOLCHAIN=local go install ./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

mkdir -p %{_GOPATH}/src/github.com/percona
cp -r vendor/github.com/percona/kardianos-service %{_GOPATH}/src/github.com/percona/

rm -f models/*_reform.go
go generate ./...
rm -fr api/*.pb.* api/swagger/*.json api/swagger/client api/swagger/models
protoc -Iapi -Igoogleapis-%{googleapis_branch} api/*.proto --go_out=plugins=grpc:api
protoc -Iapi -Igoogleapis-%{googleapis_branch} api/*.proto --grpc-gateway_out=logtostderr=true,allow_colon_final_segments=true,request_context=true,allow_delete_body=true:api

GOTOOLCHAIN=local go build -ldflags "${LDFLAGS:-} -s -w -B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') -X 'github.com/shatteredsilicon/ssm-managed/utils.Version=%{version}-%{release}'" -a -v -x ./cmd/ssm-managed


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