# 🤖 SapientIA Agent Guide (VS Code)

## Princípios
1. **Contexto primeiro**: leia `CONTEXT_*` nesta pasta.
2. **Arquitetura estável**: respeite módulos e boundaries.
3. **Qualidade**: testes, logs claros e erros bem tratados.
4. **Docs vivas**: se impactar usuário, atualize README/Docs.
5. **Licenças**: Core sob Apache-2.0; Enterprise sob EULA.

## Tarefas prioritárias (MVP)
- CLI `sapientia run` conectando estágios do pipeline.
- Gerador de coocorrência e export GEXF.
- Árvore semântica (cooc → comunidades → hierarquia).
- Relatório HTML básico.
- Cliente do worker NLP (HTTP/gRPC).

## Estilo de commit
`type(scope): descrição` — ex: `feat(lexical): add chi-square`

## Revisão automática sugerida
- Lint sem erros.
- Testes passando.
- Sem dependências pesadas desnecessárias.
- Nomes e comentários humanos e claros.
