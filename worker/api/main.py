"""
SapientIA NLP Worker - API HTTP
Expõe endpoints para processamento de linguagem natural.
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(
    title="SapientIA NLP Worker",
    description="Worker Python para processamento de linguagem natural",
    version="0.1.0",
)

# CORS para comunicação com Go e UI
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Ajustar em produção
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/")
async def root():
    return {
        "service": "SapientIA NLP Worker",
        "version": "0.1.0",
        "status": "online",
    }


@app.get("/health")
async def health():
    return {"status": "healthy"}


# Endpoints futuros:
# POST /lemmatize - Lematização de texto
# POST /embeddings - Geração de embeddings
# POST /keywords - Extração de keywords
# POST /transcribe - Transcrição de áudio (Whisper)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)
