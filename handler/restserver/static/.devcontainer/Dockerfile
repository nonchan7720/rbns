FROM node:14-alpine

RUN yarn global add @vue/cli \
    && mkdir /src \
    && apk --no-cache add tzdata gcc libc-dev git make bash \
    && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime \
    && echo "Asia/Tokyo" >  /etc/timezone \
    && apk del tzdata

WORKDIR /src/

CMD [ "sh" ]