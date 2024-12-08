FROM golang:1.23.4
WORKDIR /app 
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "cmd/api/main.go"]

