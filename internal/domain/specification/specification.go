package specification

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type (
	Specification struct {
		id             string
		ownerID        string
		testCampaignID string
		loadedAt       time.Time

		author      string
		title       string
		description string
		stories     map[string]Story
	}

	Builder struct {
		id             string
		ownerID        string
		loadedAt       time.Time
		testCampaignID string
		author         string
		title          string
		description    string
		storyFactories []storyFactory
	}

	storyFactory func() (Story, error)
)

func (s *Specification) ID() string {
	return s.id
}

func (s *Specification) OwnerID() string {
	return s.ownerID
}

func (s *Specification) TestCampaignID() string {
	return s.testCampaignID
}

func (s *Specification) LoadedAt() time.Time {
	return s.loadedAt
}

func (s *Specification) Author() string {
	return s.author
}

func (s *Specification) Title() string {
	return s.title
}

func (s *Specification) Description() string {
	return s.description
}

func (s *Specification) Scenarios() []Scenario {
	scenarios := make([]Scenario, 0, s.ScenariosCount())

	for _, story := range s.stories {
		for _, scenario := range story.scenarios {
			scenarios = append(scenarios, scenario)
		}
	}

	return scenarios
}

func (s *Specification) ScenariosCount() int {
	count := 0

	for _, story := range s.stories {
		for range story.scenarios {
			count++
		}
	}

	return count
}

func (s *Specification) Story(slug string) (story Story, ok bool) {
	story, ok = s.stories[slug]

	return
}

func (s *Specification) Stories(slugs ...string) ([]Story, error) {
	if shouldGetAll(slugs) {
		return s.allStories(), nil
	}

	return s.filteredStories(slugs)
}

func (s *Specification) allStories() []Story {
	stories := make([]Story, 0, len(s.stories))

	for _, story := range s.stories {
		stories = append(stories, story)
	}

	return stories
}

func (s *Specification) filteredStories(slugs []string) ([]Story, error) {
	stories := make([]Story, 0, len(slugs))

	var err error

	for _, slug := range slugs {
		if story, ok := s.Story(slug); ok {
			stories = append(stories, story)
		} else {
			err = multierr.Append(err, NewNoSuchSlugError(NewStorySlug(slug)))
		}
	}

	return stories, err
}

func shouldGetAll(slugs []string) bool {
	return len(slugs) == 0
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build() (*Specification, error) {
	spec := &Specification{
		id:             b.id,
		ownerID:        b.ownerID,
		testCampaignID: b.testCampaignID,
		loadedAt:       b.loadedAt,
		author:         b.author,
		title:          b.title,
		description:    b.description,
		stories:        make(map[string]Story, len(b.storyFactories)),
	}

	if len(b.storyFactories) == 0 {
		return spec, NewBuildSpecificationError(NewNoSpecificationStoriesError())
	}

	var err error

	for _, stryFactory := range b.storyFactories {
		stry, stryErr := stryFactory()
		if _, ok := spec.stories[stry.Slug().Story()]; ok {
			err = multierr.Append(err, NewSlugAlreadyExistsError(stry.Slug()))

			continue
		}

		err = multierr.Append(err, stryErr)

		spec.stories[stry.Slug().Story()] = stry
	}

	return spec, NewBuildSpecificationError(err)
}

func (b *Builder) ErrlessBuild() *Specification {
	s, _ := b.Build()

	return s
}

func (b *Builder) Reset() {
	b.author = ""
	b.title = ""
	b.description = ""
	b.storyFactories = nil
}

func (b *Builder) WithID(id string) *Builder {
	b.id = id

	return b
}

func (b *Builder) WithOwnerID(ownerID string) *Builder {
	b.ownerID = ownerID

	return b
}

func (b *Builder) WithTestCampaignID(testCampaignID string) *Builder {
	b.testCampaignID = testCampaignID

	return b
}

func (b *Builder) WithLoadedAt(loadedAt time.Time) *Builder {
	b.loadedAt = loadedAt

	return b
}

func (b *Builder) WithAuthor(author string) *Builder {
	b.author = author

	return b
}

func (b *Builder) WithTitle(title string) *Builder {
	b.title = title

	return b
}

func (b *Builder) WithDescription(description string) *Builder {
	b.description = description

	return b
}

func (b *Builder) WithStory(slug string, buildFn func(b *StoryBuilder)) *Builder {
	sb := NewStoryBuilder()
	buildFn(sb)

	b.storyFactories = append(b.storyFactories, func() (Story, error) {
		return sb.Build(NewStorySlug(slug))
	})

	return b
}

type buildSpecificationError struct {
	err error
}

func NewBuildSpecificationError(err error) error {
	if err == nil {
		return nil
	}

	return errors.WithStack(buildSpecificationError{
		err: err,
	})
}

func IsBuildSpecificationError(err error) bool {
	var target buildSpecificationError

	return errors.As(err, &target)
}

func (e buildSpecificationError) Cause() error {
	return e.err
}

func (e buildSpecificationError) Unwrap() error {
	return e.err
}

func (e buildSpecificationError) NestedErrors() []error {
	return multierr.Errors(e.err)
}

func (e buildSpecificationError) CommonError() string {
	return "specification"
}

func (e buildSpecificationError) Error() string {
	return fmt.Sprintf("specification: %s", e.err)
}

var errNoSpecificationStories = errors.New("no stories")

func NewNoSpecificationStoriesError() error {
	return errNoSpecificationStories
}

func IsNoSpecificationStoriesError(err error) bool {
	return errors.Is(err, errNoSpecificationStories)
}
