FROM golang:1.19-alpine AS gobuild

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -o ./bin/app ./cmd/main.go

FROM scratch

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV APP_ROOT /opt/app
ENV PATH /opt/app

COPY --from=gobuild /build/bin $APP_ROOT

EXPOSE 8080

USER 1000

# executable binary setted @ helm chart or docker-compose
ENTRYPOINT ["/opt/app/app"]