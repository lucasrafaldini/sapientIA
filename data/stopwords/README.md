# Stopwords pt-BR

Este diretório contém listas de stopwords para português do Brasil com três perfis:

- **pt-core.txt**: curto e conservador (artigos, pronomes, preposições, conjunções, alguns advérbios funcionais)
- **pt-extended.txt**: ampliado (inclui pt-core + advérbios frequentes e algumas formas de auxiliares)
- **pt-aggressive.txt**: agressivo (inclui pt-extended + verbos de altíssima frequência e variações)

## Guia de uso

- Comece com `pt-extended.txt` (bom equilíbrio para limpar "é", "são", "a", "o", etc.)
- Se precisar preservar mais verbos, use `pt-core.txt`.
- Se quiser focar em termos de conteúdo e reduzir ruído ao máximo, use `pt-aggressive.txt`.

## Como configurar

### Opção 1: Via flag CLI (recomendado)

Sobrescreve o perfil de stopwords para todo o pipeline sem editar o YAML:

```bash
# Usa pt-core.txt (padrão)
sapientia run pipelines/exemplo.yaml

# Usa pt-extended.txt
sapientia run pipelines/exemplo.yaml --stopwords extended

# Usa pt-aggressive.txt
sapientia run pipelines/exemplo.yaml --stopwords aggressive
```

### Opção 2: No arquivo YAML do pipeline

Configura explicitamente no step `lexical` ou `graph`:

```yaml
params:
  stopwords_file: "data/stopwords/pt-extended.txt"
```

**Ordem de prioridade:**
1. Se o YAML especifica `stopwords_file` → usa o arquivo do YAML
2. Se não especifica → aplica o perfil da flag `--stopwords` (padrão: `core`)

## Formato dos arquivos

Cada arquivo é texto simples, uma palavra por linha. Linhas iniciadas com `#` são comentários e são ignoradas pelo parser atual (que trata toda linha como candidata; evite comentários se quiser 100% compatibilidade).
