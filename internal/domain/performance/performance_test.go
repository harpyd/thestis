package performance_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/domain/performance"
	"github.com/harpyd/thestis/internal/domain/specification"
)

func TestPerformance(t *testing.T) {
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
				(&specification.Builder{}).
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
				Specification: (&specification.Builder{}).
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
				(&specification.Builder{}).
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
				(&specification.ScenarioBuilder{}).
					WithThesis("too", func(b *specification.ThesisBuilder) {}).
					Build(specification.NewScenarioSlug("moo", "koo")),
			},
		},
		{
			Performance: performance.Unmarshal(performance.Params{
				ID: "foo",
				Specification: (&specification.Builder{}).
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
				(&specification.ScenarioBuilder{}).
					WithThesis("doo", func(b *specification.ThesisBuilder) {}).
					Build(specification.NewScenarioSlug("boo", "zoo")),
				(&specification.ScenarioBuilder{}).
					WithThesis("poo", func(b *specification.ThesisBuilder) {}).
					Build(specification.NewScenarioSlug("boo", "koo")),
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
					require.ErrorIs(t, c.Performance.ShouldBeStarted(), performance.ErrNotStarted)
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
		PerformanceFactory func() *performance.Performance
		ExpectedSteps      []performance.Step
		ShouldBeErr        bool
		ExpectedErr        error
	}{
		{
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
					"",
					(&specification.Builder{}).ErrlessBuild(),
				)
			},
			ExpectedSteps: nil,
			ShouldBeErr:   false,
		},
		{
			PerformanceFactory: func() *performance.Performance {
				return performance.Unmarshal(performance.Params{
					Started: true,
				})
			},
			ShouldBeErr: true,
			ExpectedErr: performance.ErrAlreadyStarted,
		},
		{
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
					"ddq",
					(&specification.Builder{}).
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
				)
			},
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
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
					"bvs",
					(&specification.Builder{}).
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
				)
			},
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
					performance.WrapWithTerminatedError(
						errors.New("expected failing"),
						performance.FiredFail,
					),
					specification.NewThesisSlug("que", "pue", "due"),
					performance.AssertionPerformer,
					performance.FiredFail,
				),
				performance.NewScenarioStepWithErr(
					performance.WrapWithTerminatedError(
						errors.New("expected failing"),
						performance.FiredFail,
					),
					specification.NewScenarioSlug("que", "pue"),
					performance.FiredFail,
				),
			},
		},
		{
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
					"daq",
					(&specification.Builder{}).
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
				)
			},
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
					performance.WrapWithTerminatedError(
						errors.New("expected canceling"),
						performance.FiredCancel,
					),
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.HTTPPerformer,
					performance.FiredCancel,
				),
				performance.NewScenarioStepWithErr(
					performance.WrapWithTerminatedError(
						errors.New("expected canceling"),
						performance.FiredCancel,
					),
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredCancel,
				),
			},
		},
		{
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
					"hpd",
					(&specification.Builder{}).
						WithStory("foo", func(b *specification.StoryBuilder) {
							b.WithScenario("bar", func(b *specification.ScenarioBuilder) {
								b.WithThesis("baz", func(b *specification.ThesisBuilder) {})
							})
						}).
						ErrlessBuild(),
				)
			},
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
					performance.WrapWithTerminatedError(
						performance.NewRejectedError(
							performance.UnknownPerformer,
						),
						performance.FiredCrash,
					),
					specification.NewThesisSlug("foo", "bar", "baz"),
					performance.UnknownPerformer,
					performance.FiredCrash,
				),
				performance.NewScenarioStepWithErr(
					performance.WrapWithTerminatedError(
						performance.NewRejectedError(
							performance.UnknownPerformer,
						),
						performance.FiredCrash,
					),
					specification.NewScenarioSlug("foo", "bar"),
					performance.FiredCrash,
				),
			},
			ShouldBeErr: false,
		},
		{
			PerformanceFactory: func() *performance.Performance {
				return performance.FromSpecification(
					"jqd",
					(&specification.Builder{}).
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
				)
			},
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
					performance.WrapWithTerminatedError(
						errors.New("expected crashing"),
						performance.FiredCrash,
					),
					specification.NewThesisSlug("rod", "dod", "zod"),
					performance.AssertionPerformer,
					performance.FiredCrash,
				),
				performance.NewScenarioStepWithErr(
					performance.WrapWithTerminatedError(
						errors.New("expected crashing"),
						performance.FiredCrash,
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

			ctx := context.Background()

			steps, err := c.PerformanceFactory().Start(ctx)

			must := func() {
				_ = c.PerformanceFactory().MustStart(ctx)
			}

			if c.ShouldBeErr {
				t.Run("err", func(t *testing.T) {
					require.ErrorIs(t, err, c.ExpectedErr)
				})

				t.Run("panics", func(t *testing.T) {
					require.PanicsWithError(t, c.ExpectedErr.Error(), must)
				})
			} else {
				t.Run("not_panics", func(t *testing.T) {
					require.NotPanics(t, must)
				})

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

	require.ErrorIs(
		t,
		err,
		performance.ErrAlreadyStarted,
		"Err is not already started error",
	)
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

	spec := (&specification.Builder{}).
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
					performance.WrapWithTerminatedError(
						context.Canceled,
						performance.FiredCancel,
					),
					specification.AnyScenarioSlug(),
					performance.FiredCancel,
				),
			},
		},
		{
			CancelAfterReadFirstStep: true,
			ExpectedIncludedSteps: []performance.Step{
				performance.NewScenarioStepWithErr(
					performance.WrapWithTerminatedError(
						context.Canceled,
						performance.FiredCancel,
					),
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

var errTest = errors.New("test")

func TestIsWrappedInTerminatedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError error
		ExpectedIs bool
	}{
		{
			GivenError: nil,
			ExpectedIs: false,
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				nil,
				performance.FiredCrash,
			),
			ExpectedIs: false,
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				errors.New("foo"),
				performance.FiredFail,
			),
			ExpectedIs: false,
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				errTest,
				performance.FiredCrash,
			),
			ExpectedIs: true,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, c.ExpectedIs, errors.Is(c.GivenError, errTest))
		})
	}
}

func TestAsWrappedInTerminatedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError error
		ExpectedAs bool
	}{
		{
			GivenError: nil,
			ExpectedAs: false,
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				nil,
				performance.FiredFail,
			),
			ExpectedAs: false,
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				errors.New("foo"),
				performance.FiredFail,
			),
			ExpectedAs: false,
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				testError{},
				performance.FiredCrash,
			),
			ExpectedAs: true,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target testError

			require.Equal(t, c.ExpectedAs, errors.As(c.GivenError, &target))
		})
	}
}

func TestAsTerminatedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError        error
		ShouldBeWrapped   bool
		ExpectedEvent     performance.Event
		ExpectedUnwrapped error
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &performance.TerminatedError{},
			ShouldBeWrapped:   true,
			ExpectedEvent:     performance.NoEvent,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				errors.New("foo"),
				performance.FiredCrash,
			),
			ShouldBeWrapped:   true,
			ExpectedEvent:     performance.FiredCrash,
			ExpectedUnwrapped: errors.New("foo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *performance.TerminatedError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &target))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &target)

				t.Run("unwrap", func(t *testing.T) {
					if c.ExpectedUnwrapped != nil {
						require.EqualError(t, target.Unwrap(), c.ExpectedUnwrapped.Error())

						return
					}

					require.NoError(t, target.Unwrap())
				})

				t.Run("event", func(t *testing.T) {
					require.Equal(t, c.ExpectedEvent, target.Event())
				})
			})
		})
	}
}

func TestFormatTerminatedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError:          &performance.TerminatedError{},
			ExpectedErrorString: "performance has terminated",
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				errors.New("foo"),
				performance.NoEvent,
			),
			ExpectedErrorString: "performance has terminated: foo",
		},
		{
			GivenError: performance.WrapWithTerminatedError(
				errors.New("bar"),
				performance.FiredCrash,
			),
			ExpectedErrorString: "performance has terminated due to `crash` event: bar",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}

func TestAsRejectedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError            error
		ShouldBeWrapped       bool
		ExpectedPerformerType performance.PerformerType
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:            &performance.RejectedError{},
			ShouldBeWrapped:       true,
			ExpectedPerformerType: performance.NoPerformer,
		},
		{
			GivenError:            performance.NewRejectedError(performance.HTTPPerformer),
			ShouldBeWrapped:       true,
			ExpectedPerformerType: performance.HTTPPerformer,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var rerr *performance.RejectedError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &rerr))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &rerr)

				t.Run("performer_type", func(t *testing.T) {
					require.Equal(t, c.ExpectedPerformerType, rerr.PerformerType())
				})
			})
		})
	}
}

func TestFormatRejectedError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError: performance.NewRejectedError(
				performance.UnknownPerformer,
			),
			ExpectedErrorString: "rejected performer with `!` type",
		},
		{
			GivenError:          performance.NewRejectedError("foo"),
			ExpectedErrorString: "rejected performer with `foo` type",
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			require.EqualError(t, c.GivenError, c.ExpectedErrorString)
		})
	}
}

func validSpecification(t *testing.T) *specification.Specification {
	t.Helper()

	spec, err := (&specification.Builder{}).
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

type testError struct{}

func (e testError) Error() string {
	return "test"
}
