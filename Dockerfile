# Builder
FROM golang:1.13-alpine3.10 as builder

RUN apk update && apk upgrade && \
    apk --update add git

RUN mkdir -p /home/projects

WORKDIR /home/projects

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/news-service *.go


# Distribution
FROM alpine:3.10

RUN apk update && apk add --no-cache tzdata ca-certificates

RUN rm -rf /var/cache/apk/*

RUN update-ca-certificates 2>/dev/null || true

ENV TZ=Asia/Jakarta

WORKDIR /app

COPY --from=builder /home/projects/build /app

CMD ["/app/news-service"]
