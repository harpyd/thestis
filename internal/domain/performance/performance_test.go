package performance_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestPerformanceCreation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Performance              *performance.Performance
		ExpectedID               string
		ExpectedSpecificationID  string
		ExpectedOwnerID          string
		ExpectedStarted          bool
		ExpectedWorkingScenarios []specification.Scenario
	}{
		{
			Performance:              performance.FromSpecification("", nil),
			ExpectedID:               "",
			ExpectedSpecificationID:  "",
			ExpectedOwnerID:          "",
			ExpectedStarted:          false,
			ExpectedWorkingScenarios: nil,
		},
		{
			Performance: performance.FromSpecification(
				"",
				specification.NewBuilder().
					ErrlessBuild(),
			),
			ExpectedID:               "",
			ExpectedSpecificationID:  "",
			ExpectedOwnerID:          "",
			ExpectedStarted:          false,
			ExpectedWorkingScenarios: nil,
		},
		{
			Performance: performance.Unmarshal(performance.Params{
				Specification: specification.NewBuilder().
					ErrlessBuild(),
			}),
			ExpectedID:               "",
			ExpectedSpecificationID:  "",
			ExpectedOwnerID:          "",
			ExpectedStarted:          false,
			ExpectedWorkingScenarios: nil,
		},
		{
			Performance: performance.FromSpecification(
				"foo",
				specification.NewBuilder().
					WithID("bar").
					WithOwnerID("baz").
					WithStory("moo", func(b *specification.StoryBuilder) {
						b.WithScenario("koo", func(b *specification.ScenarioBuilder) {
							b.WithThesis("too", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild(),
			),
			ExpectedID:              "foo",
			ExpectedSpecificationID: "bar",
			ExpectedOwnerID:         "baz",
			ExpectedStarted:         false,
			ExpectedWorkingScenarios: []specification.Scenario{
				specification.NewScenarioBuilder().
					WithThesis("too", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("moo", "koo")),
			},
		},
		{
			Performance: performance.Unmarshal(performance.Params{
				ID: "foo",
				Specification: specification.NewBuilder().
					WithID("spc").
					WithStory("boo", func(b *specification.StoryBuilder) {
						b.WithScenario("zoo", func(b *specification.ScenarioBuilder) {
							b.WithThesis("doo", func(b *specification.ThesisBuilder) {})
						})
						b.WithScenario("koo", func(b *specification.ScenarioBuilder) {
							b.WithThesis("poo", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild(),
				OwnerID: "djr",
				Started: true,
			}),
			ExpectedID:              "foo",
			ExpectedSpecificationID: "spc",
			ExpectedOwnerID:         "djr",
			ExpectedStarted:         true,
			ExpectedWorkingScenarios: []specification.Scenario{
				specification.NewScenarioBuilder().
					WithThesis("doo", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("boo", "zoo")),
				specification.NewScenarioBuilder().
					WithThesis("poo", func(b *specification.ThesisBuilder) {}).
					ErrlessBuild(specification.NewScenarioSlug("boo", "koo")),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			t.Run("id", func(t *testing.T) {
				require.Equal(t, c.ExpectedID, c.Performance.ID())
			})

			t.Run("specification_id", func(t *testing.T) {
				require.Equal(t, c.ExpectedSpecificationID, c.Performance.SpecificationID())
			})

			t.Run("owner_id", func(t *testing.T) {
				require.Equal(t, c.ExpectedOwnerID, c.Performance.OwnerID())
			})

			t.Run("started", func(t *testing.T) {
				require.Equal(t, c.ExpectedStarted, c.Performance.Started())
			})

			if c.ExpectedStarted {
				t.Run("should_be_started", func(t *testing.T) {
					require.NoError(t, c.Performance.ShouldBeStarted())
				})
			} else {
				t.Run("not_started_error", func(t *testing.T) {
					require.True(t, performance.IsNotStartedError(c.Performance.ShouldBeStarted()))
				})
			}

			t.Run("working_scenarios", func(t *testing.T) {
				require.ElementsMatch(t, c.ExpectedWorkingScenarios, c.Performance.WorkingScenarios())
			})
		})
	}
}

func TestStartPerformance(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenPerformance *performance.Performance
		ExpectedSteps    []performance.Step
		ShouldBeErr      bool
		IsErr            func(err error) bool
	}{
		{
			GivenPerformance: performance.FromSpecification(
				"",
				specification.NewBuilder().ErrlessBuild(),
			),
			ExpectedSteps: nil,
			ShouldBeErr:   false,
		},
		{
			GivenPerformance: performance.Unmarshal(performance.Params{
				Started: true,
			}),
			ShouldBeErr: true,
			IsErr:       performance.IsAlreadyStartedError,
		},
		{
			GivenPerformance: performance.FromSpecification(
				"ddq",
				specification.NewBuilder().
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
							b.WithThesis("baz", func(b *specification.ThesisBuilder) {
								b.WithHTTP(func(b *specification.HTTPBuilder) {
									b.WithRequest(func(b *specification.HTTPRequestBuilder) {
										b.WithMethod(specification.GET)
										b.WithURL("https://some-url.com")
									})
								})
							})
						})
					}).
					ErrlessBuild(),
				performance.WithHTTP(performance.PassingPerformer()),
			),
			ExpectedSteps: []performance.Step{
				performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.HTTPPerformer,
					performance.FiredPass,
				),
				performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPass,
				),
			},
			ShouldBeErr: false,
		},
		{
			GivenPerformance: performance.FromSpecification(
				"bvs",
				specification.NewBuilder().
					WithStory("que", func(b *specification.StoryBuilder) {
						b.WithScenario("pue", func(b *specification.ScenarioBuilder) {
							b.WithThesis("due", func(b *specification.ThesisBuilder) {
								b.WithAssertion(func(b *specification.AssertionBuilder) {
									b.WithMethod(specification.JSONPath)
								})
							})
						})
					}).
					ErrlessBuild(),
				performance.WithAssertion(performance.FailingPerformer()),
			),
			ExpectedSteps: []performance.Step{
				performance.NewScenarioStep(
					specification.NewScenarioSlug("que", "pue"),
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("que", "pue", "due"),
					performance.AssertionPerformer,
					performance.FiredPerform,
				),
				performance.NewThesisStepWithErr(
					performance.NewFailedError(
						errors.New("expected failing"),
					),
					specification.NewThesisSlug("que", "pue", "due"),
					performance.AssertionPerformer,
					performance.FiredFail,
				),
				performance.NewScenarioStepWithErr(
					performance.NewFailedError(
						errors.New("expected failing"),
					),
					specification.NewScenarioSlug("que", "pue"),
					performance.FiredFail,
				),
			},
		},
		{
			GivenPerformance: performance.FromSpecification(
				"daq",
				specification.NewBuilder().
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
							b.WithThesis("baz", func(b *specification.ThesisBuilder) {
								b.WithHTTP(func(b *specification.HTTPBuilder) {
									b.WithRequest(func(b *specification.HTTPRequestBuilder) {
										b.WithMethod(specification.GET)
										b.WithURL("https://test.com")
									})
								})
							})
						})
					}).
					ErrlessBuild(),
				performance.WithHTTP(performance.CancelingPerformer()),
			),
			ExpectedSteps: []performance.Step{
				performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				),
				performance.NewThesisStepWithErr(
					performance.NewCanceledError(
						errors.New("expected canceling"),
					),
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.HTTPPerformer,
					performance.FiredCancel,
				),
				performance.NewScenarioStepWithErr(
					performance.NewCanceledError(
						errors.New("expected canceling"),
					),
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredCancel,
				),
			},
		},
		{
			GivenPerformance: performance.FromSpecification(
				"hpd",
				specification.NewBuilder().
					WithStory("foo", func(b *specification.StoryBuilder) {
						b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
							b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
						})
					}).
					ErrlessBuild(),
			),
			ExpectedSteps: []performance.Step{
				performance.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.UnknownPerformer,
					performance.FiredPerform,
				),
				performance.NewThesisStepWithErr(
					performance.NewCrashedError(
						performance.NewNoSatisfyingPerformerError(
							performance.UnknownPerformer,
						),
					),
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.UnknownPerformer,
					performance.FiredCrash,
				),
				performance.NewScenarioStepWithErr(
					performance.NewCrashedError(
						performance.NewNoSatisfyingPerformerError(
							performance.UnknownPerformer,
						),
					),
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredCrash,
				),
			},
			ShouldBeErr: false,
		},
		{
			GivenPerformance: performance.FromSpecification(
				"jqd",
				specification.NewBuilder().
					WithStory("rod", func(b *specification.StoryBuilder) {
						b.WithScenario("dod", func(b *specification.ScenarioBuilder) {
							b.WithThesis("mod", func(b *specification.ThesisBuilder) {
								b.WithStatement(specification.Given, "mod")
								b.WithHTTP(func(b *specification.HTTPBuilder) {
									b.WithRequest(func(b *specification.HTTPRequestBuilder) {
										b.WithMethod(specification.GET)
										b.WithURL("https://test-api.net")
									})
								})
							})
							b.WithThesis("zod", func(b *specification.ThesisBuilder) {
								b.WithStatement(specification.Then, "zod")
								b.WithAssertion(func(b *specification.AssertionBuilder) {
									b.WithMethod(specification.JSONPath)
								})
							})
						})
						b.WithScenario("sod", func(b *specification.ScenarioBuilder) {
							b.WithThesis("nod", func(b *specification.ThesisBuilder) {
								b.WithStatement(specification.Given, "nod")
								b.WithHTTP(func(b *specification.HTTPBuilder) {
									b.WithRequest(func(b *specification.HTTPRequestBuilder) {
										b.WithURL("https://last-api.com")
									})
								})
							})
						})
					}).
					ErrlessBuild(),
				performance.WithHTTP(performance.PassingPerformer()),
				performance.WithAssertion(performance.CrashingPerformer()),
			),
			ExpectedSteps: []performance.Step{
				performance.NewScenarioStep(
					specification.NewScenarioSlug("rod", "sod"),
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("rod", "sod", "nod"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("rod", "sod", "nod"),
					performance.HTTPPerformer,
					performance.FiredPass,
				),
				performance.NewScenarioStep(
					specification.NewScenarioSlug("rod", "sod"),
					performance.FiredPass,
				),
				performance.NewScenarioStep(
					specification.NewScenarioSlug("rod", "dod"),
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("rod", "dod", "mod"),
					performance.HTTPPerformer,
					performance.FiredPerform,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("rod", "dod", "mod"),
					performance.HTTPPerformer,
					performance.FiredPass,
				),
				performance.NewThesisStep(
					specification.NewThesisSlug("rod", "dod", "zod"),
					performance.AssertionPerformer,
					performance.FiredPerform,
				),
				performance.NewThesisStepWithErr(
					performance.NewCrashedError(
						errors.New("expected crashing"),
					),
					specification.NewThesisSlug("rod", "dod", "zod"),
					performance.AssertionPerformer,
					performance.FiredCrash,
				),
				performance.NewScenarioStepWithErr(
					performance.NewCrashedError(
						errors.New("expected crashing"),
					),
					specification.NewScenarioSlug("rod", "dod"),
					performance.FiredCrash,
				),
			},
			ShouldBeErr: false,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			steps, err := c.GivenPerformance.Start(context.Background())

			if c.ShouldBeErr {
				t.Run("err", func(t *testing.T) {
					require.True(t, c.IsErr(err))
				})
			} else {
				t.Run("steps", func(t *testing.T) {
					require.NoError(t, err)

					requireStepsMatch(t, c.ExpectedSteps, steps)
				})
			}
		})
	}
}

func TestOnePerformingAtATime(t *testing.T) {
	t.Parallel()

	perf := performance.FromSpecification("foo", validSpecification(t))

	_, err := perf.Start(context.Background())
	require.NoError(t, err, "Start with error")

	_, err = perf.Start(context.Background())

	require.True(t, performance.IsAlreadyStartedError(err), "Err is not already started error")
}

func TestPerformanceStartByStart(t *testing.T) {
	t.Parallel()

	perf := performance.FromSpecification("bar", validSpecification(t))

	finish := make(chan bool)

	go func() {
		select {
		case <-finish:
		case <-time.After(10 * time.Millisecond):
			require.Fail(t, "Timeout exceeded, test is not finished")
		}
	}()

	s1, err := perf.Start(context.Background())

	require.NoError(t, err, "First start with error")

	for range s1 {
		// read flow steps of first call
	}

	s2, err := perf.Start(context.Background())

	require.NoError(t, err, "Second start with error")

	for range s2 {
		// read flow steps of second call
	}

	finish <- true
}

func TestCancelPerformanceContext(t *testing.T) {
	t.Parallel()

	spec := specification.NewBuilder().
		WithStory("foo", func(b *specification.StoryBuilder) {
			b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
				b.WithThesis("saz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Given, "saz")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod(specification.POST)
							b.WithURL("https://preparing.net")
						})
					})
				})
				b.WithThesis("faz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.When, "faz")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod(specification.GET)
							b.WithURL("https://testing.net")
						})
					})
				})
				b.WithThesis("daz", func(b *specification.ThesisBuilder) {
					b.WithStatement(specification.Then, "daz")
					b.WithAssertion(func(b *specification.AssertionBuilder) {
						b.WithMethod(specification.JSONPath)
					})
				})
			})
		}).
		ErrlessBuild()

	testCases := []struct {
		CancelBeforeStart        bool
		CancelAfterReadFirstStep bool
		ExpectedIncludedSteps    []performance.Step
	}{
		{
			CancelBeforeStart: true,
			ExpectedIncludedSteps: []performance.Step{
				performance.NewScenarioStepWithErr(
					performance.NewCanceledError(context.Canceled),
					specification.AnyScenarioSlug(),
					performance.FiredCancel,
				),
			},
		},
		{
			CancelAfterReadFirstStep: true,
			ExpectedIncludedSteps: []performance.Step{
				performance.NewScenarioStepWithErr(
					performance.NewCanceledError(context.Canceled),
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredCancel,
				),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			perf := performance.FromSpecification(
				"foo",
				spec,
				performance.WithHTTP(performance.PassingPerformer()),
				performance.WithAssertion(performance.FailingPerformer()),
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if c.CancelBeforeStart {
				cancel()
			}

			steps, err := perf.Start(ctx)
			require.NoError(t, err)

			if c.CancelAfterReadFirstStep {
				<-steps

				cancel()
			}

			requireStepsContain(t, steps, c.ExpectedIncludedSteps...)

			_, err = perf.Start(context.Background())

			require.NoError(t, err, "Second start with error")
		})
	}
}

func TestPerformanceErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name     string
		Err      error
		IsErr    func(err error) bool
		Reversed bool
	}{
		{
			Name:  "performance_already_started_error",
			Err:   performance.NewAlreadyStartedError(),
			IsErr: performance.IsAlreadyStartedError,
		},
		{
			Name:     "NON_performance_already_started_error",
			Err:      errors.New("performance already started"),
			IsErr:    performance.IsAlreadyStartedError,
			Reversed: true,
		},
		{
			Name:  "performance_not_started_error",
			Err:   performance.NewNotStartedError(),
			IsErr: performance.IsNotStartedError,
		},
		{
			Name:     "NON_performance_not_started_error",
			Err:      errors.New("performance not started"),
			IsErr:    performance.IsNotStartedError,
			Reversed: true,
		},
		{
			Name:  "no_satisfying_performer_error",
			Err:   performance.NewNoSatisfyingPerformerError(performance.HTTPPerformer),
			IsErr: performance.IsNoSatisfyingPerformerError,
		},
		{
			Name:     "NON_no_satisfying_performer_error",
			Err:      errors.New("no satisfying performer for HTTP"),
			IsErr:    performance.IsNoSatisfyingPerformerError,
			Reversed: true,
		},
	}

	for _, c := range testCases {
		c := c

		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			if c.Reversed {
				require.False(t, c.IsErr(c.Err))

				return
			}

			require.True(t, c.IsErr(c.Err))
		})
	}
}

func validSpecification(t *testing.T) *specification.Specification {
	t.Helper()

	spec, err := specification.NewBuilder().
		WithStory("story", func(b *specification.StoryBuilder) {
			b.WithScenario("scenario", func(b *specification.ScenarioBuilder) {
				b.WithThesis("a", func(b *specification.ThesisBuilder) {
					b.WithStatement("given", "state")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod("POST")
							b.WithURL("https://some-api/endpoint")
						})
						b.WithResponse(func(b *specification.HTTPResponseBuilder) {
							b.WithAllowedCodes([]int{201})
							b.WithAllowedContentType("application/json")
						})
					})
				})
				b.WithThesis("b", func(b *specification.ThesisBuilder) {
					b.WithStatement("when", "action")
					b.WithHTTP(func(b *specification.HTTPBuilder) {
						b.WithRequest(func(b *specification.HTTPRequestBuilder) {
							b.WithMethod("GET")
							b.WithURL("https://some-api/endpoint")
						})
						b.WithResponse(func(b *specification.HTTPResponseBuilder) {
							b.WithAllowedCodes([]int{200})
						})
					})
				})
				b.WithThesis("c", func(b *specification.ThesisBuilder) {
					b.WithStatement("then", "check")
					b.WithAssertion(func(b *specification.AssertionBuilder) {
						b.WithMethod("jsonpath")
						b.WithAssert("some", "some")
					})
				})
			})
		}).
		Build()

	require.NoError(t, err)

	return spec
}

func requireStepsMatch(t *testing.T, expected []performance.Step, actual <-chan performance.Step) {
	t.Helper()

	require.ElementsMatch(
		t,
		mapStepsSliceToStrings(expected),
		mapStepsChanToStrings(actual, len(expected)),
	)
}

func requireStepsContain(
	t *testing.T,
	steps <-chan performance.Step,
	contain ...performance.Step,
) {
	t.Helper()

	require.Subset(
		t,
		mapStepsChanToStrings(steps, len(contain)),
		mapStepsSliceToStrings(contain),
	)
}

func mapStepsSliceToStrings(steps []performance.Step) []string {
	strs := make([]string, 0, len(steps))

	for _, step := range steps {
		strs = append(strs, step.String())
	}

	return strs
}

func mapStepsChanToStrings(steps <-chan performance.Step, capacity int) []string {
	strs := make([]string, 0, capacity)

	for step := range steps {
		strs = append(strs, step.String())
	}

	return strs
}
