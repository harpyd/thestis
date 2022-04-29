package yaml

import (
	specification2 "github.com/harpyd/thestis/internal/core/domain/specification"
)

type (
	specificationSchema struct {
		Author      string                 `yaml:"author"`
		Title       string                 `yaml:"title"`
		Description string                 `yaml:"description"`
		Stories     map[string]storySchema `yaml:"stories"`
	}

	storySchema struct {
		Description string                    `yaml:"description"`
		AsA         string                    `yaml:"asA"`
		InOrderTo   string                    `yaml:"inOrderTo"`
		WantTo      string                    `yaml:"wantTo"`
		Scenarios   map[string]scenarioSchema `yaml:"scenarios"`
	}

	scenarioSchema struct {
		Description string                  `yaml:"description"`
		Theses      map[string]thesisSchema `yaml:"theses"`
	}

	thesisSchema struct {
		Given     string          `yaml:"given"`
		When      string          `yaml:"when"`
		Then      string          `yaml:"then"`
		After     []string        `yaml:"after"`
		HTTP      httpSchema      `yaml:"http"`
		Assertion assertionSchema `yaml:"assertion"`
	}

	httpSchema struct {
		Request  httpRequestSchema  `yaml:"request"`
		Response httpResponseSchema `yaml:"response"`
	}

	httpRequestSchema struct {
		Method      specification2.HTTPMethod  `yaml:"method"`
		URL         string                     `yaml:"url"`
		ContentType specification2.ContentType `yaml:"contentType"`
		Body        map[string]interface{}     `yaml:"body"`
	}

	httpResponseSchema struct {
		AllowedCodes       []int                      `yaml:"allowedCodes"`
		AllowedContentType specification2.ContentType `yaml:"allowedContentType"`
	}

	assertionSchema struct {
		Method specification2.AssertionMethod `yaml:"with"`
		Assert []assertSchema                 `yaml:"assert"`
	}

	assertSchema struct {
		Actual   string      `yaml:"actual"`
		Expected interface{} `yaml:"expected"`
	}
)
