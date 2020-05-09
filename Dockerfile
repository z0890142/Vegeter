FROM golang:1.14-alpine as build
LABEL stage=builder
RUN apk add ca-certificates git gcc g++ libc-dev
WORKDIR /app

ENV GO111MODULE=on

COPY go.mod /
COPY go.sum .
RUN go mod download

FROM build_base as server_builder

COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o /vegeter .

FROM heroku/heroku:16
EXPOSE 80
COPY --from=server_builder /vegeter /vegeter
CMD ["/vegeter"]
