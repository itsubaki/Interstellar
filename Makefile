BUILD := ${PWD}/_build

package:
	set -x
	-rm -rf ${BUILD}

	mkdir -p ${BUILD}/controller/{bin,conf}
	mkdir -p ${BUILD}/broker/{cache,compute,database,project,network}/{bin,conf}

	cd controller/controller;        go build -o ${BUILD}/controller/bin/controller
	cd broker/provider/aws/project;  go build -o ${BUILD}/broker/project/bin/project
	cd broker/provider/aws/network;  go build -o ${BUILD}/broker/network/bin/network
	cd broker/provider/aws/database; go build -o ${BUILD}/broker/database/bin/database
	cd broker/provider/aws/cache;    go build -o ${BUILD}/broker/cache/bin/cache
	cd broker/provider/aws/compute;  go build -o ${BUILD}/broker/compute/bin/compute

	cp controller/controller/index.html          ${BUILD}/controller/conf
	cp broker/provider/aws/project/template.yml  ${BUILD}/broker/project/conf
	cp broker/provider/aws/network/template.yml  ${BUILD}/broker/network/conf
	cp broker/provider/aws/database/template.yml ${BUILD}/broker/database/conf
	cp broker/provider/aws/cache/template.yml    ${BUILD}/broker/cache/conf
	cp broker/provider/aws/compute/template.yml  ${BUILD}/broker/compute/conf

up-controller:
	set -x
	PORT=:9080 INDEX=./controller/controller/index.html            ${BUILD}/controller/bin/controller

up-project:
	set -x
	PORT=:9081 TEMPLATE=${BUILD}/broker/project/conf/template.yml  ${BUILD}/broker/project/bin/project

up-environ:
	set -x
	PORT=:9082 TEMPLATE=${BUILD}/broker/network/conf/template.yml  ${BUILD}/broker/network/bin/network

up-database:
	set -x
	PORT=:9083 TEMPLATE=${BUILD}/broker/database/conf/template.yml ${BUILD}/broker/database/bin/database

up-cache:
	set -x
	PORT=:9084 TEMPLATE=${BUILD}/broker/cache/conf/template.yml    ${BUILD}/broker/cache/bin/cache

up-compute:
	set -x
	PORT=:9085 TEMPLATE=${BUILD}/broker/compute/conf/template.yml  ${BUILD}/broker/compute/bin/compute

register-localhost:
	set -x

	curl -sX POST localhost:9080/v1/register -d '{"url": "http://localhost:9081"}'  | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://localhost:9082"}'  | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://localhost:9083"}'  | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://localhost:9084"}'  | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://localhost:9085"}'  | jq .

build:
	set -x

	cd controller/controller;        docker build -t controller          .
	cd broker/provider/aws/project;  docker build -t broker.aws.project  .
	cd broker/provider/aws/network;  docker build -t broker.aws.network  .
	cd broker/provider/aws/database; docker build -t broker.aws.database .
	cd broker/provider/aws/cache;    docker build -t broker.aws.cache    .
	cd broker/provider/aws/compute;  docker build -t broker.aws.compute  .

	docker images

up:
	set -x

	docker-compose up -d
	docker ps

down:
	set -x

	docker-compose down
	docker ps

register:
	set -x

	curl -sX POST localhost:9080/v1/register -d '{"url": "http://interstellar_broker.aws.project_1:8080"}'  | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://interstellar_broker.aws.network_1:8080"}'  | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://interstellar_broker.aws.database_1:8080"}' | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://interstellar_broker.aws.cache_1:8080"}'    | jq .
	curl -sX POST localhost:9080/v1/register -d '{"url": "http://interstellar_broker.aws.compute_1:8080"}'  | jq .

catalog:
	set -x

	curl -s localhost:9081/v1/catalog | jq .
	curl -s localhost:9082/v1/catalog | jq .
	curl -s localhost:9083/v1/catalog | jq .
	curl -s localhost:9084/v1/catalog | jq .
	curl -s localhost:9085/v1/catalog | jq .

service:
	set -x

	curl -s localhost:9080/v1/service | jq .
	curl -s localhost:9080/v1/service/$(shell curl -s localhost:9080/v1/service | jq -r '.service[0].service_id') | jq .

instance:
	set -x

	curl -s localhost:9080/v1/instance | jq .

create:
	set -x

	curl -s X POST  localhost:9080/v1/instance -d '{"service_id": "$(shell curl -s localhost:9080/v1/service | jq -r '.service[0].service_id')", "name": "develop01", "parameter": {"project_name": "myproject", "region": "us-east-1", "domain": "example.com"}}' | jq .

prune:
	set -x

	docker image prune --force
	docker images

.PHONY:
