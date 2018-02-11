# use same go as kubernetes
FROM ubuntu:yakkety

USER root

RUN apt-get update && apt-get install -y ca-certificates

ADD openstack-token-test /usr/local/bin

CMD ["openstack-token-test"]
