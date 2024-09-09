FROM golang:alpine3.20 as builder

LABEL version="1.0.0"

ENV PATH_PROJECT=/tgapi
ENV GO111MODULE=on
ENV GOSUMDB=off
ENV TARGET=/app

WORKDIR ${PATH_PROJECT}
COPY . ${PATH_PROJECT}
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/$TARGET/

FROM alpine:3.20
COPY --from=builder /tgapi/$TARGET /bin
CMD ["/bin/app"]
EXPOSE  8000
