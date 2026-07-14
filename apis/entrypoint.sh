#!/bin/sh
set -eu

text-embeddings-router \
  --model-id "${MODEL_ID}" \
  --hostname 127.0.0.1 \
  --port 10000 &
tei_pid=$!

/app/api &
api_pid=$!

terminate() {
  kill -TERM "${api_pid}" "${tei_pid}" 2>/dev/null || true
}

trap terminate INT TERM

# Exit when either process exits, forwarding termination to the other process.
while kill -0 "${api_pid}" 2>/dev/null && kill -0 "${tei_pid}" 2>/dev/null; do
  sleep 1
done

terminate
wait "${api_pid}" 2>/dev/null || true
wait "${tei_pid}" 2>/dev/null || true
