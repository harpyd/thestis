package performance

import (
	"context"

	"github.com/harpyd/thestis/internal/domain/specification"
)

type (
	ScenarioSyncGroup struct {
		slug   specification.Slug
		theses map[specification.Slug]thesisSync
	}

	thesisSync struct {
		done chan struct{}
		deps []specification.Slug
	}
)

func SyncDependencies(scenario specification.Scenario) ScenarioSyncGroup {
	theses := scenario.Theses()

	syncs := make(map[specification.Slug]thesisSync, len(theses))

	for _, thesis := range theses {
		var (
			deps   = thesis.Dependencies()
			before = thesisBefore(scenario, thesis)
		)

		allDeps := make([]specification.Slug, 0, len(deps)+len(before))
		allDeps = append(allDeps, deps...)
		allDeps = append(allDeps, before...)

		syncs[thesis.Slug()] = thesisSync{
			done: make(chan struct{}),
			deps: allDeps,
		}
	}

	return ScenarioSyncGroup{
		slug:   scenario.Slug(),
		theses: syncs,
	}
}

func thesisBefore(scenario specification.Scenario, thesis specification.Thesis) []specification.Slug {
	beforeStages := thesis.Statement().Stage().Before()

	theses := scenario.ThesesByStages(beforeStages...)

	slugs := make([]specification.Slug, 0, len(theses))
	for _, before := range theses {
		slugs = append(slugs, before.Slug())
	}

	return slugs
}

type DependenciesSnapshot map[specification.Slug][]specification.Slug

func (s DependenciesSnapshot) Equal(other DependenciesSnapshot) bool {
	if len(s) != len(other) {
		return false
	}

	for sk, sv := range s {
		v, ok := other[sk]
		if !ok {
			return false
		}

		if !equalDependencies(slugsSet(sv), slugsSet(v)) {
			return false
		}
	}

	return true
}

func equalDependencies(left, right map[specification.Slug]bool) bool {
	if len(left) != len(right) {
		return false
	}

	for lk := range left {
		if !right[lk] {
			return false
		}
	}

	return true
}

func slugsSet(slugs []specification.Slug) map[specification.Slug]bool {
	set := make(map[specification.Slug]bool, len(slugs))
	for _, slug := range slugs {
		set[slug] = true
	}

	return set
}

func (g ScenarioSyncGroup) Snapshot() DependenciesSnapshot {
	if g.theses == nil || len(g.theses) == 0 {
		return nil
	}

	snp := make(DependenciesSnapshot, len(g.theses))

	for slug, sync := range g.theses {
		if len(sync.deps) == 0 {
			continue
		}

		snp[slug] = make([]specification.Slug, len(sync.deps))
		copy(snp[slug], sync.deps)
	}

	return snp
}

func (g ScenarioSyncGroup) Slug() specification.Slug {
	return g.slug
}

func (g ScenarioSyncGroup) WaitThesisDependencies(
	ctx context.Context,
	slug specification.Slug,
) error {
	for _, dep := range g.theses[slug].deps {
		thesis, ok := g.theses[dep]
		if !ok {
			continue
		}

		select {
		case <-thesis.done:
		case <-ctx.Done():
			return NewCanceledError(ctx.Err())
		}
	}

	return nil
}

func (g ScenarioSyncGroup) ThesisDone(slug specification.Slug) {
	if thesis, ok := g.theses[slug]; ok {
		close(thesis.done)
	}
}
