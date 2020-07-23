# Build the Go API
FROM golang:latest
ADD . /app
WORKDIR /app/worker
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /worker .
CMD ["/worker"]