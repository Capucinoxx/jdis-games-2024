FROM haproxy:2.7

USER root

COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg

COPY generate.sh /usr/local/bin/generate.sh
RUN chmod +x /usr/local/bin/generate.sh

RUN mkdir -p /etc/haproxy/ssl && chmod -R 755 /etc/haproxy/ssl

ENTRYPOINT ["/usr/local/bin/generate.sh"]
CMD ["haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg"]
