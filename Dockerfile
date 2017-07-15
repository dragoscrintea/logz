FROM scratch

# requires statically linked go binary:
COPY deamon /deamon

# Statically linked go binary requires CA certs for
# SSL HTTP connections, fix this file into place.
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/deamon"]