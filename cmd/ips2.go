/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/gkwa/eachdodge/ips2"
	"github.com/spf13/cobra"
)

var outfile string

var ips2Cmd = &cobra.Command{
	Use:   "ips2",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ips2.IPs2(outfile)
	},
}

func init() {
	rootCmd.AddCommand(ips2Cmd)
	ips2Cmd.Flags().StringVar(&outfile, "outfile", "ips.json", "Output file for IP addresses")
}
