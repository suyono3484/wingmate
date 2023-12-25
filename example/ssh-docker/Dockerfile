FROM suyono/wingmate:alpine as source

FROM alpine:3.19

RUN apk update && apk add tzdata openssh-server && \
    ln -s /usr/share/zoneinfo/Australia/Sydney /etc/localtime && ssh-keygen -A

COPY --from=source /usr/local/bin/wingmate /usr/local/bin/
COPY --from=source /usr/local/bin/wmpidproxy /usr/local/bin/
ADD --chmod=755 example/ssh-docker/entry.sh /usr/local/bin/entry.sh
ADD --chmod=755 example/ssh-docker/etc /etc

ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]