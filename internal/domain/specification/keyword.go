package specification

type Keyword string

const (
	UnknownKeyword Keyword = ""
	Given          Keyword = "Given"
	When           Keyword = "When"
	Then           Keyword = "Then"
)
