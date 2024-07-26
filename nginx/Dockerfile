FROM nginx:alpine

RUN apk add --no-cache openssl

ARG DOMAIN
ENV CERT_DIR=/etc/nginx/ssl

COPY generate.sh /usr/local/bin/generate.sh
RUN chmod +x /usr/local/bin/generate.sh

COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80 443

CMD ["/bin/sh", "-c", "/usr/local/bin/generate.sh && nginx -g 'daemon off;'"]
