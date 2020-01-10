.PHONY: build
SVC := rimegate
WEB_ALPINE_VERSION := 3.11
WEB_SVC := web-rimegate-dogelink-com
COMMIT := $(shell git log -1 --pretty='%h')
REPOSITORY := 172.16.16.2:2443/go
PUBLIC_REPOSITIORY = icydoge/web

.PHONY: pull build push

all: pull build push clean

publish: pull build push-public

web: web-pull web-build web-push clean

build:
	docker build -t ${SVC} .

pull:
	docker pull golang:alpine

push:
	docker tag ${SVC}:latest ${REPOSITORY}:${SVC}-${COMMIT}
	docker push ${REPOSITORY}:${SVC}-${COMMIT}

push-public:
	docker tag ${SVC}:latest ${PUBLIC_REPOSITIORY}:${SVC}
	docker push ${PUBLIC_REPOSITIORY}:${SVC}

clean:
	docker image prune -f

web-build:
	docker build -t ${WEB_SVC} --build-arg ALPINE_VERSION=${WEB_ALPINE_VERSION} ./web

web-pull:
	docker pull alpine:${WEB_ALPINE_VERSION}

web-push:
	docker tag ${WEB_SVC}:latest icydoge/web:${WEB_SVC}-${COMMIT}
	docker push icydoge/web:${WEB_SVC}-${COMMIT}