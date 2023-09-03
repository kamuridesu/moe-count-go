FROM golang:1.21-alpine as build
ENV CGO_ENABLED=1

RUN apk add --no-cache \
    gcc \
    musl-dev

WORKDIR /workspace

COPY . /workspace/

RUN go mod tidy
RUN go build -ldflags='-s -w -extldflags "-static"' -o "moe-count"

FROM scratch
ENV GIN_MODE=release

WORKDIR /app/
COPY --from=build /workspace/moe-count /usr/local/bin/moe-count
COPY ./static /app/static
COPY ./template /app/template

ENTRYPOINT [ "moe-count", "-dbfile=/app/users.db" ]
