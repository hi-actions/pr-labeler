FROM golang as builder

WORKDIR /go/builds
COPY . .

RUN CGO_ENABLED=0 go build

FROM gcr.io/distroless/static

COPY --from=builder /go/builds/pr-labeler /

ENTRYPOINT ["/pr-labeler"]
