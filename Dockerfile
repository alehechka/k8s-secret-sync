# BUILD SERVER

FROM golang:1.19-alpine as go-builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

ARG RELEASE_VERSION=latest

RUN go build -ldflags "-X github.com/alehechka/kube-secret-sync.Version=${RELEASE_VERSION}" cmd/kube-secret-sync/main.go 

# SERVE

FROM scratch

COPY --from=go-builder /app/main kube-secret-sync

CMD [ "/kube-secret-sync", "start" ]
