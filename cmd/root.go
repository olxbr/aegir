package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:     "aegir",
	Short:   "Aegir is a generic admission controller for Kubernetes",
	Long:    `A generic admission controller to validate Kubernetes resources using LIVR rules.`,
	Version: "0.1.0",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
