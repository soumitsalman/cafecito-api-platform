# Build a selected Cafecito API and run it with a co-located TEI embedder.
# Build with: docker build --build-arg SERVICE=beans -f Dockerfile.tei .
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


# Bake the embedding model into the image so the runtime does not need Hub access.
FROM python:3.11-slim AS model-downloader

ARG MODEL_ID=avsolatorio/GIST-small-Embedding-v0
ENV HF_HUB_ENABLE_HF_TRANSFER=1

RUN pip install --no-cache-dir "huggingface_hub[hf_transfer]"
COPY . /build-context/

# A local cache may contain the model directly or under .models/<model-name>.
RUN set -eu; \
    model_cache=/build-context/.models; \
    model_name="$(basename "$MODEL_ID")"; \
    if [ -f "${model_cache}/config.json" ]; then \
        mv "${model_cache}" /model; \
    elif [ -f "${model_cache}/${model_name}/config.json" ]; then \
        mv "${model_cache}/${model_name}" /model; \
    else \
        hf download "$MODEL_ID" --local-dir /model; \
    fi


# Run TEI and the selected API in one container. Only the API port is published.
FROM ghcr.io/huggingface/text-embeddings-inference:cpu-latest-grpc

COPY --from=model-downloader /model /data
COPY --from=api-builder /out/api /app/api
COPY entrypoint.sh /app/entrypoint.sh

ENV MODEL_ID=/data
ENV PORT=8080
ENV EMBEDDER_BASE_URL=localhost:10000
ENV GIN_MODE=release

RUN chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
