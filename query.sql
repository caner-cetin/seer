-- name: GetLanguages :many
SELECT *
FROM languages;

-- name: GetLanguageCount :one
SELECT COUNT(id)
FROM languages;