package mock

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type Performer func(c *performance.Context, t specification.Thesis) (fail error, err error)

func (p Performer) Perform(c *performance.Context, t specification.Thesis) (fail error, err error) {
	return p(c, t)
}
