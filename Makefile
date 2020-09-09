ENV ?= local
KERNEL := $$(uname | tr '[:upper:]' '[:lower:]')
GOPKGS := $$(go list ./...)
VERSION ?= latest
KUBECTL_VERSION = "v1.17.4"
SLACK_TOKEN ?= none
GIT_REVISION := $$(git rev-parse --verify HEAD)
CLUSTER_ENDPOINT :=
AEGIR_SVC_NAME = "aegir.default.svc"
OS_FAMILY := $$(uname | tr "[:upper:]" "[:lower:]")


default: generate test build

build: generate
	CGO_ENABLED=0 go build -a -v -ldflags "-linkmode internal -extldflags -static"

build/dynamic: generate
	go build -v

docker_clean:
	-rm -v aegir
	-docker rmi -f vivareal/aegir

deploy: kubectl
	sed -i "s/<IMAGE-TAG>/${GIT_REVISION}/g" deploy/deployment.yaml
	./kubectl apply --record -f deploy/
	test -n ${ENV} -a -d deploy/${ENV}/ && ./kubectl apply --record -f deploy/${ENV}/
	./kubectl rollout status deploy/aegir || ./kubectl rollout undo deploy/aegir

docker_build:
	docker build --target builder -t vivareal/aegir:${GIT_REVISION}-build .
	docker build --target dry-app -t vivareal/aegir:${GIT_REVISION} .

docker_test:
	docker run --rm --entrypoint /bin/sh vivareal/aegir:${GIT_REVISION}-build -c "make test"

docker_publish:
	docker build --target dry-app . vivareal/aegir:${GIT_REVISION}
	docker tag vivareal/aegir
	docker push vivareal/aegir:${GIT_REVISION}

deps:
	-mkdir tools
	curl -L https://github.com/instrumenta/kubeval/releases/latest/download/kubeval-${OS_FAMILY}-amd64.tar.gz -o tools/kubeval.tar.gz
	tar xf tools/kubeval.tar.gz -C tools
	chmod +x tools/kubeval
	curl -Lo tools/kind https://kind.sigs.k8s.io/dl/v0.8.1/kind-${OS_FAMILY}-amd64
	chmod +x tools/kind
	curl -o tools/kubectl -L https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/${OS_FAMILY}/amd64/kubectl
	chmod +x tools/kubectl

clean_deps:
	rm -rf tools

lint_kube_manifests:
	./tools/kubeval -d kube-manifests

.ONESHELL:
smoke_test: deps lint_kube_manifests
	@mkdir tmp-tls ; ./genkey.sh tmp-tls ${AEGIR_SVC_NAME}
	./tools/kind create cluster && ./tools/kind load docker-image vivareal/aegir:${GIT_REVISION}
	./tools/kubectl create configmap --from-file etc/rules.yaml aegir-rules
	./tools/kubectl create secret tls aegir-tls --cert tmp-tls/webhook-server-tls.crt --key tmp-tls/webhook-server-tls.key
	sed -e "s/__REPO_IMAGE_TAG__/vivareal\/aegir\:${GIT_REVISION}/" -e "s/imagePullPolicy\:\ Always/imagePullPolicy\:\ Never/" kube-manifests/aegir.yaml | ./tools/kubectl apply -f -
	sed "s/__BASE64_CABUNDLE__/$$(base64 -w0 < tmp-tls/ca.crt)/" kube-manifests/validationwebhook.yaml | ./tools/kubectl apply -f -
	./tools/kubectl rollout status deploy aegir --timeout 5m
	-./tools/kubectl apply -f kube-manifests/bad-deployment.yaml > output.log 2>&1
	grep -i "error" output.log || false

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
	./kubectl cluster-info

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
