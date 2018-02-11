FROM registry.access.redhat.com/rhel7/rhel:7.3

RUN INSTALL_PKGS="git make" && \
    yum install --setopt=tsflags=nodocs -y $INSTALL_PKGS --enablerepo=rhel-7-server-rpms,rhel-7-server-optional-rpms && \
    rpm -V $INSTALL_PKGS && \
    mkdir -p  /tmp/go/src/openstack-token-test && \
    cd /tmp/go/ && \
    curl -LO https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz && \
    tar -zxf go1.8.3.linux-amd64.tar.gz && \
    mv go /usr/local/

USER root

RUN export GOROOT=/usr/local/go && \
    export GOPATH=/tmp/go && \
    export GOOS=linux && \
    export GOARCH=amd64 && \
    export CGO_ENABLED=1 && \
    /usr/local/go/bin/go version && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go && \
    go build -v \
    yum clean all && \
    rm -rf /tmp/go && \
    cp openstack-token-test /usr/local/bin/
    
CMD ["openstack-token-test"]
