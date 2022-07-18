FROM node:16-alpine as frontend

COPY handler/restserver/static /src/

WORKDIR /src/

RUN yarn global add @vue/cli \
    && yarn install \
    && yarn build

FROM golang:1.18.4-buster as build
ENV TZ=Asia/Tokyo

WORKDIR /src/

COPY auth auth
COPY config config
COPY domain domain
COPY handler/grpcserver handler/grpcserver
COPY handler/restserver handler/restserver
COPY --from=frontend src/dist handler/restserver/dist
COPY handler handler
COPY infra infra
COPY proto proto
COPY protoconv protoconv
COPY service service
COPY utilsconv utilsconv
COPY *.go ./
COPY go.mod go.sum Makefile ./

RUN make bin/api-rback/embed/static

FROM alpine:3

RUN addgroup -g 70 -S api-rback \
    && adduser -u 70 -S -D -G api-rback -H -h /var/lib/api-rback -s /bin/sh api-rback \
    && mkdir -p /var/lib/api-rback \
    && chown -R api-rback:api-rback /var/lib/api-rback

COPY --from=build --chown=api-rback:api-rback /src/bin/api-rback .

RUN chmod +x api-rback \
    && mv api-rback /usr/local/bin/

STOPSIGNAL SIGINT

USER api-rback

CMD [ "api-rback" ]