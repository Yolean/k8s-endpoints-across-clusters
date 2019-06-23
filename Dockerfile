FROM golang:1.12.6-alpine3.10@sha256:c750d6718009f2e94cb20f56a87884f601f175d43c9418ae0fa21ea00ad6a2ff
RUN apk add --no-cache git

WORKDIR /src

COPY go.* ./
RUN go mod download

COPY . .

#RUN go test

#RUN sed -i 's/zap.NewDevelopment()/zap.NewProduction()/' main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags '-w -extldflags "-static"'

FROM scratch

COPY --from=0 /src/k8s-endpoints-across-clusters /usr/local/bin/k8s-endpoints-across-clusters

ENTRYPOINT ["k8s-endpoints-across-clusters"]
