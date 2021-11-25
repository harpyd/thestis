package specification

import (
	"go.uber.org/multierr"

	"github.com/harpyd/thestis/pkg/deepcopy"
)

type (
	Specification struct {
		id          string
		author      string
		title       string
		description string
		stories     map[string]Story
	}

	Story struct {
		slug        string
		description string
		asA         string
		inOrderTo   string
		wantTo      string
		scenarios   map[string]Scenario
	}

	Scenario struct {
		slug        string
		description string
		theses      map[string]Thesis
	}

	Thesis struct {
		slug      string
		statement Statement
		http      HTTP
		assertion Assertion
	}

	Statement struct {
		keyword  Keyword
		behavior string
	}

	HTTP struct {
		method   HTTPMethod
		url      string
		request  HTTPRequest
		response HTTPResponse
	}

	HTTPRequest struct {
		contentType ContentType
		body        map[string]interface{}
	}

	HTTPResponse struct {
		allowedCodes       []int
		allowedContentType ContentType
	}

	Assertion struct {
		typeOf  AssertionType
		method  AssertionMethod
		asserts []Assert
	}

	Assert struct {
		expected string
		actual   interface{}
	}
)

func (s *Specification) ID() string {
	return s.id
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
			err = multierr.Append(err, NewNoStoryError(slug))
		}
	}

	return stories, err
}

func (s *Specification) Story(slug string) (story Story, ok bool) {
	story, ok = s.stories[slug]

	return
}

func (s Story) Slug() string {
	return s.slug
}

func (s Story) Description() string {
	return s.description
}

func (s Story) AsA() string {
	return s.asA
}

func (s Story) InOrderTo() string {
	return s.inOrderTo
}

func (s Story) WantTo() string {
	return s.wantTo
}

func (s Story) Scenarios(slugs ...string) ([]Scenario, error) {
	if shouldGetAll(slugs) {
		return s.allScenarios(), nil
	}

	return s.filteredScenarios(slugs)
}

func (s Story) allScenarios() []Scenario {
	scenarios := make([]Scenario, 0, len(s.scenarios))

	for _, scenario := range s.scenarios {
		scenarios = append(scenarios, scenario)
	}

	return scenarios
}

func (s Story) filteredScenarios(slugs []string) ([]Scenario, error) {
	scenarios := make([]Scenario, 0, len(slugs))

	var err error

	for _, slug := range slugs {
		if scenario, ok := s.Scenario(slug); ok {
			scenarios = append(scenarios, scenario)
		} else {
			err = multierr.Append(err, NewNoScenarioError(slug))
		}
	}

	return scenarios, err
}

func (s Story) Scenario(slug string) (scenario Scenario, ok bool) {
	scenario, ok = s.scenarios[slug]

	return
}

func (s Scenario) Slug() string {
	return s.slug
}

func (s Scenario) Description() string {
	return s.description
}

func (s Scenario) Theses(slugs ...string) ([]Thesis, error) {
	if shouldGetAll(slugs) {
		return s.allTheses(), nil
	}

	return s.filteredTheses(slugs)
}

func (s Scenario) allTheses() []Thesis {
	theses := make([]Thesis, 0, len(s.theses))

	for _, thesis := range s.theses {
		theses = append(theses, thesis)
	}

	return theses
}

func (s Scenario) filteredTheses(slugs []string) ([]Thesis, error) {
	theses := make([]Thesis, 0, len(slugs))

	var err error

	for _, slug := range slugs {
		if thesis, ok := s.Thesis(slug); ok {
			theses = append(theses, thesis)
		} else {
			err = multierr.Append(err, NewNoThesisError(slug))
		}
	}

	return theses, err
}

func (s Scenario) Thesis(slug string) (thesis Thesis, ok bool) {
	thesis, ok = s.theses[slug]

	return
}

func (t Thesis) Slug() string {
	return t.slug
}

func (t Thesis) Statement() Statement {
	return t.statement
}

func (t Thesis) HTTP() HTTP {
	return t.http
}

func (t Thesis) Assertion() Assertion {
	return t.assertion
}

func newStatement(keyword Keyword, behavior string) (Statement, error) {
	if keyword == UnknownKeyword {
		return Statement{}, ErrUnknownKeyword
	}

	return Statement{
		keyword:  keyword,
		behavior: behavior,
	}, nil
}

func (s Statement) Keyword() Keyword {
	return s.keyword
}

func (s Statement) Behavior() string {
	return s.behavior
}

func (h HTTP) Method() HTTPMethod {
	return h.method
}

func (h HTTP) URL() string {
	return h.url
}

func (h HTTP) Request() HTTPRequest {
	return h.request
}

func (h HTTP) Response() HTTPResponse {
	return h.response
}

func (r HTTPRequest) ContentType() ContentType {
	return r.contentType
}

func (r HTTPRequest) Body() map[string]interface{} {
	return deepcopy.StringInterfaceMap(r.body)
}

func (r HTTPResponse) AllowedCodes() []int {
	return deepcopy.IntSlice(r.allowedCodes)
}

func (r HTTPResponse) AllowedContentType() ContentType {
	return r.allowedContentType
}

func (a Assertion) Type() AssertionType {
	return a.typeOf
}

func (a Assertion) Method() AssertionMethod {
	return a.method
}

func (a Assertion) Asserts() []Assert {
	asserts := make([]Assert, 0, len(a.asserts))
	copy(asserts, a.asserts)

	return asserts
}

func (a Assert) Expected() string {
	return a.expected
}

func (a Assert) Actual() interface{} {
	return deepcopy.Interface(a.actual)
}

func shouldGetAll(slugs []string) bool {
	return len(slugs) == 0
}
