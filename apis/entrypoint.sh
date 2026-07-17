#!/bin/sh
set -eu

llama_server="${LLAMA_SERVER:-/app/.tools/llama-cpp/llama-server}"
llama_model="${LLAMA_MODEL:-/app/.models/F2LLM-v2-80M.Q8_0.gguf}"
api_pid=""
llama_pid=""

if [ ! -x "${llama_server}" ]; then
  echo "llama-server is not executable: ${llama_server}" >&2
  exit 1
fi

if [ ! -f "${llama_model}" ]; then
  echo "GGUF model was not found: ${llama_model}" >&2
  exit 1
fi

terminate() {
  if [ -n "${api_pid}" ]; then
    kill -TERM "${api_pid}" 2>/dev/null || true
  fi
  if [ -n "${llama_pid}" ]; then
    kill -TERM "${llama_pid}" 2>/dev/null || true
  fi
}

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

trap terminate INT TERM

/app/api &
api_pid=$!

api_status=0
wait "${api_pid}" || api_status=$?

# The API is the authoritative process. Its exit always shuts down llama-server.
kill -TERM "${llama_pid}" 2>/dev/null || true
wait "${api_pid}" 2>/dev/null || true
wait "${llama_pid}" 2>/dev/null || true

exit "${api_status}"
