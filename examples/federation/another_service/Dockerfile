FROM node:10-alpine

WORKDIR /src

COPY *.json  /src/

RUN npm ci

COPY . /src/

CMD ["node", "/src/index.js"]