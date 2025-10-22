#!/usr/bin/env zsh
set -euo pipefail

if ! command -v gh >/dev/null 2>&1; then
  echo "Erro: GitHub CLI (gh) não encontrado." >&2
  exit 1
fi

OWNER="lucasrafaldini"
PROJECT_NUM=1

echo "Configurando campos personalizados no Project #$PROJECT_NUM..."

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
  echo "Erro: Projeto #$PROJECT_NUM não encontrado para $OWNER." >&2
  exit 1
fi

echo "Project ID: $PROJECT_ID"

# Adicionar campo Status (single-select)
echo "Adicionando campo 'Status'..."
gh api graphql -f query="
mutation {
  createProjectV2Field(input: {
    projectId: \"$PROJECT_ID\"
    dataType: SINGLE_SELECT
    name: \"Status\"
    singleSelectOptions: [
      {name: \"📋 To do\", description: \"Not started\", color: GRAY}
      {name: \"🚧 In progress\", description: \"Work in progress\", color: YELLOW}
      {name: \"👀 In review\", description: \"Under review\", color: BLUE}
      {name: \"✅ Done\", description: \"Completed\", color: GREEN}
    ]
  }) {
    projectV2Field {
      ... on ProjectV2SingleSelectField {
        id
        name
      }
    }
  }
}" >/dev/null || echo "Campo 'Status' já existe ou erro ao criar."

# Adicionar campo Área (single-select)
echo "Adicionando campo 'Área'..."
gh api graphql -f query="
mutation {
  createProjectV2Field(input: {
    projectId: \"$PROJECT_ID\"
    dataType: SINGLE_SELECT
    name: \"Área\"
    singleSelectOptions: [
      {name: \"core\", description: \"Core pipeline\", color: BLUE}
      {name: \"api\", description: \"API/Services\", color: PURPLE}
      {name: \"ui\", description: \"Frontend\", color: PINK}
      {name: \"worker\", description: \"NLP Worker\", color: GREEN}
      {name: \"report\", description: \"Reports/Exports\", color: YELLOW}
      {name: \"infra\", description: \"Infra/DX/CI\", color: GRAY}
    ]
  }) {
    projectV2Field {
      ... on ProjectV2SingleSelectField {
        id
        name
      }
    }
  }
}" >/dev/null || echo "Campo 'Área' já existe ou erro ao criar."

# Adicionar campo Prioridade (single-select)
echo "Adicionando campo 'Prioridade'..."
gh api graphql -f query="
mutation {
  createProjectV2Field(input: {
    projectId: \"$PROJECT_ID\"
    dataType: SINGLE_SELECT
    name: \"Prioridade\"
    singleSelectOptions: [
      {name: \"🔴 Alta\", description: \"High priority\", color: RED}
      {name: \"🟡 Média\", description: \"Medium priority\", color: YELLOW}
      {name: \"🟢 Baixa\", description: \"Low priority\", color: GREEN}
    ]
  }) {
    projectV2Field {
      ... on ProjectV2SingleSelectField {
        id
        name
      }
    }
  }
}" >/dev/null || echo "Campo 'Prioridade' já existe ou erro ao criar."

echo "Campos personalizados configurados com sucesso!"
echo "Acesse: https://github.com/users/$OWNER/projects/$PROJECT_NUM"
