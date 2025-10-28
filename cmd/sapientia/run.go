package main

import (
	"fmt"

	"github.com/lucasrafaldini/sapientIA/internal/pipeline"
	"github.com/spf13/cobra"
)

var (
	stopwordsProfile string
)

var runCmd = &cobra.Command{
	Use:   "run [pipeline.yaml]",
	Short: "Executa um pipeline completo",
	Long: `Executa um pipeline YAML completo de análise.
	
Exemplo:
  sapientia run pipelines/exemplo.yaml
  sapientia run pipelines/exemplo.yaml --stopwords extended`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pipelinePath := args[0]

		// Carregar pipeline
		p, err := pipeline.LoadPipeline(pipelinePath)
		if err != nil {
			return fmt.Errorf("erro ao carregar pipeline: %w", err)
		}

		// Criar executor
		executor := pipeline.NewExecutor(p)

		// Aplicar stopwords profile da CLI se fornecido
		if stopwordsProfile != "" {
			executor.SetStopwordsProfile(stopwordsProfile)
		}

		// Executar
		if err := executor.Execute(); err != nil {
			return fmt.Errorf("erro na execução: %w", err)
		}

		return nil
	},
}

func init() {
	runCmd.Flags().StringVar(&stopwordsProfile, "stopwords", "core", "Perfil de stopwords: core, extended, aggressive (default: core)")
	rootCmd.AddCommand(runCmd)
}
