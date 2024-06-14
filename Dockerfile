FROM golang:1.22-alpine as builder
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./bin/app ./cmd/app/main.go
RUN go build -o ./bin/worker ./cmd/worker/main.go

FROM golang:1.22-alpine as app
WORKDIR /app
COPY --from=builder /build/bin/app .
RUN chmod +x ./app
RUN mkdir migrations && mkdir api
COPY migrations ./migrations
COPY api ./api
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
CMD ["./app"]

FROM golang:1.22-alpine as worker
RUN apk --update --no-cache add tzdata dcron libcap
RUN cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime && echo "Europe/Moscow" > /etc/timezone
WORKDIR /worker
COPY --from=builder /build/bin/worker .
RUN chmod +x ./worker
COPY crontab /etc/crontabs/root
CMD ["/usr/sbin/crond", "-f", "-c", "/etc/crontabs", "-L", "/dev/stdout"]
