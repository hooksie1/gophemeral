FROM golang:alpine as builder
WORKDIR /app
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
ADD . /app/
RUN mkdir /database
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -installsuffix cgo -o gophemeral .


FROM scratch

COPY --from=builder /app/gophemeral .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /database .

CMD ["./gophemeral", "start", "--nats"]
