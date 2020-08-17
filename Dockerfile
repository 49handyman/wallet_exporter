FROM golang:1.13 as builder

ADD . /app/
WORKDIR /app/
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/zcashd_exporter /usr/local/bin/
ENTRYPOINT ["zcashd_exporter"]
CMD ["--help"]
