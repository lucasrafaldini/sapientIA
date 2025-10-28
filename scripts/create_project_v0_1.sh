#!/usr/bin/env zsh
set -euo pipefail

if ! command -v gh >/dev/null 2>&1; then
  echo "Erro: GitHub CLI (gh) não encontrado. Instale: https://cli.github.com/" >&2
  exit 1
fi

OWNER="$(gh repo view --json owner -q .owner.login)"
TITLE="SapientIA v0.1"

# Descobre se já existe um projeto com este título
EXISTING_NUM="$(gh project list --owner "$OWNER" --format json -L 200 -q ".[] | select(.title==\"$TITLE\") | .number" || true)"

if [[ -n "${EXISTING_NUM:-}" ]]; then
  echo "Projeto '$TITLE' já existe (#$EXISTING_NUM) para $OWNER."
  PROJECT_NUM="$EXISTING_NUM"
else
  echo "Criando projeto '$TITLE' em $OWNER..."
  PROJECT_NUM="$(gh project create --owner "$OWNER" --title "$TITLE" --format json -q .number)"
  echo "Projeto criado: #$PROJECT_NUM"
fi

# Adiciona todas as issues do milestone v0.1 ao projeto
ISSUE_URLS=($(gh issue list --milestone "v0.1" --state all --json url -q '.[].url'))

if [[ ${#ISSUE_URLS[@]} -eq 0 ]]; then
  echo "Nenhuma issue encontrada no milestone v0.1."
  exit 0
fi

echo "Adicionando ${#ISSUE_URLS[@]} issues ao projeto #$PROJECT_NUM..."
for u in "${ISSUE_URLS[@]}"; do
  echo " - $u"
  # Tenta adicionar; ignora erros (ex.: já adicionada)
  gh project item-add "$PROJECT_NUM" --owner "$OWNER" --url "$u" >/dev/null || true
done

echo "Projeto '$TITLE' pronto com issues do milestone v0.1."
