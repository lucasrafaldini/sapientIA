#!/usr/bin/env zsh
set -euo pipefail

if ! command -v gh >/dev/null 2>&1; then
  echo "Erro: GitHub CLI (gh) não encontrado. Instale: https://cli.github.com/" >&2
  exit 1
fi

# Garante milestone e labels
"$(dirname $0)/bootstrap_labels.sh"
"$(dirname $0)/bootstrap_milestone_v0_1.sh"

MS="v0.1"

function mk_issue() {
  local title="$1"; shift
  local body="$1"; shift
  local labels="$1"; shift
  gh issue create --title "$title" --body-file "$body" --label $labels --milestone "$MS"
}

ROOT="docs/roadmap/issues/v0.1"

mk_issue "v0.1: CLI 'sapientia run' + parser de pipeline YAML" "$ROOT/01-cli-pipeline-runner.md" "area:core,kind:feature"
mk_issue "v0.1: Léxico básico (freq, n-grams, TF-IDF)" "$ROOT/02-lexico-basico.md" "area:core,kind:feature"
mk_issue "v0.1: Grafo de coocorrência + export GEXF/GraphML/CSV/JSON" "$ROOT/03-grafo-coocorrencia-exports.md" "area:core,kind:feature"
mk_issue "v0.1: Árvore semântica inicial" "$ROOT/04-arvore-semantica.md" "area:core,kind:feature"
mk_issue "v0.1: Relatório HTML mínimo" "$ROOT/05-relatorio-html.md" "area:report,kind:feature"
mk_issue "v0.1: Docker Compose (dev) + Make (DX)" "$ROOT/06-docker-compose-dx.md" "area:infra,kind:feature"
mk_issue "v0.1: Smoke test de pipeline" "$ROOT/07-smoke-test-pipeline.md" "area:infra,kind:feature"

echo "Issues de v0.1 criadas."
