---
Title: Docker
Weight: 20
---

This assumes a `Dockerfile` in the same folder as your `main` package with modules and vendoring.

```dockerfile
FROM golang:alpine as builder

RUN mkdir /build

COPY . /build

WORKDIR /build

RUN apk add --update musl-dev gcc go git mercurial

RUN env GO111MODULE=on go build -mod=vendor -o builds/bot . 

FROM alpine

RUN apk --no-cache add ca-certificates

RUN adduser -S -D -H -h /app appuser

USER appuser

RUN id

COPY --from=builder /build/builds/bot /app/

WORKDIR /app 

CMD ["./bot"] 
```