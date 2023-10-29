FROM fedora:latest
WORKDIR /

RUN dnf -y update

RUN dnf install -y iproute iputils bind-utils file hostname procps net-tools dnf-plugins-core findutils

ADD build/linux-amd64/home-dns-server /usr/sbin/home-dns-server
RUN chmod +x /usr/sbin/home-dns-server

CMD ["/usr/sbin/home-dns-server", "run", "-c", "/etc/dnsconfig.yaml"]