FROM alpine AS builder

RUN apk add --no-cache ca-certificates
ADD out/linux/devices /devices

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD out/linux/devices /devices

ENTRYPOINT ["/devices"]