package specification

import (
	"strings"

	"github.com/pkg/errors"
)

type SlugKind string

const (
	NoSlug       SlugKind = ""
	StorySlug    SlugKind = "story"
	ScenarioSlug SlugKind = "scenario"
	ThesisSlug   SlugKind = "thesis"
)

type Slug struct {
	kind     SlugKind
	story    string
	scenario string
	thesis   string
}

var (
	ErrNotStorySlug    = errors.New("not story slug")
	ErrNotScenarioSlug = errors.New("not scenario slug")
	ErrNotThesisSlug   = errors.New("not thesis slug")
)

func NewStorySlug(slug string) Slug {
	return Slug{
		story: slug,
		kind:  StorySlug,
	}
}

func NewScenarioSlug(storySlug, scenarioSlug string) Slug {
	return Slug{
		story:    storySlug,
		scenario: scenarioSlug,
		kind:     ScenarioSlug,
	}
}

func NewThesisSlug(storySlug, scenarioSlug, thesisSlug string) Slug {
	return Slug{
		story:    storySlug,
		scenario: scenarioSlug,
		thesis:   thesisSlug,
		kind:     ThesisSlug,
	}
}

func (s Slug) ToStoryKind() Slug {
	if s.IsZero() {
		return s
	}

	return NewStorySlug(s.story)
}

func (s Slug) ToScenarioKind() Slug {
	if s.IsZero() {
		return s
	}

	return NewScenarioSlug(s.story, s.scenario)
}

func (s Slug) ToThesisKind() Slug {
	if s.IsZero() {
		return s
	}

	return NewThesisSlug(s.story, s.scenario, s.thesis)
}

func (s Slug) ShouldBeStoryKind() error {
	if s.kind == StorySlug {
		return nil
	}

	return ErrNotStorySlug
}

func (s Slug) ShouldBeScenarioKind() error {
	if s.kind == ScenarioSlug {
		return nil
	}

	return ErrNotScenarioSlug
}

func (s Slug) ShouldBeThesisKind() error {
	if s.kind == ThesisSlug {
		return nil
	}

	return ErrNotThesisSlug
}

func (s Slug) Story() string {
	return s.story
}

func (s Slug) Scenario() string {
	return s.scenario
}

func (s Slug) Thesis() string {
	return s.thesis
}

func (s Slug) Kind() SlugKind {
	return s.kind
}

// Partial returns a Kind part of the Slug.
// If Slug is zero, return an empty string.
func (s Slug) Partial() string {
	switch s.kind {
	case StorySlug:
		return s.story
	case ScenarioSlug:
		return s.scenario
	case ThesisSlug:
		return s.thesis
	case NoSlug:
		return ""
	}

	return ""
}

func (s Slug) IsZero() bool {
	return s == Slug{}
}

const (
	emptySlugReplace = "*"
	slugsSeparator   = "."
)

func (s Slug) String() string {
	var slugs []string

	switch s.kind {
	case StorySlug:
		slugs = []string{s.story}
	case ScenarioSlug:
		slugs = []string{s.story, s.scenario}
	case ThesisSlug:
		slugs = []string{s.story, s.scenario, s.thesis}
	case NoSlug:
	}

	return strings.Join(mapSlugs(slugs, replaceSlugIfEmpty), slugsSeparator)
}

func mapSlugs(slugs []string, fn func(string) string) []string {
	res := make([]string, 0, len(slugs))

	for _, s := range slugs {
		res = append(res, fn(s))
	}

	return res
}

func replaceSlugIfEmpty(s string) string {
	if s == "" {
		return emptySlugReplace
	}

	return s
}
