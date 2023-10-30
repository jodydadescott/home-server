FROM fedora:latest
WORKDIR /

RUN dnf -y update

RUN dnf install -y iproute iputils bind-utils file hostname procps net-tools dnf-plugins-core findutils

ADD build/linux-amd64/home-server /usr/sbin/home-server
RUN chmod +x /usr/sbin/home-server

CMD ["/usr/sbin/home-server", "run", "-c", "/etc/home-server-config.yaml"]