package mock

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type Performer func(env *performance.Environment, t specification.Thesis) (fail, crash error)

func (p Performer) Perform(env *performance.Environment, t specification.Thesis) (fail, crash error) {
	return p(env, t)
}
