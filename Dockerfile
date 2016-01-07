FROM scratch

# requires compiled statically linked go binary
COPY deamon /deamon

ENTRYPOINT ["/deamon"]