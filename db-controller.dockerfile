# Build the Go API
FROM golang:latest
ADD . /app
WORKDIR /app/db-controller
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /db-controller .
CMD ["/db-controller"]