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
	fmt.Printf("   ✅ Concluído em %s\n\n", formatElapsed(elapsed))
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
	inputPath := e.resolvePath(step.Input)
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("ingest: erro ao ler input '%s': %w", inputPath, err)
	}

	// Criar corpus JSON mínimo
	corpus := fmt.Sprintf(`{"documents":[{"id":"doc1","text":%q}]}`, string(data))

	// Escrever output
	if step.Output == "" {
		return fmt.Errorf("ingest: output é obrigatório")
	}
	outputPath := e.resolvePath(step.Output)
	if err := e.ensureOutputDir(outputPath); err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, []byte(corpus), 0644); err != nil {
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
	inputPath := e.resolvePath(step.Input)
	b, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("lexical: erro ao ler input '%s': %w", inputPath, err)
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
	outputPath := e.resolvePath(step.Output)
	if err := e.ensureOutputDir(outputPath); err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, []byte(out), 0644); err != nil {
		return fmt.Errorf("erro ao criar output: %w", err)
	}

	// Exportar CSVs opcionais
	if exportCSV {
		baseDir := filepath.Dir(outputPath)
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

	// Parâmetros
	window := 2
	minEdge := 1
	exportCSV := false
	corpusPath := ""
	stopwordsFile := ""
	if step.Params != nil {
		if v, ok := step.Params["cooccurrence_window"].(int); ok { window = v } else if fv, ok := step.Params["cooccurrence_window"].(float64); ok { window = int(fv) }
		if v, ok := step.Params["min_edge_weight"].(int); ok { minEdge = v } else if fv, ok := step.Params["min_edge_weight"].(float64); ok { minEdge = int(fv) }
		if v, ok := step.Params["export_csv"].(bool); ok { exportCSV = v }
		if v, ok := step.Params["corpus"].(string); ok { corpusPath = v }
		if v, ok := step.Params["stopwords_file"].(string); ok { stopwordsFile = v }
	}
	if window < 2 { window = 2 }
	if minEdge < 1 { minEdge = 1 }

	inputPath := e.resolvePath(step.Input)
	outputPath := e.resolvePath(step.Output)
	
	// Descobrir corpus.json se não fornecido
	if corpusPath == "" {
		baseDir := filepath.Dir(inputPath)
		guess := filepath.Join(baseDir, "corpus.json")
		if _, err := os.Stat(guess); err == nil {
			corpusPath = guess
		}
	} else {
		corpusPath = e.resolvePath(corpusPath)
	}

	// Ler tokens do corpus ou, como fallback, derivar arestas de ngrams_2 do lexical
	var tokens []string
	var freq map[string]int
	if corpusPath != "" {
		cb, err := os.ReadFile(corpusPath)
		if err != nil {
			return fmt.Errorf("graph: erro lendo corpus '%s': %w", corpusPath, err)
		}
		text := extractFirstTextField(string(cb))
		tokens = simpleTokenizePTWithExternal(text, stopwordsFile)
		freq = countFreq(tokens)
	} else {
		// Fallback: usar lexical.json para nós (freq) e arestas (ngrams_2) se existir
		lb, err := os.ReadFile(inputPath)
		if err != nil { return fmt.Errorf("graph: erro lendo lexical '%s': %w", inputPath, err) }
		var parsed map[string]any
		if err := json.Unmarshal(lb, &parsed); err != nil {
			return fmt.Errorf("graph: erro parse lexical.json: %w", err)
		}
		// freq
		freq = map[string]int{}
		if f, ok := parsed["freq"].(map[string]any); ok {
			for k, v := range f {
				switch t := v.(type) {
				case float64: freq[k] = int(t)
				case int: freq[k] = t
				}
			}
		}
		// tokens vazio; usaremos ngrams_2 como arestas
	}

	// Coocorrência
	edges := map[[2]string]int{}
	if len(tokens) > 0 {
		for i := 0; i < len(tokens); i++ {
			for d := 1; d < window && i+d < len(tokens); d++ {
				a, b := tokens[i], tokens[i+d]
				if a == b { continue }
				// aresta não-direcionada com ordenação estável
				if a > b { a, b = b, a }
				edges[[2]string{a, b}]++
			}
		}
	} else {
		// Fallback: usar ngrams_2 como arestas
		lb, err := os.ReadFile(inputPath)
		if err != nil { return fmt.Errorf("graph: erro lendo lexical '%s': %w", inputPath, err) }
		var parsed map[string]any
		if err := json.Unmarshal(lb, &parsed); err != nil {
			return fmt.Errorf("graph: erro parse lexical.json: %w", err)
		}
		if ng2, ok := parsed["ngrams_2"].(map[string]any); ok {
			for k, v := range ng2 {
				parts := strings.Split(k, " ")
				if len(parts) != 2 { continue }
				a, b := parts[0], parts[1]
				if a > b { a, b = b, a }
				w := 0
				switch t := v.(type) {
				case float64: w = int(t)
				case int: w = t
				}
				if w > 0 {
					edges[[2]string{a, b}] += w
				}
			}
		}
	}

	// Filtrar por min_edge_weight
	for k, w := range edges {
		if w < minEdge {
			delete(edges, k)
		}
	}

	// Geração de IDs de nós
	nodeID := map[string]int{}
	nodes := []string{}
	nextID := 0
	// Inicializar nós a partir das arestas
	for pair := range edges {
		for _, term := range pair {
			if _, ok := nodeID[term]; !ok {
				nodeID[term] = nextID
				nodes = append(nodes, term)
				nextID++
			}
		}
	}

	// Criar diretório de output
	if err := e.ensureOutputDir(outputPath); err != nil { return err }

	// Escrever GEXF
	gexf := buildGEXF(nodeID, nodes, edges, freq)
	if err := os.WriteFile(outputPath, []byte(gexf), 0644); err != nil {
		return fmt.Errorf("erro ao escrever GEXF: %w", err)
	}

	// CSVs opcionais
	if exportCSV {
		baseDir := filepath.Dir(outputPath)
		if err := writeEdgesCSV(filepath.Join(baseDir, "graph_edges.csv"), edges); err != nil { return err }
		if err := writeNodesCSV(filepath.Join(baseDir, "graph_nodes.csv"), nodeID, freq); err != nil { return err }
	}
	return nil
}

// executeTree executa geração de árvore semântica
func (e *Executor) executeTree(step *Step) error {
	fmt.Printf("   📥 Input: %s\n", step.Input)
	fmt.Printf("   📤 Output: %s\n", step.Output)
	fmt.Println("   ⚠️  Tree ainda não implementado (v0.1)")

	if step.Output != "" {
		outputPath := e.resolvePath(step.Output)
		if err := e.ensureOutputDir(outputPath); err != nil {
			return err
		}
		if err := os.WriteFile(outputPath, []byte("{}"), 0644); err != nil {
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
		outputPath := e.resolvePath(step.Output)
		if err := e.ensureOutputDir(outputPath); err != nil {
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
		if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
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

// resolvePath resolve paths relativos ao workDir
func (e *Executor) resolvePath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(e.workDir, p)
}

// formatElapsed retorna uma string amigável com maior precisão temporal
// - < 1ms: microsegundos
// - < 1s: milissegundos com 2 casas
// - >= 1s: segundos com 3 casas
func formatElapsed(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		us := float64(d) / float64(time.Microsecond)
		return fmt.Sprintf("%.2fµs", us)
	}
	if d < time.Second {
		ms := float64(d) / float64(time.Millisecond)
		return fmt.Sprintf("%.2fms", ms)
	}
	return fmt.Sprintf("%.3fs", d.Seconds())
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
		// Filtrar números puros (cantos, versos, marcadores)
		if isNumericOnly(p) {
			continue
		}
		if stop[p] {
			continue
		}
		out = append(out, p)
	}
	return out
}

func isNumericOnly(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
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
		if len(p) <= 1 {
			continue
		}
		// Filtrar números puros
		if isNumericOnly(p) {
			continue
		}
		if base[p] {
			continue
		}
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

// --- Graph helpers ---
func buildGEXF(nodeID map[string]int, nodes []string, edges map[[2]string]int, freq map[string]int) string {
	var sb strings.Builder
	sb.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	sb.WriteString("<gexf xmlns=\"http://www.gexf.net/1.2draft\" version=\"1.2\">\n")
	sb.WriteString("  <graph mode=\"static\" defaultedgetype=\"undirected\">\n")
	sb.WriteString("    <nodes>\n")
	for _, term := range nodes {
		id := nodeID[term]
		w := freq[term]
		sb.WriteString(fmt.Sprintf("      <node id=\"%d\" label=\"%s\">", id, xmlEscape(term)))
		if w > 0 {
			sb.WriteString("<attvalues><attvalue for=\"weight\" value=\"" + fmt.Sprint(w) + "\"/></attvalues>")
		}
		sb.WriteString("</node>\n")
	}
	sb.WriteString("    </nodes>\n")
	sb.WriteString("    <edges>\n")
	eid := 0
	for pair, w := range edges {
		src := nodeID[pair[0]]
		dst := nodeID[pair[1]]
		sb.WriteString(fmt.Sprintf("      <edge id=\"%d\" source=\"%d\" target=\"%d\" weight=\"%d\"/>\n", eid, src, dst, w))
		eid++
	}
	sb.WriteString("    </edges>\n")
	sb.WriteString("  </graph>\n")
	sb.WriteString("</gexf>\n")
	return sb.String()
}

func writeEdgesCSV(path string, edges map[[2]string]int) error {
	if err := eEnsureDir(path); err != nil { return err }
	var sb strings.Builder
	sb.WriteString("source,target,weight\n")
	for pair, w := range edges {
		sb.WriteString(fmt.Sprintf("%s,%s,%d\n", pair[0], pair[1], w))
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func writeNodesCSV(path string, nodeID map[string]int, freq map[string]int) error {
	if err := eEnsureDir(path); err != nil { return err }
	var sb strings.Builder
	sb.WriteString("id,label,weight\n")
	for term, id := range nodeID {
		w := freq[term]
		sb.WriteString(fmt.Sprintf("%d,%s,%d\n", id, term, w))
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func xmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
