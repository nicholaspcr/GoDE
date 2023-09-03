package commands

import (
	"fmt"
	"log/slog"
	_ "net/http/pprof"

	"github.com/spf13/cobra"
)

var format string

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Shows web server configuration",
	RunE: func(*cobra.Command, []string) error {
		var b []byte
		var err error
		switch format {
		case "json":
			b, err = cfg.JSON()
		case "yaml", "yml":
			b, err = cfg.YAML()
		default:
			b, err = cfg.JSON()
		}

		if err != nil {
			slog.Error(
				"Failed to parse config to desired format",
				slog.String("format", format),
				slog.String("error_msg", err.Error()),
			)
			return err
		}
		_, err = fmt.Println(string(b))
		return err
	},
}

func init() {
	configCmd.Flags().StringVarP(
		&format, "format", "f", "json",
		"Format to print config in. Options: json and yaml",
	)
}
