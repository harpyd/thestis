package specification

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
