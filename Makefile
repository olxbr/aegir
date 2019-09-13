ENV ?= local
KERNEL := $$(uname | tr '[:upper:]' '[:lower:]')
GOPKGS := $$(go list ./... | grep -v /vendor/)
VERSION ?= latest
KUBECTL_VERSION = "v1.13.10"
SLACK_TOKEN ?= none

default: generate test build

build: generate
	CGO_ENABLED=0 go build -a -v -ldflags "-linkmode internal -extldflags -static"

build/dynamic: generate
	go build -v

clean:
	-rm -v aegir
	-docker rmi -f vivareal/aegir

deploy: kubectl
	sed -i "s/<IMAGE-TAG>/${VERSION}/g" deploy/deployment.yaml
	./kubectl apply --record -f deploy/
	test -n ${ENV} -a -d deploy/${ENV}/ && ./kubectl apply --record -f deploy/${ENV}/
	./kubectl rollout status deploy/aegir || ./kubectl rollout undo deploy/aegir

docker_image:
	docker build -t vivareal/aegir:latest .
  
docker_publish:
	docker tag vivareal/aegir:latest vivareal/aegir:${VERSION}
	docker push vivareal/aegir:${VERSION}

docker_test:
	docker run --rm --entrypoint /bin/sh vivareal/aegir:latest -c "make test"

generate:
	go generate ${GOPKGS}

install: generate
	CGO_ENABLED=0 go install -a -v -ldflags "-linkmode external -extldflags -static"

kubeconfig: kubectl
	@test ${ENV} || { echo ENV not set; exit 1; }
	@test ${CLUSTER_ENDPOINT} || { echo CLUSTER_ENDPOINT not set; exit 1; }
	@test ${CLUSTER_TOKEN} || { echo CLUSTER_TOKEN not set; exit 1; }
	@test ${NAMESPACE} || { echo NAMESPACE not set; exit 1; }
	@./kubectl config set-cluster aegir-${ENV} --server "${CLUSTER_ENDPOINT}" --insecure-skip-tls-verify
	@./kubectl config set-credentials aegir-${ENV} --token "${CLUSTER_TOKEN}"
	@./kubectl config set-context aegir-${ENV} --cluster aegir-${ENV} --user aegir-${ENV} --namespace ${NAMESPACE}
	@./kubectl config use-context aegir-${ENV}

kubectl:
	wget -O ./kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/${KERNEL}/amd64/kubectl
	chmod +x kubectl

notify-slack:
	docker run --net=host --rm byrnedo/alpine-curl -X POST -H 'Content-type:application/json' --data '{"text":"${text}","channel":"squad-platform","username":"circleci"}' https://hooks.slack.com/services/${SLACK_TOKEN}

tunnel:
	aws s3 cp s3://develop-br/scripts/circleci/circleci-v4-tunnel.sh ~/ && chmod +x ~/circleci-v4-tunnel.sh && ~/circleci-v4-tunnel.sh > /dev/null

test:
	go test -cover ${GOPKGS}

.PHONY: build build/dynamic clean deploy docker_image docker_publish docker_test generate install kubeconfig notify-slack tunnel test
