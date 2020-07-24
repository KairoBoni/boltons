# Build the Go API
FROM golang:1.13
ADD . /app
WORKDIR /app/arquivei-api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /arquivei-api .
CMD ["/arquivei-api"]