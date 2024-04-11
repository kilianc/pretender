FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN make test
RUN make build
# RUN CGO_ENABLED=0 GOOS=linux go build -o /pretender

FROM scratch AS release-stage

WORKDIR /

COPY --from=build-stage /app/bin/pretender /pretender

EXPOSE 8080

ENTRYPOINT ["/pretender"]
