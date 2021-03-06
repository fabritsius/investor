module github.com/fabritsius/investor/plugins/tinkoff

go 1.16

require (
	github.com/TinkoffCreditSystems/invest-openapi-go-sdk v0.6.1
	github.com/fabritsius/envar v1.1.0
	github.com/fabritsius/investor/messages v0.0.0-20210224185547-f97a6a021c10
	google.golang.org/grpc v1.36.0
)

replace github.com/fabritsius/investor/messages => ../../messages/
