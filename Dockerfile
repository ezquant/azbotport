FROM golang:1.19.3 as builder

WORKDIR /go/src/github.com/ezquant/azbot

COPY go.mod go.sum ./
COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /bin/azbot github.com/ezquant/azbot/cmd/azbot
RUN mkdir -p /usr/share/azbot

FROM scratch

WORKDIR /app

COPY --from=builder /bin/azbot /bin/azbot
COPY --from=builder /usr/share/azbot /usr/share/azbot
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["azbot"]
