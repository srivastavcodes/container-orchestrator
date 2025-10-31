FROM golang:1.24.9-alpine

WORKDIR /app

RUN apk add --no-cache bash curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 4000
CMD ["/bin/bash"]
