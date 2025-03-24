-- +goose Up
-- +goose StatementBegin
CREATE TYPE language_type AS ENUM ('data', 'programming', 'markup', 'prose');
CREATE TABLE public.languages (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  fs_name VARCHAR(255),
  "type" language_type,
  aliases TEXT [],
  ace_mode VARCHAR(100),
  codemirror_mode VARCHAR(100),
  codemirror_mime_type VARCHAR(100),
  wrap BOOLEAN DEFAULT false,
  extensions TEXT [],
  filenames TEXT [],
  interpreters TEXT [],
  language_id INTEGER UNIQUE NOT NULL,
  color CHAR(7),
  tm_scope VARCHAR(255),
  "group" VARCHAR(255)
);
COMMENT ON TABLE languages IS 'Stores programming language definitions and metadata';
COMMENT ON COLUMN languages.name IS 'Primary name of the language';
COMMENT ON COLUMN languages.fs_name IS 'Optional field. Only necessary as a replacement for the sample directory name if the language name is not a valid filename under the Windows filesystem (e.g., if it contains an asterisk)';
COMMENT ON COLUMN languages.type IS 'Category of the language: data, programming, markup, prose, or null';
COMMENT ON COLUMN languages.aliases IS 'An array of additional aliases (implicitly includes name.downcase)';
COMMENT ON COLUMN languages.ace_mode IS 'A String name of the Ace Mode used for highlighting whenever a file is edited. This must match one of the filenames in https://gh.io/acemodes. Use "text" if a mode does not exist';
COMMENT ON COLUMN languages.codemirror_mode IS 'A String name of the CodeMirror Mode used for highlighting whenever a file is edited. This must match a mode from https://git.io/vi9Fx';
COMMENT ON COLUMN languages.codemirror_mime_type IS 'A String name of the file mime type used for highlighting whenever a file is edited. This should match the `mime` associated with the mode from https://git.io/f4SoQ';
COMMENT ON COLUMN languages.wrap IS 'Boolean value to enable line wrapping (default: false)';
COMMENT ON COLUMN languages.extensions IS 'An array of associated extensions (the first one is considered the primary extension, the others should be listed alphabetically)';
COMMENT ON COLUMN languages.filenames IS 'An array of filenames commonly associated with the language';
COMMENT ON COLUMN languages.interpreters IS 'An array of associated interpreters';
COMMENT ON COLUMN languages.language_id IS 'Integer used as a language-name-independent indexed field so that we can rename languages in Linguist without reindexing all the code on GitHub. Must not be changed for existing languages without the explicit permission of GitHub staff';
COMMENT ON COLUMN languages.color IS 'CSS hex color to represent the language. Only used if type is "programming" or "markup"';
COMMENT ON COLUMN languages.tm_scope IS 'The TextMate scope that represents this programming language. This should match one of the scopes listed in the grammars.yml file. Use "none" if there is no grammar for this language';
COMMENT ON COLUMN languages.group IS 'Name of the parent language. Languages in a group are counted in the statistics as the parent language';
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.languages;
DROP TYPE IF EXISTS language_type;
-- +goose StatementEnd