module github.com/fabritsius/investor/aggregator

go 1.16

require (
	github.com/fabritsius/envar v1.0.1
	github.com/fabritsius/investor/messages v0.0.0-20210128174819-fd4734fa4c44
	github.com/gocql/gocql v0.0.0-20210129204804-4364a4b9cfdd
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	google.golang.org/grpc v1.35.0
)

replace github.com/fabritsius/investor/messages => ../messages/
