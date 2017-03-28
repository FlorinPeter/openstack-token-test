.PHONY:	build push validate deploy build-vcasyslogp01 validate-vcasyslogp01 push-vcasyslogp01 deploy-vcasyslogp01

TAG := $(shell date +%Y-%m-%d-%H-%M)

build:
	docker run --rm -v "${PWD}":/go/src/github.com/xrl/openstack-token-test -w /go/src/github.com/xrl/openstack-token-test golang:1.7 go build -v
	docker build -t registry.usw1.viasat.cloud/openstack-token-test:${TAG} .

push:
	docker push registry.usw1.viasat.cloud/openstack-token-test:${TAG}

deploy: build push
	kubectl --namespace="default" set image deploy/openstack-token-test main=registry.usw1.viasat.cloud/openstack-token-test:${TAG}
