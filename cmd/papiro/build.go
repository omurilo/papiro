package main

import (
	"fmt"
	"os"

	"github.com/omurilo/papiro/internal/builder"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Gera o site estático à partir do markdown",
	Long:  "Lê os arquivos da pasta 'content', processa o markdown e gera o HTML final na pasta 'public'.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Construindo o seu site...")

		err := builder.BuildSite()
		if err != nil {
			fmt.Printf("Falha na construção: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Site gerado com sucesso na pasta '/public'!")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
