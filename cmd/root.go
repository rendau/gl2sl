package cmd

import (
	"context"
	"github.com/rendau/gl2sl/internal/adapters/http_api"
	"github.com/rendau/gl2sl/internal/domain/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "gl2sl",
	Run: func(cmd *cobra.Command, args []string) {
		cr := core.NewSt(viper.GetString("slack_webhook_url"))

		httpApi := http_api.NewApi(
			viper.GetString("http_listen"),
			cr,
		)

		go func() {
			err := httpApi.Start()
			if err != nil {
				log.Fatal(err)
			}
		}()

		log.Println("Started", "http_listen", viper.GetString("http_listen"))

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		log.Println("Shutting down...")

		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
		err := httpApi.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetDefault("http_listen", ":80")

	viper.AutomaticEnv() // read in environment variables that match

	if viper.GetString("conf_path") != "" {
		viper.SetConfigFile(viper.GetString("conf_path"))
		_ = viper.ReadInConfig()
	}
}
