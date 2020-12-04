
Name: exim-amqp-pipe
Version: 0.0.1
Release: 2%{?dist}
Summary: Publish emails or attachments to RabbitMQ

License: MIT License
Source0: exim-amqp-pipe
Source1: exim-amqp-pipe.yaml
URL: https://github.com/boreevyuri/exim-amqp-pipe
BuildArch: x86_64

%description
Exim-amqp-pipe receives mails by pipe transport from exim (or any other pipe).
Parse attachments or embedded files and publishes them to RabbitMQ in base64

%build
echo "OK"

%install
mkdir -p %{buildroot}{%{_bindir},%{_sysconfdir}}
install -p -m 755 %{SOURCE0} %{buildroot}%{_bindir}
install -p -m 644 %{SOURCE1} %{buildroot}%{_sysconfdir}

%files
%defattr(0775,root,root,-)
%{_bindir}/exim-amqp-pipe
%attr(0644, root, root) %config(noreplace) %{_sysconfdir}/exim-amqp-pipe.yaml

%changelog
* Tue Nov 24 2020 Boreev Yuri <boreevyuri@gmail.com>
- first version