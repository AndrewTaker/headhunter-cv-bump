FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

FROM node:22-alpine

WORKDIR /app

COPY --from=builder /app/build /app/build
COPY --from=builder /app/package*.json /app/

ENV NODE_ENV production
CMD ["node", "build"]
