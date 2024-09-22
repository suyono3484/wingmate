FROM suyono/wingmate:test AS source


FROM debian:bookworm
COPY --from=source /usr/local/bin/ /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/entry.sh" ]
CMD [ "/usr/local/bin/wingmate" ]
