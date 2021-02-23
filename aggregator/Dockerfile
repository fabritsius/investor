 FROM golang:1.16.0-alpine
 
 WORKDIR $GOPATH/src/github.com/fabritsius/investor/aggregator
 COPY . .
 
 # current go.mod has a replace() statement which requires local messages module
 # the next couple lines fix the installation by re-replacing the module
 # to leave replace() as it is future versions of "go install" should have an option to skip replaces
 RUN go mod edit -replace github.com/fabritsius/investor/messages=github.com/fabritsius/investor/messages@v0.0.0-20210218183334-c9662315d9b0
 RUN go mod download github.com/fabritsius/investor/messages
 
 RUN go install
 
 # ports for gRPC connection
 ENV TINKOFF_PORT 7702
 EXPOSE 7702
 ENV ETHEREUM_PORT 7703
 EXPOSE 7703
 
 # cassandra ports
 EXPOSE 7000
 EXPOSE 7001
 EXPOSE 7199
 EXPOSE 9042
 EXPOSE 9160
 
 CMD ["aggregator"]
