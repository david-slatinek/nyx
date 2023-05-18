FROM golang:1.20-alpine3.17 AS build

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -v main.go


FROM alpine:3.17

WORKDIR /api

COPY --from=build /api/main ./
COPY --from=build /api/index.html ./
RUN mkdir env
COPY --from=build /api/env/.env ./env

CMD ["./main"]