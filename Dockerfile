# 1. Build Image
FROM golang:1.18-buster AS builder

WORKDIR /root

COPY cmd ./cmd

RUN go build -tags=jsoniter -o ./cmd/whatisthissong

# 2. Production Image
FROM ubuntu:20.04

WORKDIR /

COPY --from=builder /root/cmd/whatisthissong/config ./config
COPY --from=builder /root/cmd/whatisthissong/whatisthissong ./whatisthissong


EXPOSE 8080
#./dbaas-api.init start
ENTRYPOINT ["./whatisthissong"]