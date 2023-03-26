FROM alpine

RUN apk update && apk add --no-cache tzdata

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY ./bld/solaredge-exporter-$TARGETOS-$TARGETARCH /app/solaredge-exporter

EXPOSE 8080
ENTRYPOINT ["/app/solaredge-exporter"]
CMD []
