FROM golang:1.19-alpine as builder

WORKDIR /app

COPY ./* ./

RUN go build

FROM alpine:latest

COPY --from=builder /app/yals .

EXPOSE 8086

CMD [ "/yals" ]