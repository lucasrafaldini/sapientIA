package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "0.1.0-dev"
	rootCmd = &cobra.Command{
		Use:   "sapientia",
		Short: "SapientIA - Análise de conversas e dados textuais",
		Long: `SapientIA é uma plataforma para análise de conversas corporativas.
Transforma briefs e entrevistas em árvores semânticas, grafos e insights acionáveis.`,
		Version: version,
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
