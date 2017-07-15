FROM golang

ADD server /go/src/github.com/nanashi74/VNLiner/server

WORKDIR /go/src/github.com/nanashi74/VNLiner/server
RUN go get ./...

RUN go install github.com/nanashi74/VNLiner/server

RUN apt-get update && apt-get install -y redis-server

RUN sed -i 's/bind 127.0.0.1/bind 0.0.0.0/' /etc/redis/redis.conf

RUN echo "service redis-server start" > /start.sh
RUN echo "sleep 10" >> /start.sh
RUN echo "server" >> /start.sh

CMD ["/bin/sh", "/start.sh"]

EXPOSE 8080
EXPOSE 6379