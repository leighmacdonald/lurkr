package cmd

import (
	"context"
	"github.com/leighmacdonald/lurkr/internal"
	"github.com/leighmacdonald/lurkr/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		internal.Start(ctx)
		defer internal.Stop()
		<-ctx.Done()
	},
}

func Execute() {
	cobra.OnInitialize(func() {
		if err := config.Read(""); err != nil {
			log.Errorf("Failed to load config: %v", err)
			os.Exit(1)
		}
	})
	if err := rootCmd.Execute(); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
}
