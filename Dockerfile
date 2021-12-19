#https://github.com/codefresh-contrib/golang-sample-app/blob/master/Dockerfile.multistage
# Build image
FROM golang:1.16-alpine as build

WORKDIR /tmp/app

COPY *.go .
COPY go.* ./

RUN go mod download
COPY . .

RUN go build -o out/app .

# App image

FROM alpine:3.14

RUN mkdir /app
COPY --from=build /tmp/app/out/app /app/app

CMD [ "/app/app" ]
