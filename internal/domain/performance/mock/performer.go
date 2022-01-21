package mock

import (
	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

type Performer func(env *performance.Environment, t specification.Thesis) performance.Result

func (p Performer) Perform(env *performance.Environment, t specification.Thesis) performance.Result {
	return p(env, t)
}
