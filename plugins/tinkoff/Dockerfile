 FROM golang:1.16.0-alpine
 
 WORKDIR $GOPATH/src/github.com/fabritsius/investor/plugin/tinkoff
 COPY . .
 
 # current go.mod has a replace() statement which requires local messages module
 # the next couple lines fix the installation by re-replacing the module
 # to leave replace() as it is future versions of "go install" should have an option to skip replaces
 RUN go mod edit -replace github.com/fabritsius/investor/messages=github.com/fabritsius/investor/messages@v0.0.0-20210224185547-f97a6a021c10
 RUN go mod download github.com/fabritsius/investor/messages
 
 RUN go install
 
 # port for gRPC connection
 ENV TINKOFF_PORT 7702
 EXPOSE 7702
 
 CMD ["tinkoff"]
