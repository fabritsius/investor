module github.com/fabritsius/investor/aggregator

go 1.16

require (
	github.com/fabritsius/envar v1.1.0
	github.com/fabritsius/investor/messages v0.0.0-20210224185547-f97a6a021c10
	github.com/gocql/gocql v0.0.0-20210129204804-4364a4b9cfdd
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	google.golang.org/grpc v1.36.0
)

replace github.com/fabritsius/investor/messages => ../messages/
