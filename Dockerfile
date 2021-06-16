FROM golang:1.16.5-alpine3.13

WORKDIR /app
COPY . .

RUN apk update && \
    apk add git && \
    go get github.com/cespare/reflex && \
    go get -u github.com/gin-gonic/gin

EXPOSE 8080
CMD ["reflex", "-c", "reflex.conf"]