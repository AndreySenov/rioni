FROM alpine:3.23

ARG TARGETPLATFORM

ENV HOME=/home/rioni
ENV RIONI=/usr/local/bin/rioni
ENV RIONI_RELAY_UPSTREAM=https://dns.google/dns-query,https://cloudflare-dns.com/dns-query
ENV RIONI_RELAY_CLIENT_DNS=8.8.8.8,1.1.1.1
ENV RIONI_SERVER_DNS_ADDRESS=:53
ENV RIONI_SERVER_HTTP_ADDRESS=:443
ENV RIONI_SERVER_HTTP_TLS_CERT_FILE=/etc/rioni/tls/rioni.crt
ENV RIONI_SERVER_HTTP_TLS_KEY_FILE=/etc/rioni/tls/rioni.key

EXPOSE 53/udp
EXPOSE 53/tcp
EXPOSE 443/tcp

COPY --chown=root:root ${TARGETPLATFORM}/rioni $RIONI

RUN apk --no-cache add ca-certificates \
    && addgroup rioni \
    && adduser -S -G rioni -h $HOME rioni \
    && mkdir -p /etc/rioni $HOME \
    && chown -R rioni:rioni $HOME /etc/rioni \
    && chmod 755 $RIONI \
    && rioni --version

USER rioni
WORKDIR $HOME

ENTRYPOINT ["rioni"]
