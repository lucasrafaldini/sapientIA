# Stopwords pt-BR

Este diretório contém listas de stopwords para português do Brasil com três perfis:

- pt-core.txt: curto e conservador (artigos, pronomes, preposições, conjunções, alguns advérbios funcionais)
- pt-extended.txt: ampliado (inclui pt-core + advérbios frequentes e algumas formas de auxiliares)
- pt-aggressive.txt: agressivo (inclui pt-extended + verbos de altíssima frequência e variações)

Guidas de uso:
- Comece com `pt-extended.txt` (bom equilíbrio para limpar “é”, “são”, “a”, “o”, etc.)
- Se precisar preservar mais verbos, use `pt-core.txt`.
- Se quiser focar em termos de conteúdo e reduzir ruído ao máximo, use `pt-aggressive.txt`.

Como configurar no pipeline (lexical/graph):

```yaml
params:
  stopwords_file: "data/stopwords/pt-extended.txt"
```

Cada arquivo é texto simples, uma palavra por linha. Linhas iniciadas com `#` são comentários e são ignoradas pelo parser atual (que trata toda linha como candidata; evite comentários se quiser 100% compatibilidade).
