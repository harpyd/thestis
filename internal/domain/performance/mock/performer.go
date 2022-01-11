package mock

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type Performer func(c *performance.Context, t specification.Thesis)

func (p Performer) Perform(c *performance.Context, t specification.Thesis) {
	p(c, t)
}
