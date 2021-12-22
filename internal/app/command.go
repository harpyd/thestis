package app

import (
	"time"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	CreateTestCampaignCommand struct {
		OwnerID  string
		ViewName string
		Summary  string
	}

	LoadSpecificationCommand struct {
		TestCampaignID string
		LoadedByID     string
		Content        []byte
	}
)

type ParserOption func(b *specification.Builder)

func WithSpecificationID(specID string) ParserOption {
	return func(b *specification.Builder) {
		b.WithID(specID)
	}
}

func WithSpecificationOwnerID(ownerID string) ParserOption {
	return func(b *specification.Builder) {
		b.WithOwnerID(ownerID)
	}
}

func WithSpecificationLoadedAt(loadedAt time.Time) ParserOption {
	return func(b *specification.Builder) {
		b.WithLoadedAt(loadedAt)
	}
}
