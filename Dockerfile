FROM golang:1.21.6

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/githubnemo/CompileDaemon@latest

ENTRYPOINT CompileDaemon --build="go build -o /main main.go" --command=/main