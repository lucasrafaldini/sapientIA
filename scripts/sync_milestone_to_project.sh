#!/usr/bin/env zsh
set -euo pipefail

if ! command -v gh >/dev/null 2>&1; then
  echo "Erro: GitHub CLI (gh) não encontrado." >&2
  exit 1
fi

OWNER="lucasrafaldini"
MILESTONE="v0.1"
PROJECT_NUM=1

echo "Sincronizando issues do milestone '$MILESTONE' para o Project #$PROJECT_NUM..."

# Obter Project ID
PROJECT_ID=$(gh api graphql -f query="
{
  user(login: \"$OWNER\") {
    projectV2(number: $PROJECT_NUM) {
      id
    }
  }
}" -q '.data.user.projectV2.id')

if [[ -z "$PROJECT_ID" ]]; then
  echo "Erro: Projeto #$PROJECT_NUM não encontrado." >&2
  exit 1
fi

# Listar issues do milestone que ainda não estão no projeto
ISSUE_URLS=($(gh issue list --milestone "$MILESTONE" --state all --json url -q '.[].url'))

if [[ ${#ISSUE_URLS[@]} -eq 0 ]]; then
  echo "Nenhuma issue encontrada no milestone '$MILESTONE'."
  exit 0
fi

echo "Encontradas ${#ISSUE_URLS[@]} issues no milestone '$MILESTONE'."

# Listar items já no projeto
EXISTING_ITEMS=$(gh api graphql -f query="
{
  user(login: \"$OWNER\") {
    projectV2(number: $PROJECT_NUM) {
      items(first: 100) {
        nodes {
          content {
            ... on Issue {
              url
            }
          }
        }
      }
    }
  }
}" -q '.data.user.projectV2.items.nodes[].content.url')

# Adicionar issues que ainda não estão no projeto
for url in "${ISSUE_URLS[@]}"; do
  if echo "$EXISTING_ITEMS" | grep -q "$url"; then
    echo "  ✓ Já existe: $url"
  else
    echo "  + Adicionando: $url"
    gh project item-add "$PROJECT_NUM" --owner "$OWNER" --url "$url" >/dev/null || true
  fi
done

echo "Sincronização concluída!"
