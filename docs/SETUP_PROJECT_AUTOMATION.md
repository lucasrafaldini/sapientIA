# Configuração de Secrets para Automação do Project

Para que a GitHub Action `add-to-project.yml` funcione, você precisa criar um Personal Access Token (PAT) com permissões de Projects.

## Passos para configurar

### 1. Criar Personal Access Token (Classic)

1. Acesse: https://github.com/settings/tokens
2. Clique em **"Generate new token"** → **"Generate new token (classic)"**
3. Configure:
   - **Note**: `SapientIA Project Automation`
   - **Expiration**: 90 days (ou custom)
   - **Scopes** (marque):
     - ✅ `repo` (Full control of private repositories)
     - ✅ `project` (Full control of projects)
     - ✅ `read:org` (Read org and team membership, read org projects)
4. Clique em **"Generate token"** e copie o token (só aparece uma vez!)

### 2. Adicionar Secret no Repositório

1. Acesse: https://github.com/lucasrafaldini/sapientIA/settings/secrets/actions
2. Clique em **"New repository secret"**
3. Configure:
   - **Name**: `ADD_TO_PROJECT_PAT`
   - **Value**: cole o token copiado
4. Clique em **"Add secret"**

### 3. Testar a Automação

Depois de configurar o secret, a Action será acionada automaticamente quando:

- Uma nova issue for criada com milestone v0.1
- Uma issue existente for marcada com milestone v0.1
- Uma issue for rotulada com a label v0.1

Para testar manualmente:

```bash
# Criar uma issue de teste
gh issue create --title "Teste de automação" --body "Issue de teste" --milestone "v0.1"
```

Depois verifique se ela aparece no Project: https://github.com/users/lucasrafaldini/projects/1

## Alternativa: usar GITHUB_TOKEN

Se preferir não criar um PAT, você pode tentar usar o token padrão do GitHub Actions, mas ele tem limitações com Projects v2. Edite `.github/workflows/add-to-project.yml` e troque:

```yaml
github-token: ${{ secrets.ADD_TO_PROJECT_PAT }}
```

por:

```yaml
github-token: ${{ secrets.GITHUB_TOKEN }}
```

⚠️ **Nota**: O `GITHUB_TOKEN` pode não ter permissões suficientes para Projects v2 em user-level projects. Se não funcionar, use o PAT.
