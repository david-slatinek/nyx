FROM node:20-alpine3.16

WORKDIR /app

COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile

COPY . .

RUN yarn build

ENV NODE_ENV=production

CMD ["yarn", "start"]