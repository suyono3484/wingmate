FROM golang:1.21-bookworm as builder

ADD . /root/wingmate
WORKDIR /root/wingmate/
RUN make all && make DESTDIR=/usr/local/bin/wingmate install



FROM debian:bookworm

RUN ln -sf /usr/share/zoneinfo/Australia/Sydney /etc/localtime && \
    apt update && apt install -y procps && \
    useradd -m -s /bin/bash user1
COPY --from=builder /usr/local/bin/wingmate/ /usr/local/bin/
ADD --chmod=755 docker/bookworm/entry.sh /usr/local/bin/entry.sh
ADD --chmod=755 docker/bookworm/etc /etc

ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]