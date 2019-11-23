FROM golang:latest as builder

ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/github.com/shidax-tech/speed-wifi-home-exporter
COPY . .

RUN go get -d
RUN go build -o /app .


FROM scratch

COPY --from=builder /app /app

EXPOSE 9999
ENTRYPOINT ["/app"]
CMD ["-listen=:9999"]
