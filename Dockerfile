FROM golang:1.24 as builder

WORKDIR /build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o application main.go

FROM alpine:latest as runner

WORKDIR /car_estimator_api_gateway

COPY --from=builder /build/application .
COPY --from=builder /build/.env .

CMD [ "./application" ]
