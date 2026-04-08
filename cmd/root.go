package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/devsh/internal/config"
	"github.com/devsh/internal/plugins"
	"github.com/devsh/internal/repl"
	"github.com/devsh/internal/router"
	"github.com/spf13/cobra"
)

var (
	interactive bool
	dryRun      bool
	force       bool
	explain     bool
	provider    string
)

var rootCmd = &cobra.Command{
	Use:   "devsh [natural language command]",
	Short: "devsh is a natural language terminal assistant",
	Long:  "devsh translates natural language into executable, safe terminal commands.",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Warning: Failed to load config: %v\n", err)
		}

		pluginsList := []plugins.Plugin{
			&plugins.BashPlugin{},
			&plugins.KubernetesPlugin{},
			&plugins.AWSPlugin{},
			&plugins.TerraformPlugin{},
			&plugins.GitPlugin{},
			&plugins.DockerPlugin{},
		}

		r, err := router.NewRouter(pluginsList)
		if err != nil {
			fmt.Printf("Failed to initialize devsh: %v\n", err)
			os.Exit(1)
		}

		ctx := context.Background()

		if interactive {
			session := repl.NewSession(r)
			if err := session.Start(ctx, force); err != nil {
				fmt.Printf("REPL error: %v\n", err)
				os.Exit(1)
			}
			return
		}

		if len(args) == 0 {
			cmd.Help()
			return
		}

		input := strings.Join(args, " ")
		
		_, err = r.Process(ctx, input, nil, force, dryRun)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Start REPL mode")
	rootCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Generate but do not execute command")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Bypass safety confirmations (USE WITH CAUTION)")
	rootCmd.Flags().BoolVarP(&explain, "explain", "e", false, "Provide detailed explanation of the generated command")
	rootCmd.Flags().StringVarP(&provider, "provider", "p", "", "Force a specific tool/provider (aws, k8s, terraform, etc.)")

	rootCmd.AddCommand(configCmd)
}
