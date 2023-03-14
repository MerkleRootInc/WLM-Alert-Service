# syntax=docker/dockerfile:1.2
FROM golang:1.19-alpine AS builder

# install ssh client and git
RUN apk add --no-cache openssh-client git

# needed for unit tests
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

# download public key for github.com
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# download public and private modules via ssh context
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN --mount=type=ssh go mod download

COPY . ./

RUN go test -v ./...

RUN go build -o /alert-service

FROM alpine:latest
WORKDIR /
COPY --from=builder /alert-service /alert-service
COPY --from=builder /lib /lib
COPY --from=builder /usr/share /usr/share
COPY --from=builder /usr/include /usr/include
COPY --from=builder /usr/lib /usr/lib

EXPOSE 8080

CMD ["/alert-service"]