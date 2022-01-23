FROM ubuntu:latest
COPY techtrainingcamp-AppUpgrade /root/server
COPY ./public/index.html /root/public/index.html
COPY ./redis.conf /root/redis.conf
EXPOSE 8080
EXPOSE 11451
ENV IS_DOCKER 1
RUN apt-get update
RUN apt-get install -y redis 
RUN redis-server --version 
RUN redis-server /root/redis.conf
CMD /root/server
