FROM golang:1.18-alpine
WORKDIR /app
ADD ./ ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o service .

FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai
RUN apk add --no-cache ca-certificates
ADD entrypoint.sh /app/
COPY --from=0 /app/service /app/gexservice/service
ADD conf /app/gexservice/conf
ADD www /app/gexservice/www
EXPOSE 3831
ENTRYPOINT ["/app/entrypoint.sh"]
