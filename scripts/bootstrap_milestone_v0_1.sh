#!/usr/bin/env zsh
set -euo pipefail

if ! command -v gh >/dev/null 2>&1; then
  echo "Erro: GitHub CLI (gh) não encontrado. Instale: https://cli.github.com/" >&2
  exit 1
fi

# Cria/atualiza milestone v0.1 (idempotente) via API
TITLE="v0.1"
DESC="Alpha: CLI + pipeline YAML + coocorrência + exports"
DUE_ISO="2025-11-30T00:00:00Z"

REPO="$(gh repo view --json nameWithOwner -q .nameWithOwner)"
EXISTING_NUM="$(gh api -X GET "/repos/$REPO/milestones?state=all" -q ".[] | select(.title==\"$TITLE\") | .number" || true)"

if [[ -n "${EXISTING_NUM:-}" ]]; then
  echo "Milestone '$TITLE' já existe (#$EXISTING_NUM). Atualizando..."
  gh api -X PATCH "/repos/$REPO/milestones/$EXISTING_NUM" -f title="$TITLE" -f description="$DESC" -f due_on="$DUE_ISO" >/dev/null
else
  echo "Criando milestone '$TITLE'..."
  gh api -X POST "/repos/$REPO/milestones" -f title="$TITLE" -f description="$DESC" -f due_on="$DUE_ISO" >/dev/null
fi

echo "Milestone '$TITLE' pronto."
