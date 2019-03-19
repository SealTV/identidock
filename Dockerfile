FROM debian:latest

COPY templates ./templates
ADD identidock ./

EXPOSE 5000:5000

ENTRYPOINT /identidock

