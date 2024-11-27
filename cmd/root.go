package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"log/slog"
	"os"
	"os/signal"
)

func Start() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	rootCmd := &cobra.Command{}
	commands := []*cobra.Command{
		{
			Use:   "serve-http",
			Short: "Run HTTP server",
			Run: func(cmd *cobra.Command, args []string) {
				runHTTPServer(ctx)
			},
		},
	}
	rootCmd.AddCommand(commands...)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatalln(err)
	}

}
