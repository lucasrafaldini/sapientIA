package pipeline

import (
	"reflect"
	"testing"
)

func TestSimpleTokenizePT(t *testing.T) {
	text := "Olá, mundo! Isto é um teste simples de tokenização em PT-BR."
	tokens := simpleTokenizePT(text)
	if len(tokens) == 0 {
		t.Fatalf("esperava tokens, obtive 0")
	}
	for _, tk := range tokens {
		if len(tk) <= 1 {
			t.Errorf("token muito curto: %q", tk)
		}
	}
}

func TestCountFreq(t *testing.T) {
	tokens := []string{"a", "b", "a", "c", "b", "a"}
	m := countFreq(tokens)
	if m["a"] != 3 || m["b"] != 2 || m["c"] != 1 {
		t.Errorf("freq incorreta: %#v", m)
	}
}

func TestNGrams(t *testing.T) {
	tokens := []string{"go", "é", "legal"}
	m := ngrams(tokens, 2)
	exp := map[string]int{"go é": 1, "é legal": 1}
	if !reflect.DeepEqual(m, exp) {
		t.Errorf("ngrams inválido: got=%#v exp=%#v", m, exp)
	}
}

func TestTFIDFSingle(t *testing.T) {
	tokens := []string{"um", "dois", "um"}
	m := computeTFIDFSingle(tokens)
	if len(m) != 2 {
		t.Fatalf("esperava 2 termos, got=%d", len(m))
	}
	if m["um"] <= m["dois"] {
		t.Errorf("tf esperado de 'um' > 'dois': %#v", m)
	}
}
