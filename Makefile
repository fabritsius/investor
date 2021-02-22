export TINKOFF_PORT=7702
export ETHEREUM_PORT=7703

include .env
export $(shell sed 's/=.*//' .env)


build:
	docker build -t tinkoff-plugin ./plugins/tinkoff
	docker build -t ethereum-plugin ./plugins/ethereum
	docker build -t aggregator ./aggregator

start:
	docker start main-tinkoff-plugin || docker run -d -p 7702:7702 --name main-tinkoff-plugin tinkoff-plugin
	docker start main-ethereum-plugin || docker run -d -p 7703:7703 --env-file .env --name main-ethereum-plugin ethereum-plugin
	docker start main-aggregator || docker run -d --network="host" --name main-aggregator aggregator

stop:
	docker stop main-tinkoff-plugin
	docker stop main-aggregator

tinkoff:
	cd plugins/tinkoff && go run main.go

ethereum:
	cd plugins/ethereum && go run main.go

aggregate:
	cd aggregator && go run main.go

client:
	cd clients/core && go run main.go

# use 4.0-beta3 cause :latest doesn't support platform (linux/arm64/v8)
# version 4.0-beta4 doesn't work because of this: https://issues.apache.org/jira/browse/CASSANDRA-16424
db:
	docker run --name main-cassandra -p 7000:7000 -p 7001:7001 -p 7199:7199 -p 9042:9042 -p 9160:9160 -d cassandra:4.0-beta3

keyspace:
	docker run -it --network="host" --rm cassandra:4.0-beta3 cqlsh --execute=\
	"CREATE KEYSPACE IF NOT EXISTS investor WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1 };"

cqlsh:
	docker run --name cqlsh-cassandra -it --network="host" --rm cassandra:4.0-beta3 cqlsh
