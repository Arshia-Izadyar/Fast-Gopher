FROM golang:1.22.1-alpine3.18
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
WORKDIR /app/src/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /api
CMD ["/api"]