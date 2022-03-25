package specification

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
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

func (s *Specification) Story(slug string) (story Story, ok bool) {
	story, ok = s.stories[slug]

	return
}

func (s *Specification) Stories() []Story {
	stories := make([]Story, 0, len(s.stories))

	for _, story := range s.stories {
		stories = append(stories, story)
	}

	return stories
}

func (s *Specification) StoriesCount() int {
	return len(s.stories)
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
		count += len(story.scenarios)
	}

	return count
}

func (s *Specification) Theses() []Thesis {
	theses := make([]Thesis, 0, s.ThesesCount())

	for _, story := range s.stories {
		for _, scenario := range story.scenarios {
			for _, thesis := range scenario.theses {
				theses = append(theses, thesis)
			}
		}
	}

	return theses
}

func (s *Specification) ThesesCount() int {
	count := 0

	for _, story := range s.stories {
		for _, scenario := range story.scenarios {
			count += len(scenario.theses)
		}
	}

	return count
}

var ErrNoSpecificationStories = errors.New("no stories")

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

	var w BuildErrorWrapper

	if len(b.storyFactories) == 0 {
		w.WithError(ErrNoSpecificationStories)
	}

	for _, storyFry := range b.storyFactories {
		story, err := storyFry()
		w.WithError(err)

		if _, ok := spec.stories[story.Slug().Story()]; ok {
			w.WithError(NewDuplicatedError(story.Slug()))

			continue
		}

		spec.stories[story.Slug().Story()] = story
	}

	return spec, w.Wrap("specification")
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
	var sb StoryBuilder

	buildFn(&sb)

	b.storyFactories = append(b.storyFactories, func() (Story, error) {
		return sb.Build(NewStorySlug(slug))
	})

	return b
}

type BuildErrorWrapper struct {
	errs []error
}

func (w *BuildErrorWrapper) WithError(err error) *BuildErrorWrapper {
	if err == nil {
		return w
	}

	w.errs = append(w.errs, err)

	return w
}

func (w *BuildErrorWrapper) Wrap(msg string) error {
	if len(w.errs) == 0 {
		return nil
	}

	return errors.WithStack(&BuildError{
		msg:  msg,
		errs: w.errs,
	})
}

func (w *BuildErrorWrapper) SluggedWrap(slug Slug) error {
	return w.Wrap(slug.String())
}

type BuildError struct {
	msg  string
	errs []error
}

func (e *BuildError) Message() string {
	return e.msg
}

func (e *BuildError) Errors() []error {
	errs := make([]error, len(e.errs))
	copy(errs, e.errs)

	return errs
}

const (
	errorMsgSeparator       = ": "
	leftNestedErrorsBorder  = "["
	rightNestedErrorsBorder = "]"
	nestedErrorsSeparator   = "; "
)

func (e *BuildError) Error() string {
	if e == nil {
		return ""
	}

	var b strings.Builder

	_, _ = b.WriteString(e.msg + errorMsgSeparator)
	_, _ = b.WriteString(leftNestedErrorsBorder)

	lastErrIdx := len(e.errs) - 1

	for _, err := range e.errs[:lastErrIdx] {
		_, _ = b.WriteString(fmt.Sprintf("%v", err))
		_, _ = b.WriteString(nestedErrorsSeparator)
	}

	_, _ = b.WriteString(fmt.Sprintf("%v", e.errs[lastErrIdx]))
	_, _ = b.WriteString(rightNestedErrorsBorder)

	return b.String()
}

func (e *BuildError) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (e *BuildError) As(target interface{}) bool {
	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}
