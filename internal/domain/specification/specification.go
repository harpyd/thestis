package specification

import (
	"fmt"
	"io"
	"strings"
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

	for _, storyFry := range b.storyFactories {
		story, storyErr := storyFry()
		if _, ok := spec.stories[story.Slug().Story()]; ok {
			err = multierr.Append(err, NewSlugAlreadyExistsError(story.Slug()))

			continue
		}

		err = multierr.Append(err, storyErr)

		spec.stories[story.Slug().Story()] = story
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

type Error struct {
	msg  string
	errs []error
}

func WrapErrors(msg string, errs ...error) error {
	nonNilErrs := make([]error, 0, len(errs))

	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	}

	return errors.WithStack(&Error{
		msg:  msg,
		errs: nonNilErrs,
	})
}

func WrapErrorsFromSlug(slug Slug, errs ...error) error {
	return WrapErrors(slug.String(), errs...)
}

func (e *Error) Message() string {
	return e.msg
}

func (e *Error) Errors() []error {
	errs := make([]error, len(e.errs))
	copy(errs, e.errs)

	return errs
}

func (e *Error) Format(f fmt.State, verb rune) {
	if verb == 'v' && f.Flag('+') {
		e.writeMultiLine(f)

		return
	}

	e.writeSingleLine(f)
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	var b strings.Builder

	e.writeSingleLine(&b)

	return b.String()
}

const (
	singleLineMsgErrorSeparator = ": "
	singleLineErrorSeparator    = "; "

	multiLineMsgErrorSeparator = ":"
	multiLineErrorSeparator    = "\n    "
	multiLineErrorIndent       = "    "
)

func (e *Error) writeSingleLine(w io.Writer) {
	_, _ = io.WriteString(w, e.msg+singleLineMsgErrorSeparator)

	lastErrIdx := len(e.errs) - 1

	for _, err := range e.errs[:lastErrIdx] {
		_, _ = io.WriteString(w, fmt.Sprintf("%v", err))
		_, _ = io.WriteString(w, singleLineErrorSeparator)
	}

	_, _ = io.WriteString(w, fmt.Sprintf("%v", e.errs[lastErrIdx]))
}

func (e *Error) writeMultiLine(w io.Writer) {
	_, _ = io.WriteString(w, e.msg+multiLineMsgErrorSeparator)

	for _, err := range e.errs {
		_, _ = io.WriteString(w, multiLineErrorSeparator)
		writePrefixLine(w, multiLineErrorIndent, fmt.Sprintf("%v", err))
	}
}

func writePrefixLine(w io.Writer, prefix, s string) {
	first := true

	for len(s) > 0 {
		if first {
			first = false
		} else {
			_, _ = io.WriteString(w, prefix)
		}

		idx := strings.IndexByte(s, '\n')
		if idx < 0 {
			idx = len(s) - 1
		}

		_, _ = io.WriteString(w, s[:idx+1])
		s = s[idx+1:]
	}
}

func (e *Error) Is(target error) bool {
	for _, err := range e.errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (e *Error) As(target interface{}) bool {
	for _, err := range e.errs {
		if errors.As(err, target) {
			return true
		}
	}

	return false
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
