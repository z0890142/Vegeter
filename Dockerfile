FROM golang:1.14-alpine as build_base
LABEL stage=builder
RUN apk add ca-certificates git gcc g++ libc-dev
WORKDIR /app

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
COPY /config /config
RUN go mod download

FROM build_base as server_builder

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /vegeter .

FROM alpine

WORKDIR /vegeter
COPY --from=server_builder /vegeter ./
COPY --from=server_builder /config ./config
EXPOSE 8088

CMD ./vegeter
