#!/bin/sh
set -eu

llama_server="${LLAMA_SERVER:-/app/llama-cpp/llama-server}"
llama_model="${LLAMA_MODEL:-/app/models/F2LLM-v2-80M.Q8_0.gguf}"

if [ ! -x "${llama_server}" ]; then
  echo "llama-server is not executable: ${llama_server}" >&2
  exit 1
fi

if [ ! -f "${llama_model}" ]; then
  echo "GGUF model was not found: ${llama_model}" >&2
  exit 1
fi

"${llama_server}" \
  --model "${llama_model}" \
  --embedding \
  --pooling last \
  --embd-normalize 2 \
  --verbosity 1 \
  --ctx-size "${LLAMA_CTX_SIZE:-16384}" \
  --parallel "${LLAMA_PARALLEL:-32}" \
  --batch-size "${LLAMA_BATCH_SIZE:-2048}" \
  --ubatch-size "${LLAMA_UBATCH_SIZE:-512}" \
  --host "${LLAMA_HOST:-127.0.0.1}" \
  --port "${LLAMA_PORT:-10000}" &
llama_pid=$!

attempt=0
until curl --fail --silent --show-error \
  "http://${LLAMA_HOST:-127.0.0.1}:${LLAMA_PORT:-10000}/health" >/dev/null; do
  if ! kill -0 "${llama_pid}" 2>/dev/null; then
    echo "llama-server exited before becoming ready" >&2
    wait "${llama_pid}" 2>/dev/null || true
    exit 1
  fi

  attempt=$((attempt + 1))
  if [ "${attempt}" -ge "${LLAMA_STARTUP_TIMEOUT:-120}" ]; then
    echo "llama-server did not become ready within ${attempt} seconds" >&2
    kill -TERM "${llama_pid}" 2>/dev/null || true
    wait "${llama_pid}" 2>/dev/null || true
    exit 1
  fi
  sleep 1
done

/app/api &
api_pid=$!

terminate() {
  kill -TERM "${api_pid}" "${llama_pid}" 2>/dev/null || true
}

trap terminate INT TERM

# Exit when either process exits, forwarding termination to the other process.
while kill -0 "${api_pid}" 2>/dev/null && kill -0 "${llama_pid}" 2>/dev/null; do
  sleep 1
done

terminate
wait "${api_pid}" 2>/dev/null || true
wait "${llama_pid}" 2>/dev/null || true
