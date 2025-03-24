package cli

import (
	"context"
	"time"

	"github.com/caner-cetin/seer/internal"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// ResourceType defines the type of resource to be initialized
type ResourceType int

// Resource types for initialization
const (
	// ResourceDatabase represents database resource type
	ResourceDatabase ResourceType = iota
)

// ResourceConfig holds the configuration for resource initialization
type ResourceConfig struct {
	Resources []ResourceType
}

// WrapCommandWithResources wraps a Cobra command function with resource initialization and cleanup logic based on the provided configuration
func WrapCommandWithResources(fn func(cmd *cobra.Command, args []string), config ResourceConfig) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMs*int(time.Millisecond)))
		defer cancel()

		appCtx := internal.AppCtx{}

		for _, resource := range config.Resources {
			switch resource {
			case ResourceDatabase:
				cfg.SetDBURL()
				if err := appCtx.InitializeDB(); err != nil {
					log.Error().Err(err).Msg("failed to initialize database")
					return
				}
			}

		}
		defer func() {
			if appCtx.Conn != nil {
				if err := appCtx.Conn.Close(ctx); err != nil {
					log.Error().Err(err).Msg("failed to close database connection")
					return
				}
			}
		}()
		cmd.SetContext(context.WithValue(ctx, internal.APP_CONTEXT_KEY, appCtx))
		fn(cmd, args)
	}
}

func GetApp(cmd *cobra.Command) interface{} {
	return cmd.Context().Value(internal.APP_CONTEXT_KEY)
}
