#!/usr/bin/env bash

REGISTRY=harbor.liasica.com/auroraride/edocseal:$1
PORT=26611
MAINTAIN=http://10.17.0.15:5000/maintain/stop/9geUbBHvX3caRWl1

if [ "$1" = "prod" ]; then
	PORT=26610
fi

GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -tags=jsoniter,poll_opt -gcflags "all=-N -l" -o build/release/edocseal cmd/edocseal/main.go
docker build -t "$REGISTRY" .
docker push "$REGISTRY"

ssh root@118.116.4.16 -p $PORT "
	cd /var/www
	docker pull ${REGISTRY}
	curl ${MAINTAIN}
	docker compose stop edocseal
	docker compose rm -f edocseal
	docker compose up edocseal -d
	docker image prune -f -a
	docker container prune -f
	docker volume prune -f -a
"
