FROM public.ecr.aws/docker/library/golang:1.19.2-alpine3.15 as builder

RUN apk add git
WORKDIR /src/todo_test
ADD . .

RUN apk update && apk add bash ca-certificates git gcc g++

# Unit tests
RUN GO111MODULE=on go test -tags musl $(go list ./...) -cover

# Build
RUN \
  VERSION=$(date '+%Y%m%d.%H%M%S') && \
  COMMIT=$(git rev-parse HEAD) && \
  BRANCH=$(git rev-parse --abbrev-ref HEAD) && \
  HOST=$(hostname) && \
  GO111MODULE=on \
  GOOS=linux \
  GOARCH=amd64 \
  go build \
    -tags musl \
    -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH} -X main.BUILDHOST=${HOST}" \
    -o /go/bin/. ./cmd/service

# ----  Now build final image  ----
FROM public.ecr.aws/docker/library/golang:1.19.2-alpine3.15

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ENV BIND_ADDRESS=0.0.0.0:${APP_PORT}

COPY --from=builder /go/bin/. /app/
WORKDIR /app

EXPOSE ${APP_PORT}
RUN date > BUILD_DATE
CMD ["./service"]
