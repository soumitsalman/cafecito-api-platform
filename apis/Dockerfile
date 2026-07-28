# syntax=docker/dockerfile:1.7
# Build a selected Cafecito API from apis/:
# docker build --build-arg SERVICE=beans -t cafecito-beans .

ARG SERVICE

FROM golang:1.26-alpine AS api-builder
ARG SERVICE

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
    go mod download

COPY . ./
RUN --mount=type=cache,id=go-build,target=/root/.cache/go-build \
    --mount=type=cache,id=go-mod,target=/go/pkg/mod \
    case "${SERVICE}" in \
      beans|espresso) ;; \
      *) echo "SERVICE must be one of: beans, espresso" >&2; exit 1 ;; \
    esac && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api "./${SERVICE}"

FROM golang:1.26-alpine

RUN apk add --no-cache ca-certificates

COPY --from=api-builder /out/api /app/api

ENV PORT=8080
ENV EMBEDDER_MODEL=F2LLM-v2-80M
ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["/app/api"]
