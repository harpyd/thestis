package yaml

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
		Method   string             `yaml:"method"`
		URL      string             `yaml:"url"`
		Request  httpRequestSchema  `yaml:"request"`
		Response httpResponseSchema `yaml:"response"`
	}

	httpRequestSchema struct {
		ContentType string                 `yaml:"contentType"`
		Body        map[string]interface{} `yaml:"body"`
	}

	httpResponseSchema struct {
		AllowedCodes       []int  `yaml:"allowedCodes"`
		AllowedContentType string `yaml:"allowedContentType"`
	}

	assertionSchema struct {
		Method string                 `yaml:"with"`
		Assert map[string]interface{} `yaml:"assert"`
	}
)
