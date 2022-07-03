FROM golang:alpine


WORKDIR /srv
COPY . /srv/
RUN go build -v .

FROM alpine:latest
EXPOSE 8080

ARG PROJECT=proxy
COPY --from=0 /srv/$PROJECT /srv/main

ENTRYPOINT ["/srv/main"]
