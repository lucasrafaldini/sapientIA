package pipeline

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestE2E_Briefing testa pipeline completo com briefing.yaml
func TestE2E_Briefing(t *testing.T) {
	runE2EPipeline(t, "pipelines/examples/briefing.yaml", "runs/exemplo")
}

// TestE2E_BriefingMinFreq1 testa pipeline com min_freq=1
func TestE2E_BriefingMinFreq1(t *testing.T) {
	runE2EPipeline(t, "pipelines/examples/briefing_minfreq1.yaml", "runs/exemplo_mf1")
}

// TestE2E_Lusiadas testa pipeline completo com Os Lusíadas
func TestE2E_Lusiadas(t *testing.T) {
	runE2EPipeline(t, "pipelines/examples/lusiadas.yaml", "runs/lusiadas")
}

// TestE2E_DivinaComedia testa pipeline completo com Divina Comédia
func TestE2E_DivinaComedia(t *testing.T) {
	runE2EPipeline(t, "pipelines/examples/divina-comedia.yaml", "runs/divina_comedia")
}

func runE2EPipeline(t *testing.T, pipelinePath, expectedOutputDir string) {
	// Ajustar path relativo ao root do projeto
	// Os testes rodam em internal/pipeline, então subir 2 níveis
	rootPath := filepath.Join("..", "..")
	absPath := filepath.Join(rootPath, pipelinePath)

	// Carregar pipeline
	p, err := LoadPipeline(absPath)
	if err != nil {
		t.Fatalf("Erro ao carregar pipeline %s: %v", absPath, err)
	}

	// Validar
	if err := p.Validate(); err != nil {
		t.Fatalf("Validação falhou: %v", err)
	}

	// Executar
	executor := NewExecutor(p)
	// Ajustar workDir do executor para root
	executor.workDir = rootPath
	if err := executor.Execute(); err != nil {
		t.Fatalf("Execução falhou: %v", err)
	}

	// Verificações básicas de outputs
	outDir := p.Config["output_dir"]
	if outDir == "" {
		outDir = expectedOutputDir
	}

	// Caminhos relativos ao root
	outDir = filepath.Join(rootPath, outDir)

	// Checar corpus.json
	corpusPath := filepath.Join(outDir, "corpus.json")
	if !fileExists(corpusPath) {
		t.Errorf("corpus.json não foi gerado: %s", corpusPath)
	}

	// Checar lexical.json
	lexicalPath := filepath.Join(outDir, "lexical.json")
	if !fileExists(lexicalPath) {
		t.Errorf("lexical.json não foi gerado: %s", lexicalPath)
	} else {
		// Validar que não há números puros na freq
		checkNoNumericTerms(t, lexicalPath)
	}

	// Checar graph.gexf
	graphPath := filepath.Join(outDir, "graph.gexf")
	if !fileExists(graphPath) {
		t.Errorf("graph.gexf não foi gerado: %s", graphPath)
	}

	// Checar tree.json (pode estar vazio mas deve existir)
	treePath := filepath.Join(outDir, "tree.json")
	if !fileExists(treePath) {
		t.Errorf("tree.json não foi gerado: %s", treePath)
	}

	// Checar report.html
	reportPath := filepath.Join(outDir, "report.html")
	if !fileExists(reportPath) {
		t.Errorf("report.html não foi gerado: %s", reportPath)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func checkNoNumericTerms(t *testing.T, lexicalPath string) {
	b, err := os.ReadFile(lexicalPath)
	if err != nil {
		t.Fatalf("Erro ao ler lexical.json: %v", err)
	}

	var data map[string]any
	if err := json.Unmarshal(b, &data); err != nil {
		t.Fatalf("Erro ao parse lexical.json: %v", err)
	}

	// Checar freq
	if freq, ok := data["freq"].(map[string]any); ok {
		for term := range freq {
			if isNumericOnly(term) {
				t.Errorf("Termo numérico encontrado em freq: %q", term)
			}
		}
	}

	// Checar ngrams
	for key := range data {
		if strings.HasPrefix(key, "ngrams_") {
			if ngrams, ok := data[key].(map[string]any); ok {
				for ngram := range ngrams {
					parts := strings.Fields(ngram)
					for _, p := range parts {
						if isNumericOnly(p) {
							t.Errorf("Termo numérico encontrado em %s: %q (ngram: %q)", key, p, ngram)
						}
					}
				}
			}
		}
	}
}
