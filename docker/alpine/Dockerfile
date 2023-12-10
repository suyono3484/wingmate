FROM golang:1.21-alpine as builder

ADD . /root/wingmate
WORKDIR /root/wingmate/
RUN apk add make && make all



FROM alpine:3.18

RUN apk add tzdata && ln -s /usr/share/zoneinfo/Australia/Sydney /etc/localtime
COPY --from=builder /root/wingmate/cmd/wingmate/wingmate /usr/local/bin/wingmate
COPY --from=builder /root/wingmate/cmd/experiment/dummy/dummy /usr/local/bin/wmdummy
COPY --from=builder /root/wingmate/cmd/experiment/starter/starter /usr/local/bin/wmstarter
COPY --from=builder /root/wingmate/cmd/experiment/oneshot/oneshot /usr/local/bin/wmoneshot
COPY --from=builder /root/wingmate/cmd/experiment/spawner/spawner /usr/local/bin/wmspawner
COPY --from=builder /root/wingmate/cmd/pidproxy/pidproxy /usr/local/bin/wmpidproxy
ADD --chmod=755 docker/alpine/entry.sh /usr/local/bin/entry.sh
ADD --chmod=755 docker/alpine/etc /etc

ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]