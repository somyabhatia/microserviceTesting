FROM ubuntu
WORKDIR /home
ADD ./dummi /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/dummi"]
