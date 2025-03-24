package cli

import (
	"github.com/caner-cetin/seer/internal"
	"github.com/caner-cetin/seer/pkg/db"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type migrateConfig struct {
	LinguistLanguageRemotePath string
}

var (
	migrateCmd = &cobra.Command{
		Use: "migrate",
		Run: WrapCommandWithResources(migrate, ResourceConfig{Resources: []ResourceType{ResourceDatabase}}),
	}
	migrateCfg migrateConfig
)

func getMigrateCmd() *cobra.Command {
	migrateCmd.PersistentFlags().StringVar(
		&migrateCfg.LinguistLanguageRemotePath,
		"linguist-language-remote-path",
		"https://raw.githubusercontent.com/github/linguist/master/lib/linguist/languages.yml",
		"remote path to linguist languages.yml",
	)
	return migrateCmd
}

func migrate(cmd *cobra.Command, args []string) {
	app := GetApp(cmd).(internal.AppCtx)
	if err := db.Migrate(app.StdDB); err != nil {
		log.Error().Err(err).Msg("failed to migrate database schema")
		return
	}
	log.Info().Msg("migrated database schema")
	langCnt, err := app.DB.GetLanguageCount(cmd.Context())
	if err != nil {
		log.Error().Err(err).Msg("failed to get languages")
		return
	}
	if langCnt == 0 {
		var languages db.LanguagesNonPgtype
		languages.Load(migrateCfg.LinguistLanguageRemotePath)
		var languageSource = new(db.LanguageCopyFrom)
		languageSource.Languages = languages.ToPgType()
		if _, err := app.Conn.CopyFrom(
			cmd.Context(),
			[]string{"languages"},
			db.LanguageColumns(),
			languageSource,
		); err != nil {
			log.Error().Err(err).Msg("failed to insert language batch to languages table")
			return
		}
	}
}
