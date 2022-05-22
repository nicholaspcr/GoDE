package commands

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)


var startCmd = &cobra.Command{
	Use: "start",
	Short: "Start the server responsible for handling requests to start DE instances.",
	RunE: func(_ *cobra.Command, _ []string) error {
		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		return r.Run()
	},
}
