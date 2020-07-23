# Build the Go API
FROM golang:latest
ADD . /app
WORKDIR /app/arquivei-api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /arquivei-api .
CMD ["/arquivei-api"]