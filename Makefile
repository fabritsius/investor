export TINKOFF_PORT=3000

tinkoff:
	cd plugins/tinkoff && go run main.go

aggregate:
	cd aggregator && go run main.go

# use 4.0-beta3 cause :latest doesn't support platform (linux/arm64/v8)
# version 4.0-beta4 doesn't work because of this: https://issues.apache.org/jira/browse/CASSANDRA-16424
db:
	docker run --name main-cassandra -p 7000:7000 -p 7001:7001 -p 7199:7199 -p 9042:9042 -p 9160:9160 -d cassandra:4.0-beta3

cqlsh:
	docker run --name cqlsh-cassandra -it --network="host" --rm cassandra:4.0-beta3 cqlsh
