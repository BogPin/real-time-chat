FROM golang:1.19-alpine3.16 AS server_builder

RUN apk add gcc libc-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin/ main.go

FROM scratch

COPY --from=server_builder /app/bin/main .

EXPOSE 8080

CMD ["./main"]
