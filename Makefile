.PHONY:	build push deploy create-deploy

TAG := $(shell date +%Y-%m-%d-%H-%M)

build:
	docker run --rm -v "${PWD}":/go/src/github.com/xrl/openstack-token-test -w /go/src/github.com/xrl/openstack-token-test golang:1.7 go build -v
	docker build -t registry.usw1.viasat.cloud/openstack-token-test:${TAG} .

push: build
	docker push registry.usw1.viasat.cloud/openstack-token-test:${TAG}

deploy: push
	kubectl --namespace="default" set image deploy/openstack-token-test main=registry.usw1.viasat.cloud/openstack-token-test:${TAG}
	kubectl --namespace="default" rollout status deploy/openstack-token-test

create-deploy: push
	docker tag registry.usw1.viasat.cloud/openstack-token-test:${TAG} registry.usw1.viasat.cloud/openstack-token-test:latest
	docker push registry.usw1.viasat.cloud/openstack-token-test:latest
	-kubectl --namespace="default" create -f auth-secrets.yml
	kubectl --namespace="default" create -f deployment.yml
	kubectl --namespace="default" rollout status deploy/openstack-token-test
	kubectl --namespace="default" logs --follow $(kubectl --namespace=default get pods | grep openstack-token-test | awk '{print $1}')