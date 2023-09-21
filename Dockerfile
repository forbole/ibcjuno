FROM golang:1.20-alpine AS builder
RUN apk update && apk add --no-cache make git
WORKDIR /go/src/github.com/forbole/ibcjuno
COPY . ./
RUN go mod download
RUN make build

FROM alpine:latest
WORKDIR /ibcjuno
COPY --from=builder /go/src/github.com/forbole/ibcjuno/build/ibcjuno /usr/bin/ibcjuno
CMD [ "ibcjuno" ]