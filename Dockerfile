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

RUN go build cmd/kube-secret-sync/main.go

# SERVE

FROM scratch

COPY --from=go-builder /app/main kube-secret-sync

CMD [ "/kube-secret-sync", "start" ]
# Error invoking remote method 'docker-run-container': Error: (HTTP code 400) unexpected - failed to create shim task: OCI runtime create failed: runc create failed: unable to start container process: exec: "/bin/sh": stat /bin/sh: no such file or directory: unknown