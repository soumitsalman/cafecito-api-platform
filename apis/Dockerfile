# syntax=docker/dockerfile:1.7
# Build a selected Cafecito API and run it with a co-located llama.cpp embedder.
# Build from apis/: docker build --build-arg SERVICE=beans -t cafecito-beans .
ARG SERVICE
ARG MODEL_URL=https://huggingface.co/mradermacher/F2LLM-v2-80M-GGUF/resolve/main/F2LLM-v2-80M.Q8_0.gguf
ARG LLAMA_CPP_VERSION=b10025
ARG LLAMA_CPP_SHA256=6d8fe90c31b46f43da2728d6a5059d90000857cef5bedd2878fad2494890d0cf
# ARG MODEL_SHA256=FB2A92E51DBA7120D5369704502520814775B27C116D27ED039690106BF224EE

FROM ubuntu:24.04 AS runtime-base

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl wget libgomp1 && \
    rm -rf /var/lib/apt/lists/*

FROM runtime-base AS embedder-assets
ARG MODEL_URL
ARG LLAMA_CPP_VERSION
ARG LLAMA_CPP_SHA256
# ARG MODEL_SHA256

# Optional checksum validation:
# echo "${MODEL_SHA256}  /app/.models/F2LLM-v2-80M.Q8_0.gguf" | sha256sum --check --status
RUN mkdir -p /assets/.models /assets/llama-cpp && \
    curl --fail --silent --show-error --location \
      --output /tmp/llama-cpp.tar.gz \
      "https://github.com/ggml-org/llama.cpp/releases/download/${LLAMA_CPP_VERSION}/llama-${LLAMA_CPP_VERSION}-bin-ubuntu-x64.tar.gz" && \
    echo "${LLAMA_CPP_SHA256}  /tmp/llama-cpp.tar.gz" | sha256sum --check --status && \
    tar --extract --gzip --file /tmp/llama-cpp.tar.gz \
      --strip-components=1 --directory=/assets/llama-cpp && \
    rm /tmp/llama-cpp.tar.gz && \
    wget --no-verbose --output-document=/assets/.models/F2LLM-v2-80M.Q8_0.gguf \
      "${MODEL_URL}" && \
    chmod +x /assets/llama-cpp/llama-server

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

# Run llama-server and the selected API in one Linux container.
FROM runtime-base

COPY --from=embedder-assets /assets/.models/ /app/.models/
COPY --from=embedder-assets /assets/llama-cpp/ /app/.tools/llama-cpp/

COPY --from=api-builder /out/api /app/api
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/api /app/entrypoint.sh

ENV LLAMA_SERVER=/app/.tools/llama-cpp/llama-server
ENV LLAMA_MODEL=/app/.models/F2LLM-v2-80M.Q8_0.gguf
ENV LLAMA_HOST=127.0.0.1
ENV LLAMA_PORT=10000
ENV LLAMA_CTX_SIZE=16384
ENV LLAMA_PARALLEL=32
ENV LLAMA_BATCH_SIZE=2048
ENV LLAMA_UBATCH_SIZE=512
ENV LD_LIBRARY_PATH=/app/.tools/llama-cpp
ENV PORT=8080
ENV EMBEDDER_BASE_URL=http://127.0.0.1:10000
ENV EMBEDDER_MODEL=F2LLM-v2-80M
ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
