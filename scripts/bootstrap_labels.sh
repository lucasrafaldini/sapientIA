#!/usr/bin/env zsh
set -euo pipefail

if ! command -v gh >/dev/null 2>&1; then
  echo "Erro: GitHub CLI (gh) não encontrado. Instale: https://cli.github.com/" >&2
  exit 1
fi

# Cria labels idempotentes
function ensure_label() {
  local name="$1"; local color="$2"; local desc="$3"
  gh label create "$name" --color "$color" --description "$desc" 2>/dev/null || gh label edit "$name" --color "$color" --description "$desc" >/dev/null
}

ensure_label "area:core"    "0366d6" "Core (pipeline, grafos, árvore)"
ensure_label "area:api"     "1d76db" "API/Serviços"
ensure_label "area:ui"      "a2eeef" "Interface web"
ensure_label "area:worker"  "0e8a16" "Worker Python/NLP"
ensure_label "area:report"  "fbca04" "Relatórios/exports"
ensure_label "area:infra"   "d4c5f9" "Infra, DX e CI"

ensure_label "kind:feature" "0e8a16" "Nova funcionalidade"
ensure_label "kind:bug"     "d73a4a" "Correção de bug"
ensure_label "kind:docs"    "0075ca" "Documentação"
ensure_label "kind:perf"    "c2e0c6" "Performance"
ensure_label "kind:refactor" "cfd3d7" "Refactor/Tech debt"

ensure_label "good first issue" "7057ff" "Bom para iniciantes"
ensure_label "help wanted"      "008672" "Ajuda bem-vinda"
ensure_label "blocked"          "b60205" "Bloqueado"

echo "Labels configuradas/atualizadas."
