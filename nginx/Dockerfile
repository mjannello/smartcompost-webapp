FROM nginx:latest

RUN mkdir -p /var/log/nginx && \
    chown -R nginx:nginx /var/log/nginx

COPY smartcompost.conf /etc/nginx/nginx.conf

