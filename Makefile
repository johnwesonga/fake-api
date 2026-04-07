PROJECT?=github.com/johnwesonga/fake-api
APP?=fake-api
PORT?=9000

RELEASE?=0.0.1
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
CONTAINER_IMAGE?=docker.io/johnwesonga/${APP}

GOOS?=linux
GOARCH?=arm64

.PHONY: clean build container run push minikube

clean:
	rm -f ${APP}

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP}

container: build
	docker build -t $(CONTAINER_IMAGE):$(RELEASE) .

run: container
	docker stop $(APP) || true && docker rm $(APP) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "PORT=${PORT}" \
		$(CONTAINER_IMAGE):$(RELEASE)

push: container
	docker push $(CONTAINER_IMAGE):$(RELEASE)

minikube: push
	for t in $(shell find ./kubernetes -type f -name "*.yaml"); do \
		cat $$t | \
			sed -E "s/\{\{\s*\.Release\s*\}\}/$(RELEASE)/g" | \
			sed -E "s/\{\{\s*\.ServiceName\s*\}\}/$(APP)/g"; \
		printf "\n---\n"; \
	done > tmp.yaml
	kubectl apply -f tmp.yaml --validate=false