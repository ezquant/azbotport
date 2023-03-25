FROM golang:1.19.3 as builder

WORKDIR /go/src/github.com/ezquant/ezbot

COPY go.mod go.sum ./
COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /bin/ezbot github.com/ezquant/ezbot/cmd/ezbot
RUN mkdir -p /usr/share/ezbot

FROM scratch

WORKDIR /app

COPY --from=builder /bin/ezbot /bin/ezbot
COPY --from=builder /usr/share/ezbot /usr/share/ezbot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["ezbot"]
