FROM golang:1.13 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN make



FROM alpine:3.8

RUN apk --no-cache update && apk --no-cache upgrade && apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/svc .
COPY .conf.yml ./.conf.yml

EXPOSE 80

CMD ["./svc"]
