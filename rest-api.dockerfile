# Build the Go API
FROM golang:latest
ADD . /app
WORKDIR /app/rest-api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /rest-api .
CMD ["/rest-api"]