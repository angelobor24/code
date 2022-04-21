FROM golang:1.16-alpine

WORKDIR /demo_pok

COPY ./go.mod .
COPY ./go.sum .


COPY ./cmd/*.go ./
COPY ./config/*.go ./
COPY ./db/*.go ./
COPY ./handlerMessage/*.go ./
COPY ./payment/*.go ./
COPY ./poke/*.go ./
COPY ./server/*.go ./
COPY . .
RUN apk add build-base
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o app ./cmd/main.go

EXPOSE 8080

ENTRYPOINT [ "./app" ]