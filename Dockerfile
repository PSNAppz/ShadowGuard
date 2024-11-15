FROM golang:1.19

WORKDIR /app

# Downloading dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying over source code
COPY . /app

# Compiling code
RUN CGO_ENABLED=0 GOOS=linux go build -o shadowguard cmd/main.go
CMD ["./shadowguard"]
