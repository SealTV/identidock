FROM golang

RUN groupadd -r gouser && useradd -r -g gouser gouser

WORKDIR /go/src/app
ADD ./app /go/src/app

RUN go get github.com/go-redis/redis
RUN go install app

EXPOSE 5000:5000
USER gouser

ENTRYPOINT /go/bin/app

