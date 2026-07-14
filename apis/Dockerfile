# Build a selected Cafecito API and run it with a co-located llama.cpp embedder.
# Build from apis/: docker build --build-arg SERVICE=beans -t cafecito-beans .
ARG SERVICE

FROM golang:1.26-alpine AS api-builder
ARG SERVICE

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN case "${SERVICE}" in \
      beans|espresso) ;; \
      *) echo "SERVICE must be one of: beans, espresso" >&2; exit 1 ;; \
    esac && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/api "./${SERVICE}"

# Run llama-server and the selected API in one Linux container.
FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates curl libgomp1 && \
    rm -rf /var/lib/apt/lists/*

COPY --from=api-builder /out/api /app/api
COPY llama-cpp/ubuntu-x64/ /app/llama-cpp/
COPY models/F2LLM-v2-80M.Q8_0.gguf /app/models/F2LLM-v2-80M.Q8_0.gguf
COPY entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/api /app/entrypoint.sh /app/llama-cpp/llama-server && \
    cd /app/llama-cpp && \
    ln -sf libggml.so.0.16.0 libggml.so.0 && \
    ln -sf libggml-base.so.0.16.0 libggml-base.so.0 && \
    ln -sf libllama.so.0.0.9999 libllama.so.0 && \
    ln -sf libllama-common.so.0.0.9999 libllama-common.so.0 && \
    ln -sf libmtmd.so.0.0.9999 libmtmd.so.0

ENV LLAMA_SERVER=/app/llama-cpp/llama-server
ENV LLAMA_MODEL=/app/models/F2LLM-v2-80M.Q8_0.gguf
ENV LLAMA_HOST=127.0.0.1
ENV LLAMA_PORT=10000
ENV LLAMA_CTX_SIZE=16384
ENV LLAMA_PARALLEL=32
ENV LLAMA_BATCH_SIZE=2048
ENV LLAMA_UBATCH_SIZE=512
ENV LD_LIBRARY_PATH=/app/llama-cpp
ENV PORT=8080
ENV EMBEDDER_BASE_URL=http://127.0.0.1:10000
ENV EMBEDDER_MODEL=F2LLM-v2-80M
ENV GIN_MODE=release

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
