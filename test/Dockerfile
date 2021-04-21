FROM ubuntu:latest
RUN mkdir -p /test
RUN apt-get update && apt-get install netcat -y
COPY . .
CMD nc server 12345 < test.txt 
