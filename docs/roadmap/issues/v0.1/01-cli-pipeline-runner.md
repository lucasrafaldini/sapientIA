# v0.1 — CLI `sapientia run` e Parser de Pipeline YAML

Labels sugeridas: `area:core`, `kind:feature`

## Contexto

Precisamos de um executor de pipelines YAML com DAG de etapas para habilitar o fluxo fim-a-fim.

## Objetivo

Implementar `sapientia run` que lê `pipelines/<nome>.yaml`, valida esquema e executa as etapas.

## Escopo

- CLI `sapientia` com subcomando `run`
- Parser YAML com validações (schemas, erros amigáveis)
- Execução em DAG (ingest → lexical → graph → tree → report)
- Logs claros por etapa

## Critérios de aceite

- `sapientia run pipelines/exemplo.yaml` executa sem erros
- Artefatos versionados em `runs/<nome>/`

## Notas

- Ver `docs/roadmap/v0.1-execution.md` (Trilha 1)
