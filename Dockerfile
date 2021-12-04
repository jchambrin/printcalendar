FROM golang:1.16.10-alpine3.15 as builder
WORKDIR /build
COPY . .
RUN go doc embed
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM surnet/alpine-node-wkhtmltopdf:14.17.5-0.12.6-small
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/app /app/
WORKDIR /app
CMD ["./app"]