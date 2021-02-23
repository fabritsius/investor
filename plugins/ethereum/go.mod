module github.com/fabritsius/investor/plugins/ethereum

go 1.16

require (
	github.com/fabritsius/envar v1.1.0
	github.com/fabritsius/investor/messages v0.0.0-20210224185547-f97a6a021c10
	google.golang.org/grpc v1.36.0
)

replace github.com/fabritsius/investor/messages => ../../messages/
