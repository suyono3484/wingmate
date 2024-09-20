FROM golang:1.22-alpine3.20 AS builder

ADD . /root/wingmate
WORKDIR /root/wingmate/
ARG TEST_BUILD
RUN apk update && apk add git make build-base && \
    CGO_ENABLED=1 make all && \
    make DESTDIR=/usr/local/bin/wingmate install



FROM alpine:3.20

RUN apk add tzdata && ln -s /usr/share/zoneinfo/Australia/Sydney /etc/localtime && \
    adduser -h /home/user1 -D -s /bin/sh user1 && \
    adduser -h /home/user2 -D -s /bin/sh user2
COPY --from=builder /usr/local/bin/wingmate/ /usr/local/bin/
ADD --chmod=755 docker/alpine/entry.sh /usr/local/bin/entry.sh
ADD --chmod=755 docker/alpine/etc /etc

ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]