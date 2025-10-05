FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o server .

EXPOSE 3000

CMD ["./server"]
