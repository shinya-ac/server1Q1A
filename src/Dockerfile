FROM golang:1.21
ENV TZ=Asia/Tokyo

WORKDIR /app/src

COPY . /app/src

COPY cert.pem /etc/ssl/cert.pem

RUN go mod tidy

CMD ["go","run","cmd/main.go"]