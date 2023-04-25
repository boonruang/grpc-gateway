FROM golang:1.20.3-alpine3.17 as builder

WORKDIR /builder

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-w -s" -o app cmd/main.go

FROM alpine:3.17
USER nobody
COPY --from=builder /builder/app /usr/local/bin/app
COPY .env .
EXPOSE 8080

CMD "/usr/local/bin/app"