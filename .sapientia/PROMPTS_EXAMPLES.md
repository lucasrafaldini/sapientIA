# 💬 Exemplos de Prompts (para agentes)

## Pipeline/CLI
> Implemente `sapientia run` que leia `pipelines/*.yaml`, construa a DAG e execute ingest, prep, lexical, graph e report.

## NLP Worker
> Crie um endpoint `/lemmatize` no worker que recebe texto PT-BR e retorna tokens com PoS e lemma (spaCy).

## Grafo
> Escreva função em Go que gera matriz de coocorrência por janela de 5 tokens e exporta GEXF.

## Relatório
> Crie template HTML com seções: briefing, perguntas, árvores por pergunta, insights e anexos.

## Infra
> Dockerfile multi-stage para o binário Go + imagem do worker Python; Compose com Redis e Postgres.
