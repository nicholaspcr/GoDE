package commands

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/nicholaspcr/GoDE/internal/store"
	"github.com/nicholaspcr/GoDE/pkg/api"
	"github.com/spf13/cobra"
)

// testCmd set of commands used to test specific parts of the application
// locally.
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A set of commands to test the application locally",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info(args)
		return cmd.Help()
	},
}

var dbServerCmd = &cobra.Command{
	Use:   "db-server",
	Short: "Starts a local database server",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info(args)

		ctx := cmd.Context()
		r := gin.Default()
		db, err := store.New(ctx)
		if err != nil {
			return err
		}

		r.POST("/user", func(c *gin.Context) {
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
			var usr api.User
			if err := json.Unmarshal(b, &usr); err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			err = db.Create(ctx, &usr)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
			}

			c.JSON(200, gin.H{})
		})

		r.GET("/user", func(c *gin.Context) {
			usr, err := db.Read(ctx, &api.UserID{
				Id: c.Query("id"),
			})
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, usr)
		})

		r.PUT("/user", func(c *gin.Context) {
			b, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}
			var usr api.User
			if err := json.Unmarshal(b, &usr); err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			err = db.Update(ctx, &usr)
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
			}

			c.JSON(200, gin.H{})
		})

		r.DELETE("/user", func(c *gin.Context) {
			err := db.Delete(ctx, &api.UserID{
				Id: c.Query("id"),
			})
			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(200, gin.H{})
		})

		return r.Run()
	},
}
