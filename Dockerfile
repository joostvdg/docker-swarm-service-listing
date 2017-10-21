# build stage
FROM golang:latest AS build-env
RUN go get -v github.com/docker/docker/client/...
RUN go get -v github.com/docker/docker/api/...
RUN go get -v github.com/gorilla/mux/...
ADD src/ $GOPATH/flow-proxy-service-lister
WORKDIR $GOPATH/flow-proxy-service-lister
RUN go build -o main -tags netgo main.go

# final stage
FROM alpine
ENTRYPOINT ["/app/main"]
COPY --from=build-env /go/flow-proxy-service-lister/main /app/
RUN chmod +x /app/main
