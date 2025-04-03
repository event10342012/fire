FROM ubuntu:20.04

COPY fire /app/fire
WORKDIR /app

CMD ["/app/fire"]
