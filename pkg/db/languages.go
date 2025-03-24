package db

import (
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// LanguageNonPgtype stores programming language definitions and metadata
type LanguageNonPgtype struct {
	// Primary name of the language
	Name string `yaml:"name"`
	// Optional field. Only necessary as a replacement for the sample directory name if the language name is not a valid filename
	FsName string `yaml:"fs_name"`
	// Category of the language: data, programming, markup, prose, or null
	Type string `yaml:"type"`
	// An array of additional aliases (implicitly includes name.downcase)
	Aliases []string `yaml:"aliases"`
	// A String name of the Ace Mode used for highlighting whenever a file is edited. This must match one of the filenames in https://gh.io/acemodes. Use "text" if a mode does not exist
	AceMode string `yaml:"ace_mode"`
	// A String name of the CodeMirror Mode used for highlighting whenever a file is edited. This must match a mode from https://git.io/vi9Fx
	CodemirrorMode string `yaml:"codemirror_mode"`
	// A String name of the file mime type used for highlighting whenever a file is edited. This should match the `mime` associated with the mode from https://git.io/f4SoQ
	CodemirrorMimeType string `yaml:"codemirror_mime_type"`
	// Boolean value to enable line wrapping (default: false)
	Wrap bool `yaml:"wrap"`
	// An array of associated extensions (the first one is considered the primary extension, the others should be listed alphabetically)
	Extensions []string `yaml:"extensions"`
	// An array of filenames commonly associated with the language
	Filenames []string `yaml:"filenames"`
	// An array of associated interpreters
	Interpreters []string `yaml:"interpreters"`
	// Integer used as a language-name-independent indexed field so that we can rename languages in Linguist without reindexing all the code on GitHub
	LanguageID int32 `yaml:"language_id"`
	// CSS hex color to represent the language. Only used if type is "programming" or "markup"
	Color string `yaml:"color"`
	// The TextMate scope that represents this programming language. This should match one of the scopes listed in the grammars.yml file. Use "none" if there is no grammar for this language
	TmScope string `yaml:"tm_scope"`
	// Name of the parent language. Languages in a group are counted in the statistics as the parent language
	Group string `yaml:"group"`
}

// LanguagesNonPgtype represents a collection of programming language definitions, without pgtypes.
type LanguagesNonPgtype []LanguageNonPgtype

// Load fetches programming language definitions from a remote YAML file specified by remote_path
// https://raw.githubusercontent.com/github-linguist/linguist/refs/heads/main/lib/linguist/languages.yml
// should be used
func (l *LanguagesNonPgtype) Load(remote_path string) {
	req, err := http.NewRequest(http.MethodGet, remote_path, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get linguist languages")
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}()
	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return
	}
	var languageYaml map[string]interface{}
	if err := yaml.Unmarshal(resp_bytes, &languageYaml); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal language config into map")
		return
	}
	for k, v := range languageYaml {
		var language LanguageNonPgtype
		mv, err := yaml.Marshal(v)
		if err != nil {
			log.Error().Str("key", k).Err(err).Msg("failed to marshal language config")
			return
		}
		if err := yaml.Unmarshal(mv, &language); err != nil {
			log.Error().Str("key", k).Str("value", string(mv)).Err(err).Msg("failed to unmarshal language config back to struct")
			return
		}
		language.Name = k
		*l = append(*l, language)
	}
}

// ToPgType converts a slice of LanguageNonPgtype to a slice of Language with PostgreSQL-compatible types
func (l LanguagesNonPgtype) ToPgType() []Language {
	var dbLanguages = make([]Language, 0, len(l))
	for _, lang := range l {
		var dbLanguage = Language{
			Name:         lang.Name,
			Aliases:      lang.Aliases,
			Extensions:   lang.Extensions,
			Filenames:    lang.Filenames,
			Interpreters: lang.Filenames,
			LanguageID:   lang.LanguageID,
			Wrap:         pgtype.Bool{Bool: lang.Wrap, Valid: true},
		}
		if lang.FsName != "" {
			dbLanguage.FsName = pgtype.Text{String: lang.FsName, Valid: true}
		} else {
			dbLanguage.FsName = pgtype.Text{String: "", Valid: true}
		}
		if lang.Type != "" {
			var ltype LanguageType
			switch lang.Type {
			case "data":
				ltype = LanguageTypeData
			case "programming":
				ltype = LanguageTypeData
			case "markup":
				ltype = LanguageTypeMarkup
			case "prose":
				ltype = LanguageTypeProse
			}
			dbLanguage.Type = NullLanguageType{LanguageType: ltype, Valid: true}
		} else {
			dbLanguage.Type = NullLanguageType{LanguageType: LanguageTypeData, Valid: false}
		}
		if lang.AceMode != "" {
			dbLanguage.AceMode = pgtype.Text{String: lang.AceMode, Valid: true}
		} else {
			dbLanguage.AceMode = pgtype.Text{String: "", Valid: false}
		}
		if lang.CodemirrorMode != "" {
			dbLanguage.CodemirrorMode = pgtype.Text{String: lang.CodemirrorMode, Valid: true}
		} else {
			dbLanguage.CodemirrorMode = pgtype.Text{String: "", Valid: false}
		}
		if lang.CodemirrorMimeType != "" {
			dbLanguage.CodemirrorMimeType = pgtype.Text{String: lang.CodemirrorMimeType, Valid: true}
		} else {
			dbLanguage.CodemirrorMimeType = pgtype.Text{String: "", Valid: false}
		}
		if lang.Color != "" {
			dbLanguage.Color = pgtype.Text{String: lang.Color, Valid: true}
		} else {
			dbLanguage.Color = pgtype.Text{String: "", Valid: false}
		}
		if lang.Group != "" {
			dbLanguage.Group = pgtype.Text{String: lang.Group, Valid: true}
		} else {
			dbLanguage.Group = pgtype.Text{String: "", Valid: false}
		}
		if lang.TmScope != "" {
			dbLanguage.TmScope = pgtype.Text{String: lang.TmScope, Valid: true}
		} else {
			dbLanguage.TmScope = pgtype.Text{String: "", Valid: false}
		}
		dbLanguages = append(dbLanguages, dbLanguage)
	}
	return dbLanguages
}

func LanguageColumns() []string {
	return []string{
		"id",
		"name",
		"fs_name",
		"type",
		"aliases",
		"ace_mode",
		"codemirror_mode",
		"codemirror_mime_type",
		"wrap",
		"extensions",
		"filenames",
		"interpreters",
		"language_id",
		"color",
		"tm_scope",
		"group",
	}
}

// MarshalZerologObject implements zerolog.LogObjectMarshaler interface to enable structured logging of Language objects
func (l Language) MarshalZerologObject(e *zerolog.Event) {
	e.Int32("id", l.ID).
		Str("name", l.Name).
		Str("fs_name", l.FsName.String).
		Str("type", string(l.Type.LanguageType)).
		Strs("aliases", l.Aliases).
		Str("ace_mode", l.AceMode.String).
		Str("codemirror_mode", l.CodemirrorMode.String).
		Str("codemirror_mime_type", l.CodemirrorMimeType.String).
		Bool("wrap", l.Wrap.Bool).
		Strs("extensions", l.Extensions).
		Strs("filenames", l.Filenames).
		Strs("interpreters", l.Interpreters).
		Int32("language_id", l.LanguageID).
		Str("color", l.Color.String).
		Str("tm_scope", l.TmScope.String).
		Str("group", l.Group.String)
}

// LanguageCopyFrom implements pgx.CopyFromSource interface for bulk inserting Language records
type LanguageCopyFrom struct {
	Languages []Language
	i         int
}

// Next returns true if there is another row and makes the next row data
// available to Values(). When there are no more rows available or an error
// has occurred it returns false.
func (l *LanguageCopyFrom) Next() bool {
	l.i++
	log.Trace().Int("i", l.i).Msg("incrementing LanguageCopyFrom index")
	return l.i < len(l.Languages)
}

// Values returns the values for the current row.
func (l *LanguageCopyFrom) Values() ([]any, error) {
	lang := l.Languages[l.i]
	log.Trace().Int("i", l.i).Object("language", lang).Msg("yielding from LanguageCopyFrom.Values")
	if l.i >= len(l.Languages) {
		return nil, nil
	}
	if l.i > math.MaxInt32 {
		return nil, fmt.Errorf("index %d exceeds maximum value for int32", l.i)
	}
	lang.ID = int32(l.i - 1) //nolint: gosec
	return []any{
		lang.ID,
		lang.Name,
		lang.FsName,
		lang.Type,
		lang.Aliases,
		lang.AceMode,
		lang.CodemirrorMode,
		lang.CodemirrorMimeType,
		lang.Wrap,
		lang.Extensions,
		lang.Filenames,
		lang.Interpreters,
		lang.LanguageID,
		lang.Color,
		lang.TmScope,
		lang.Group,
	}, nil
}

// Err returns any error that has been encountered by the CopyFromSource. If
// this is not nil *Conn.CopyFrom will abort the copy.
func (l *LanguageCopyFrom) Err() error {
	return nil
}
