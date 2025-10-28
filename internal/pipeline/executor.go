// Package pipeline implementa o executor de pipelines
package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Executor executa um pipeline completo
type Executor struct {
	pipeline *Pipeline
	workDir  string
}

// NewExecutor cria um novo executor para o pipeline
func NewExecutor(p *Pipeline) *Executor {
	return &Executor{
		pipeline: p,
		workDir:  ".",
	}
}

// Execute executa o pipeline completo
func (e *Executor) Execute() error {
	fmt.Printf("📋 Pipeline: %s\n", e.pipeline.Name)
	if e.pipeline.Description != "" {
		fmt.Printf("   %s\n", e.pipeline.Description)
	}
	fmt.Println()

	// Validar antes de executar
	if err := e.pipeline.Validate(); err != nil {
		return fmt.Errorf("validação falhou: %w", err)
	}

	// Criar diretório de output se especificado
	if outputDir, ok := e.pipeline.Config["output_dir"]; ok {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("erro ao criar output_dir: %w", err)
		}
		fmt.Printf("📁 Output: %s\n\n", outputDir)
	}

	// Executar steps sequencialmente
	for i, step := range e.pipeline.Steps {
		start := time.Now()
		fmt.Printf("▶️  Step %d/%d: %s (%s)\n", i+1, len(e.pipeline.Steps), step.Name, step.Type)

		if err := e.executeStep(&step); err != nil {
			return fmt.Errorf("step %s falhou: %w", step.Name, err)
		}

		elapsed := time.Since(start)
		fmt.Printf("   ✅ Concluído em %.2fs\n\n", elapsed.Seconds())
	}

	fmt.Println("🎉 Pipeline executado com sucesso!")
	return nil
}

// executeStep executa um step individual
func (e *Executor) executeStep(step *Step) error {
	switch step.Type {
	case "ingest":
		return e.executeIngest(step)
	case "lexical":
		return e.executeLexical(step)
	case "graph":
		return e.executeGraph(step)
	case "tree":
		return e.executeTree(step)
	case "report":
		return e.executeReport(step)
	default:
		return fmt.Errorf("tipo de step desconhecido: %s", step.Type)
	}
}

// executeIngest executa step de ingestão
func (e *Executor) executeIngest(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)

	// Ler texto de entrada
	if step.Input == "" {
		return fmt.Errorf("ingest: input é obrigatório")
	}
	data, err := os.ReadFile(step.Input)
	if err != nil {
		return fmt.Errorf("ingest: erro ao ler input '%s': %w", step.Input, err)
	}

	// Criar corpus JSON mínimo
	corpus := fmt.Sprintf(`{"documents":[{"id":"doc1","text":%q}]}`, string(data))

	// Escrever output
	if step.Output == "" {
		return fmt.Errorf("ingest: output é obrigatório")
	}
	if err := e.ensureOutputDir(step.Output); err != nil {
		return err
	}
	if err := os.WriteFile(step.Output, []byte(corpus), 0644); err != nil {
		return fmt.Errorf("erro ao criar output: %w", err)
	}
	return nil
}

// executeLexical executa análise lexical
func (e *Executor) executeLexical(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)

	// Ler corpus
	if step.Input == "" {
		return fmt.Errorf("lexical: input é obrigatório")
	}
	b, err := os.ReadFile(step.Input)
	if err != nil {
		return fmt.Errorf("lexical: erro ao ler input '%s': %w", step.Input, err)
	}

	// Parse simples do corpus (sem dependências): procurar campo "text"
	// Suporta schema {"documents":[{"id":"doc1","text":"..."}]}.
	text := extractFirstTextField(string(b))

	// Ler parâmetros
	minFreq := 1
	tfidfOn := true
	exportCSV := false
	ngramsList := []int{2}
	var stopwordsFile string
	if step.Params != nil {
		if v, ok := step.Params["min_freq"].(int); ok {
			minFreq = v
		} else if fv, ok := step.Params["min_freq"].(float64); ok {
			minFreq = int(fv)
		}
		if v, ok := step.Params["tfidf"].(bool); ok {
			tfidfOn = v
		}
		if v, ok := step.Params["export_csv"].(bool); ok {
			exportCSV = v
		}
		if v, ok := step.Params["stopwords_file"].(string); ok {
			stopwordsFile = v
		}
		if arr, ok := step.Params["ngrams"].([]any); ok {
			ngramsList = make([]int, 0, len(arr))
			for _, x := range arr {
				switch t := x.(type) {
				case int:
					if t >= 1 { ngramsList = append(ngramsList, t) }
				case float64:
					iv := int(t)
					if iv >= 1 { ngramsList = append(ngramsList, iv) }
				}
			}
			if len(ngramsList) == 0 {
				ngramsList = []int{2}
			}
		}
	}

	// Tokenização com stopwords opcionais
	tokens := simpleTokenizePTWithExternal(text, stopwordsFile)
	// Freq
	freq := countFreq(tokens)
	if minFreq > 1 {
		for k, v := range freq {
			if v < minFreq { delete(freq, k) }
		}
	}
	// N-grams
	ngramsByN := map[int]map[string]int{}
	for _, n := range ngramsList {
		if n <= 1 {
			continue
		}
		ng := ngrams(tokens, n)
		if minFreq > 1 {
			for k, v := range ng {
				if v < minFreq { delete(ng, k) }
			}
		}
		ngramsByN[n] = ng
	}
	// TF-IDF simples
	var tfidf map[string]float64
	if tfidfOn {
		tfidf = computeTFIDFSingle(tokens)
	} else {
		tfidf = map[string]float64{}
	}

	// Serializar resultado
	out := buildLexicalJSON(freq, ngramsByN, tfidf)

	if step.Output == "" {
		return fmt.Errorf("lexical: output é obrigatório")
	}
	if err := e.ensureOutputDir(step.Output); err != nil {
		return err
	}
	if err := os.WriteFile(step.Output, []byte(out), 0644); err != nil {
		return fmt.Errorf("erro ao criar output: %w", err)
	}

	// Exportar CSVs opcionais
	if exportCSV {
		baseDir := filepath.Dir(step.Output)
		if len(freq) > 0 {
			if err := writeCSVCounts(filepath.Join(baseDir, "lexical_freq.csv"), freq); err != nil {
				return err
			}
		}
		for n, m := range ngramsByN {
			if len(m) == 0 { continue }
			if err := writeCSVCounts(filepath.Join(baseDir, fmt.Sprintf("lexical_ngrams_%d.csv", n)), m); err != nil {
				return err
			}
		}
	}
	return nil
}

// executeGraph executa geração de grafo
func (e *Executor) executeGraph(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Graph ainda não implementado (v0.1)")

	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		// Criar GEXF mínimo como placeholder
		gexf := `<?xml version="1.0" encoding="UTF-8"?>
<gexf xmlns="http://www.gexf.net/1.2draft" version="1.2">
  <graph mode="static" defaultedgetype="undirected">
    <nodes></nodes>
    <edges></edges>
  </graph>
</gexf>`
		if err := os.WriteFile(step.Output, []byte(gexf), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// executeTree executa geração de árvore semântica
func (e *Executor) executeTree(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Tree ainda não implementado (v0.1)")

	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		if err := os.WriteFile(step.Output, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// executeReport executa geração de relatório
func (e *Executor) executeReport(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Report ainda não implementado (v0.1)")

	if step.Output != "" {
		if err := e.ensureOutputDir(step.Output); err != nil {
			return err
		}
		// Criar HTML mínimo como placeholder
		html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>SapientIA Report</title>
</head>
<body>
    <h1>Relatório SapientIA</h1>
    <p>Relatório em desenvolvimento (v0.1)</p>
</body>
</html>`
		if err := os.WriteFile(step.Output, []byte(html), 0644); err != nil {
			return fmt.Errorf("erro ao criar output: %w", err)
		}
	}
	return nil
}

// ensureOutputDir garante que o diretório de output existe
func (e *Executor) ensureOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório: %w", err)
		}
	}
	return nil
}

// --- Auxiliares simples para lexical ---

type corpusDoc struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type corpusFile struct {
	Documents []corpusDoc `json:"documents"`
}

func extractFirstTextField(s string) string {
	var c corpusFile
	if err := json.Unmarshal([]byte(s), &c); err == nil {
		if len(c.Documents) > 0 {
			return c.Documents[0].Text
		}
	}
	// Fallback: tratar como texto puro
	return s
}

func simpleTokenizePT(s string) []string {
	s = strings.ToLower(s)
	// Substitui tudo que não for letra/numero por espaço
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	s = re.ReplaceAllString(s, " ")
	parts := strings.Fields(s)
	stop := ptStopwords()
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if len(p) <= 1 {
			continue
		}
		if stop[p] {
			continue
		}
		out = append(out, p)
	}
	return out
}

func ptStopwords() map[string]bool {
	words := []string{
		"a", "à", "às", "as", "o", "os", "um", "uma", "de", "do", "da", "dos", "das",
		"e", "é", "em", "no", "na", "nos", "nas", "para", "por", "com", "sem", "se",
		"que", "quem", "quando", "onde", "como", "porque", "porquê", "mas", "ou", "também",
		"mais", "menos", "muito", "muita", "muitos", "muitas", "já", "não", "sim", "ao",
	}
	m := make(map[string]bool, len(words))
	for _, w := range words {
		m[w] = true
	}
	return m
}

func countFreq(tokens []string) map[string]int {
	m := make(map[string]int)
	for _, t := range tokens {
		m[t]++
	}
	return m
}

func ngrams(tokens []string, n int) map[string]int {
	m := make(map[string]int)
	if n <= 1 {
		return m
	}
	for i := 0; i+n <= len(tokens); i++ {
		key := strings.Join(tokens[i:i+n], " ")
		m[key]++
	}
	return m
}

func computeTFIDFSingle(tokens []string) map[string]float64 {
	tf := make(map[string]float64)
	if len(tokens) == 0 {
		return tf
	}
	for _, t := range tokens {
		tf[t] += 1.0
	}
	L := float64(len(tokens))
	for k := range tf {
		tf[k] = tf[k] / L // idf=1 em corpus unitário
	}
	return tf
}

func buildLexicalJSON(freq map[string]int, ngramsByN map[int]map[string]int, tfidf map[string]float64) string {
	payload := map[string]any{
		"freq":  freq,
		"tfidf": tfidf,
	}
	for n, m := range ngramsByN {
		payload[fmt.Sprintf("ngrams_%d", n)] = m
	}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		// Fallback muito simples
		return "{}"
	}
	return string(b)
}

func simpleTokenizePTWithExternal(s string, stopwordsPath string) []string {
	base := ptStopwords()
	if stopwordsPath != "" {
		if b, err := os.ReadFile(stopwordsPath); err == nil {
			lines := strings.Split(string(b), "\n")
			for _, ln := range lines {
				w := strings.TrimSpace(strings.ToLower(ln))
				if w != "" { base[w] = true }
			}
		}
	}
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[^\p{L}\p{N}]+`)
	s = re.ReplaceAllString(s, " ")
	parts := strings.Fields(s)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if len(p) <= 1 { continue }
		if base[p] { continue }
		out = append(out, p)
	}
	return out
}

func writeCSVCounts(path string, m map[string]int) error {
	if err := eEnsureDir(path); err != nil { return err }
	var sb strings.Builder
	sb.WriteString("term,count\n")
	// deterministic order not required now; simple range
	for k, v := range m {
		sb.WriteString(fmt.Sprintf("%s,%d\n", k, v))
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func eEnsureDir(path string) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}
