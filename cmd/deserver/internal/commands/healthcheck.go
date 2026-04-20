package commands

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "probe the server's /health endpoint and exit 0 on success",
	Long: `healthcheck issues a single HTTP GET to /health and exits with
status 0 if the server returns 2xx, status 1 otherwise. Designed for
use as a container healthcheck on distroless runtimes where curl/wget
are unavailable.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(*cobra.Command, []string) error {
		// Skip root's config/log setup — healthcheck must survive misconfig.
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		port := strings.TrimPrefix(os.Getenv("HTTP_PORT"), ":")
		if port == "" {
			port = "8081"
		}
		url := fmt.Sprintf("http://127.0.0.1:%s/health", port)

		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("healthcheck request failed: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("healthcheck got status %d", resp.StatusCode)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)
}
