FROM alpine:3.14
RUN apk add --no-cache go
WORKDIR /app
COPY . .
RUN go build main.go
CMD ["go", "run", "btcclock"]
EXPOSE 3000