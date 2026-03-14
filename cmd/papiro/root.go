package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "papiro",
	Short: "Papiro é um gerador de sites estáticos super rápido",
	Long: "Papiro é um gerador de sites estáticos à partir de papiros (markdown files) focado em simplicidade, velocidade e minimalismo. Transforme markdown em html numa rabiscada de tinta no papel.",
	Run: func (cmd *cobra.Command, args []string) {
		fmt.Println("Bem vindo ao Papiro!")
		fmt.Println("Use 'papiro --help' para ver os comandos disponíveis.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}