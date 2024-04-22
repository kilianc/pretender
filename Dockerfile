FROM golang:1.22.2 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN make test
RUN make build

FROM scratch AS release-stage

WORKDIR /

COPY --from=build-stage /app/bin/pretender /pretender

EXPOSE 8080

ENTRYPOINT ["/pretender"]
