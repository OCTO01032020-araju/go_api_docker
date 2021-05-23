FROM golang:1.16.4-buster
WORKDIR /app
COPY main .
EXPOSE 4343
CMD ["./main"]

