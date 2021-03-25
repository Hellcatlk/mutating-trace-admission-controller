IMAGE=trace-context-injector:v1
DELAY=1

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GO111MODULE=on go build -a -o bin/webhook ./cmd/webhook/main.go

docker: build
	docker rmi -f $(IMAGE)
	docker build -f build/Dockerfile -t $(IMAGE) .
	docker save -o bin/$(IMAGE).tar $(IMAGE)

install:
	hack/webhook-create-signed-cert.sh --service trace-context-injector-webhook-svc --secret trace-context-injector-webhook-certs --namespace default
	cat deploy/base/mutatingwebhook.yaml | hack/webhook-patch-ca-bundle.sh > deploy/base/mutatingwebhook-ca-bundle.yaml
	./tools/kustomize build deploy/base | kubectl apply -f -

remove:
	kubectl delete secret trace-context-injector-webhook-certs
	./tools/kustomize build deploy/base | kubectl delete -f -

.PHONY: test
test: gofmt golint govet gosec unit

gofmt:
	./hack/gofmt.sh

golint: bin/golangci-lint
	./bin/golangci-lint run ./... --timeout=10m

govet:
	go vet ./...

gosec: bin/gosec
	./bin/gosec -quiet ./...

unit:
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o coverage.html

deployment:
	kubectl apply -f test/yaml/deployment.yaml
	sleep $(DELAY)
	kubectl apply -f test/yaml/deployment_v2.yaml
	sleep $(DELAY)
	kubectl delete -f test/yaml/deployment_v2.yaml

deamonset:
	kubectl apply -f test/yaml/deamonset.yaml
	sleep $(DELAY)
	kubectl apply -f test/yaml/deamonset_v2.yaml
	sleep $(DELAY)
	kubectl delete -f test/yaml/deamonset_v2.yaml

statefulset:
	kubectl apply -f test/yaml/statefulset.yaml
	sleep $(DELAY)
	kubectl apply -f test/yaml/statefulset_v2.yaml
	sleep $(DELAY)
	kubectl delete -f test/yaml/statefulset_v2.yaml

replicaset:
	kubectl apply -f test/yaml/replicaset.yaml
	sleep $(DELAY)
	kubectl apply -f test/yaml/replicaset_v2.yaml
	sleep $(DELAY)
	kubectl delete -f test/yaml/replicaset_v2.yaml

pod:
	kubectl apply -f test/yaml/pod.yaml
	sleep $(DELAY)
	kubectl apply -f test/yaml/pod_v2.yaml
	sleep $(DELAY)
	kubectl delete -f test/yaml/pod_v2.yaml

clean:
	rm -rf bin/*
	rm -f deploy/base/mutatingwebhook-ca-bundle.yaml
	rm -f cover*
	docker rmi -f $(IMAGE)

# Install kustomize
bin/kustomize:
	./hack/install_kustomize.sh

# Install controller-gen
bin/controller-gen:
	./hack/install_controller-gen.sh

# Install golangci-lint
bin/golangci-lint:
	./hack/install_golangci-lint.sh

# Install gosec
bin/gosec:
	./hack/install_gosec.sh
