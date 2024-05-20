package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gkwa/eachdodge/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/goldbug"
)

var (
	cfgFile   string
	verbose   bool
	logFormat string
	outfile   string
	out       string
)

var rootCmd = &cobra.Command{
	Use:   "eachdodge",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		core.Run(outfile, out)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eachdodge.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		slog.Error("error binding verbose flag", "error", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "json or text (default is text)")
	err = viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format"))
	if err != nil {
		slog.Error("error binding log-format flag", "error", err)
		os.Exit(1)
	}

	rootCmd.Flags().StringVar(&outfile, "outfile", "", "Output file for IP addresses")
	err = viper.BindPFlag("outfile", rootCmd.Flags().Lookup("outfile"))
	if err != nil {
		slog.Error("error binding outfile flag", "error", err)
		os.Exit(1)
	}

	rootCmd.Flags().StringVar(&out, "out", "json", "Output format (list or json)")
	err = viper.BindPFlag("out", rootCmd.Flags().Lookup("out"))
	if err != nil {
		slog.Error("error binding out flag", "error", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".eachdodge")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	logFormat = viper.GetString("log-format")
	verbose = viper.GetBool("verbose")

	slog.Debug("using config file", "path", viper.ConfigFileUsed())
	slog.Debug("log-format", "value", logFormat)
	slog.Debug("log-format", "value", viper.GetString("log-format"))

	setupLogging()
}

func setupLogging() {
	if verbose || logFormat != "" {
		if logFormat == "json" {
			goldbug.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			goldbug.SetDefaultLoggerText(slog.LevelDebug)
		}

		slog.Debug("setup", "verbose", verbose)
	}
}
