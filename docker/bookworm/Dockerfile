FROM golang:1.21-bookworm as builder

ADD . /root/wingmate
WORKDIR /root/wingmate/
RUN make all



FROM debian:bookworm

RUN ln -sf /usr/share/zoneinfo/Australia/Sydney /etc/localtime
COPY --from=builder /root/wingmate/cmd/wingmate/wingmate /usr/local/bin/wingmate
COPY --from=builder /root/wingmate/cmd/experiment/dummy/dummy /usr/local/bin/wmdummy
COPY --from=builder /root/wingmate/cmd/experiment/starter/starter /usr/local/bin/wmstarter
COPY --from=builder /root/wingmate/cmd/experiment/oneshot/oneshot /usr/local/bin/wmoneshot
COPY --from=builder /root/wingmate/cmd/experiment/spawner/spawner /usr/local/bin/wmspawner
COPY --from=builder /root/wingmate/cmd/pidproxy/pidproxy /usr/local/bin/wmpidproxy
ADD --chmod=755 docker/bookworm/entry.sh /usr/local/bin/entry.sh
ADD --chmod=755 docker/bookworm/etc /etc

ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]