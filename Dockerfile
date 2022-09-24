FROM golang:1.19-alpine AS development

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /go/bin/server

FROM golang:1.19-alpine as production

COPY --from=development /go/bin/server /go/bin/server

EXPOSE 8081

ENTRYPOINT [ "/go/bin/server" ]