package snakecharmer_test

import (
	"path"
	"runtime"
	"testing"

	"github.com/kgoins/snakecharmer/snakecharmer"
	"github.com/spf13/cobra"
)

func getTestDataDir() string {
	_, me, _, _ := runtime.Caller(0)
	myParent := path.Join(path.Dir(me))

	return path.Join(
		myParent,
		"testdata",
	)
}

func TestConfImport(t *testing.T) {
	confPath := path.Join(getTestDataDir(), ".sctest.toml")
	sc := snakecharmer.NewSnakeCharmer("sctest", ".sctest")

	testCmd := &cobra.Command{
		Use: "testcmd",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return sc.InitConfig(cmd, confPath)
		},
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt("port")
			server, _ := cmd.Flags().GetString("server")

			if port != 1337 || server != "example.com" {
				t.Fatal()
			}
		},
	}

	testCmd.PersistentFlags().StringP("server", "s", "", "server")
	testCmd.PersistentFlags().IntP("port", "p", 443, "port")

	err := testCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}
}
