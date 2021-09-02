FROM golang:1.16 as builder

WORKDIR /go/builds
COPY . .

RUN CGO_ENABLED=0 go build -o toolapp

#FROM gcr.io/distroless/static
FROM alpine:3.12

COPY --from=builder /go/builds/toolapp /

ENTRYPOINT ["/toolapp"]
