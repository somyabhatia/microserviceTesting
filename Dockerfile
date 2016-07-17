FROM alpine
WORKDIR /home
ADD ./dummi /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/dummi"]
