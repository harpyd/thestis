package specification

import "strings"

type SlugKind string

const (
	NoSlug       SlugKind = ""
	StorySlug    SlugKind = "story"
	ScenarioSlug SlugKind = "scenario"
	ThesisSlug   SlugKind = "thesis"
)

type Slug struct {
	story    string
	scenario string
	thesis   string

	kind SlugKind
}

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

func (s Slug) IsZero() bool {
	return s == Slug{}
}

const (
	emptyReplace   = "*"
	slugsSeparator = "."
)

func (s Slug) String() string {
	switch s.kind {
	case StorySlug:
		return replaceIfEmpty(s.story)
	case ScenarioSlug:
		slugs := mapSlugs([]string{
			s.story,
			s.scenario,
		}, replaceIfEmpty)

		return strings.Join(slugs, slugsSeparator)
	case ThesisSlug:
		slugs := mapSlugs([]string{
			s.story,
			s.scenario,
			s.thesis,
		}, replaceIfEmpty)

		return strings.Join(slugs, slugsSeparator)
	case NoSlug:
	}

	return ""
}

func mapSlugs(slugs []string, fn func(string) string) []string {
	res := make([]string, 0, len(slugs))

	for _, s := range slugs {
		res = append(res, fn(s))
	}

	return res
}

func replaceIfEmpty(s string) string {
	if s == "" {
		return emptyReplace
	}

	return s
}
