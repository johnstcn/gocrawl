FROM golang:alpine as builder

RUN apk add --no-cache git
RUN go get github.com/golang/dep/cmd/dep

COPY Gopkg.lock Gopkg.toml /go/src/github.com/johnstcn/gocrawl/
WORKDIR /go/src/github.com/johnstcn/gocrawl
RUN dep ensure -vendor-only

COPY . /go/src/github.com/johnstcn/gocrawl/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gocrawld github.com/johnstcn/gocrawl/cmd/gocrawld

FROM alpine:latest
LABEL "Author"="Cian Johnston <public@cianjohnston.ie>"

ENV HOST 0.0.0.0
ENV PORT 12345
EXPOSE 12345

WORKDIR /root
COPY --from=builder /go/src/github.com/johnstcn/gocrawl/gocrawld .
# so we can establish HTTPS connections
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates
ENTRYPOINT ["./gocrawld"]
CMD ["-host", ${HOST}, "-port", ${PORT}]
