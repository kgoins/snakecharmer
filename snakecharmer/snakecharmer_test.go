package snakecharmer_test

import (
	"testing"

	"github.com/spf13/cobra"
)

func init() {
}

func TestConfImport(t *testing.T) {
	testCmd := &cobra.Command{
		Use: "testcmd",
	}

	testCmd.PersistentFlags().StringP("server", "s", "", "server")
	testCmd.PersistentFlags().IntP("port", "p", 443, "port")
}
