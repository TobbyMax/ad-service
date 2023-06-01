FROM golang:1.20-alpine
LABEL authors="ageev"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server cmd/main/main.go

EXPOSE 8080
EXPOSE 18080

CMD ["./server"]