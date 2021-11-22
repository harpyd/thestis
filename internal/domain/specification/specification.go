package specification

type Specification struct {
	author      string
	title       string
	description string
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

type Builder struct {
	Specification
}

func NewBuilder() *Builder {
	return &Builder{}
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

func (b *Builder) Build() *Specification {
	return &b.Specification
}
