FROM golang:1.21-alpine as builder

ADD . /root/wingmate
WORKDIR /root/wingmate/cmd/wingmate
RUN go build -v

WORKDIR /root/wingmate/cmd/experiment/dummy
RUN go build -v

WORKDIR /root/wingmate/cmd/experiment/starter
RUN go build -v




FROM alpine:3.18

RUN apk add tzdata && ln -s /usr/share/zoneinfo/Australia/Sydney /etc/localtime
COPY --from=builder /root/wingmate/cmd/wingmate/wingmate /usr/local/bin/wingmate
COPY --from=builder /root/wingmate/cmd/experiment/dummy/dummy /usr/local/bin/wmdummy
COPY --from=builder /root/wingmate/cmd/experiment/starter/starter /usr/local/bin/wmstarter
ADD --chmod=755 docker/alpine/entry.sh /usr/local/bin/entry.sh
ADD --chmod=755 docker/alpine/etc /etc

ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]