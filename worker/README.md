# Python NLP Worker

Worker Python responsável por processamento de linguagem natural (NLP).

## Requisitos

- Python 3.10+
- spaCy pt-BR
- KeyBERT
- Whisper (transcrição)

## Setup

```bash
cd worker
python -m venv venv
source venv/bin/activate  # ou venv\Scripts\activate no Windows
pip install -r requirements.txt
python -m spacy download pt_core_news_lg
```

## Executar

```bash
# Modo dev
python -m worker.api

# Com hot-reload
uvicorn worker.api.main:app --reload --port 8001
```

## Estrutura

- `nlp/` - Módulos de NLP (lematização, embeddings, keywords)
- `api/` - API HTTP/gRPC para comunicação com Go
- `models/` - Schemas e tipos
