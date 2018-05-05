FROM scratch

ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/cacert.pem
ADD b3bot /
ADD .b3bot.yaml /
CMD ["/b3bot"]
