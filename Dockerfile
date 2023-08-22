FROM golang:1.21-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/main/main.go


FROM ubuntu
LABEL authors="ageev"

COPY --from=build ./app/server ./server

EXPOSE 8080
EXPOSE 18080
CMD ./server