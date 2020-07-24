# Build the Go API
FROM golang:1.13
ADD . /app
WORKDIR /app/db-controller
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /db-controller .
CMD ["/db-controller"]