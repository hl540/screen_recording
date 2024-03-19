FROM golang:1.18.10-alpine3.16 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn

RUN go build -o /src/bin/server /src/cmd/server/main.go

FROM alpine:3.16

WORKDIR /app/cmd/server

COPY --from=builder /src/bin /app/cmd/server
COPY --from=builder /src/cmd/server/conf.yaml /app/cmd/server/conf.yaml
COPY --from=builder /src/template /app/template

EXPOSE 9999

CMD ["./server"]