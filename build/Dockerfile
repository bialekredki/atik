FROM golang:1.22-alpine3.19 AS build-app

WORKDIR /build

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

COPY config ./config

COPY lib ./lib

COPY models ./models

COPY pkg ./pkg

COPY web ./web

RUN templ generate .

RUN go build -o atik


FROM node:22-alpine3.19 as build-tailwind

WORKDIR /build

COPY package.json .

COPY package-lock.json .

COPY tailwind.config.js .

COPY web/templates ./web/templates

RUN npm install

RUN npx tailwindcss -o assets/css/styles.css --minify


FROM alpine:3.19

WORKDIR /opt/atik

COPY --from=build-tailwind /build/assets/css/styles.css /opt/atik/assets/css/styles.css
COPY --from=build-app /build/atik /opt/atik/atik 

COPY assets/js /opt/atik/assets/js

ENTRYPOINT [ "/opt/atik/atik" ]