package hmeili

import "github.com/meilisearch/meilisearch-go"

func BuildSearchSettings(searchableAttributes []string) *meilisearch.Settings {
	return &meilisearch.Settings{
		SearchableAttributes: searchableAttributes,
		TypoTolerance: &meilisearch.TypoTolerance{
			Enabled: true,
		},
		SeparatorTokens: []string{
			" ", "-", "_", "/", ".",
		},
	}
}
