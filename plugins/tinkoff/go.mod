module github.com/fabritsius/investor/plugins/tinkoff

go 1.16

require (
	github.com/TinkoffCreditSystems/invest-openapi-go-sdk v0.6.1
	github.com/fabritsius/envar v1.1.0
	github.com/fabritsius/investor/messages v0.0.0-20210218183334-c9662315d9b0
	google.golang.org/grpc v1.35.0
)

replace github.com/fabritsius/investor/messages => ../../messages/
