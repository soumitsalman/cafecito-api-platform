#!/usr/bin/env bash
# Inspect a Zudoku build for Pagefind, published Markdown, and llms.txt.
set -euo pipefail

ROOT="${1:-.}"
SEARCH_ROOT="$ROOT"

if [[ -d "$ROOT/docs/dist" ]]; then
  SEARCH_ROOT="$ROOT/docs/dist"
elif [[ -d "$ROOT/dist" ]]; then
  SEARCH_ROOT="$ROOT/dist"
fi

echo "Inspecting portal output under $SEARCH_ROOT"

pagefind_hit="$(find "$SEARCH_ROOT" -type d -name 'pagefind' | head -n 1 || true)"
if [[ -z "$pagefind_hit" ]]; then
  pagefind_hit="$(find "$SEARCH_ROOT" -name 'pagefind.js' -o -name 'pagefind-entry.json' | head -n 1 || true)"
fi
if [[ -z "$pagefind_hit" ]]; then
  echo "Missing Pagefind index (search.type pagefind)." >&2
  exit 1
fi
echo "Pagefind: $pagefind_hit"

llms="$(find "$SEARCH_ROOT" -name 'llms.txt' | head -n 1 || true)"
if [[ -z "$llms" ]]; then
  echo "Missing llms.txt" >&2
  exit 1
fi
echo "llms.txt: $llms"

md_hit="$(find "$SEARCH_ROOT" \( -name '*.md' -o -path '*/markdown/*' -o -path '*/md/*' \) | head -n 1 || true)"
if [[ -z "$md_hit" ]]; then
  echo "Missing published Markdown output (docs.publishMarkdown)." >&2
  exit 1
fi
echo "Markdown: $md_hit"

# Slugs must appear in llms.txt itself, not merely somewhere in dist.
required_slugs=(
  "/products/beans"
  "/products/espresso"
  "/start/api-keys"
  "/guides/pricing-limits"
  "/guides/mcp-ai-agents"
  "/api/overview"
)
for slug in "${required_slugs[@]}"; do
  if ! grep -Fqi -- "$slug" "$llms"; then
    echo "Required slug not found in llms.txt: $slug" >&2
    exit 1
  fi
done

forbidden=(
  "cupboard"
  "beansack"
  "HNSW"
  "pgvector"
  "X-API-KEY"
  "BACKEND_API_KEY"
  "same_as"
  "derived_from"
)
for term in "${forbidden[@]}"; do
  hits="$(grep -RIl -i -- "$term" "$SEARCH_ROOT" --include='*.txt' --include='*.md' --include='*.html' 2>/dev/null || true)"
  if [[ -n "$hits" ]]; then
    echo "Forbidden internal term '$term' found in published output:" >&2
    echo "$hits" >&2
    exit 1
  fi
done

echo "Portal output inspection passed."
