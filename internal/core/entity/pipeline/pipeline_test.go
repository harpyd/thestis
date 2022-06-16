package pipeline_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/harpyd/thestis/internal/core/entity/pipeline"
	"github.com/harpyd/thestis/internal/core/entity/specification"
)

func TestPipeline(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Pipeline                 *pipeline.Pipeline
		ExpectedID               string
		ExpectedSpecificationID  string
		ExpectedOwnerID          string
		ExpectedStarted          bool
		ExpectedWorkingScenarios []specification.Scenario
	}{
		{
			Pipeline:                 pipeline.Trigger("", nil),
			ExpectedID:               "",
			ExpectedSpecificationID:  "",
			ExpectedOwnerID:          "",
			ExpectedStarted:          false,
			ExpectedWorkingScenarios: nil,
		},
		{
			Pipeline: pipeline.Trigger(
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
			Pipeline: pipeline.Unmarshal(pipeline.Params{
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
			Pipeline: pipeline.Trigger(
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
			Pipeline: pipeline.Unmarshal(pipeline.Params{
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
				require.Equal(t, c.ExpectedID, c.Pipeline.ID())
			})

			t.Run("specification_id", func(t *testing.T) {
				require.Equal(t, c.ExpectedSpecificationID, c.Pipeline.SpecificationID())
			})

			t.Run("owner_id", func(t *testing.T) {
				require.Equal(t, c.ExpectedOwnerID, c.Pipeline.OwnerID())
			})

			t.Run("started", func(t *testing.T) {
				require.Equal(t, c.ExpectedStarted, c.Pipeline.Started())
			})

			if c.ExpectedStarted {
				t.Run("should_be_started", func(t *testing.T) {
					require.NoError(t, c.Pipeline.ShouldBeStarted())
				})
			} else {
				t.Run("not_started_error", func(t *testing.T) {
					require.ErrorIs(t, c.Pipeline.ShouldBeStarted(), pipeline.ErrNotStarted)
				})
			}

			t.Run("working_scenarios", func(t *testing.T) {
				require.ElementsMatch(t, c.ExpectedWorkingScenarios, c.Pipeline.WorkingScenarios())
			})
		})
	}
}

func TestStartPipeline(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		PipelineFactory func() *pipeline.Pipeline
		ExpectedSteps   []pipeline.Step
		ShouldBeErr     bool
		ExpectedErr     error
	}{
		{
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
					"",
					(&specification.Builder{}).ErrlessBuild(),
				)
			},
			ExpectedSteps: nil,
			ShouldBeErr:   false,
		},
		{
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Unmarshal(pipeline.Params{
					Started: true,
				})
			},
			ShouldBeErr: true,
			ExpectedErr: pipeline.ErrAlreadyStarted,
		},
		{
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
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
					pipeline.WithHTTP(pipeline.PassingExecutor()),
				)
			},
			ExpectedSteps: []pipeline.Step{
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.HTTPExecutor,
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.HTTPExecutor,
					pipeline.FiredPass,
				),
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredPass,
				),
			},
			ShouldBeErr: false,
		},
		{
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
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
					pipeline.WithAssertion(pipeline.FailingExecutor()),
				)
			},
			ExpectedSteps: []pipeline.Step{
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("que", "pue"),
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("que", "pue", "due"),
					pipeline.AssertionExecutor,
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStepWithErr(
					pipeline.WrapWithTerminatedError(
						errors.New("expected failing"),
						pipeline.FiredFail,
					),
					specification.NewThesisSlug("que", "pue", "due"),
					pipeline.AssertionExecutor,
					pipeline.FiredFail,
				),
				pipeline.NewScenarioStepWithErr(
					pipeline.WrapWithTerminatedError(
						errors.New("expected failing"),
						pipeline.FiredFail,
					),
					specification.NewScenarioSlug("que", "pue"),
					pipeline.FiredFail,
				),
			},
		},
		{
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
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
					pipeline.WithHTTP(pipeline.CancelingExecutor()),
				)
			},
			ExpectedSteps: []pipeline.Step{
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.HTTPExecutor,
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStepWithErr(
					pipeline.WrapWithTerminatedError(
						errors.New("expected canceling"),
						pipeline.FiredCancel,
					),
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.HTTPExecutor,
					pipeline.FiredCancel,
				),
				pipeline.NewScenarioStepWithErr(
					pipeline.WrapWithTerminatedError(
						errors.New("expected canceling"),
						pipeline.FiredCancel,
					),
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredCancel,
				),
			},
		},
		{
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
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
			ExpectedSteps: []pipeline.Step{
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.UnknownExecutor,
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStepWithErr(
					pipeline.WrapWithTerminatedError(
						pipeline.NewUndefinedExecutorError(
							pipeline.UnknownExecutor,
						),
						pipeline.FiredCrash,
					),
					specification.NewThesisSlug("foo", "bar", "baz"),
					pipeline.UnknownExecutor,
					pipeline.FiredCrash,
				),
				pipeline.NewScenarioStepWithErr(
					pipeline.WrapWithTerminatedError(
						pipeline.NewUndefinedExecutorError(
							pipeline.UnknownExecutor,
						),
						pipeline.FiredCrash,
					),
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredCrash,
				),
			},
			ShouldBeErr: false,
		},
		{
			PipelineFactory: func() *pipeline.Pipeline {
				return pipeline.Trigger(
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
					pipeline.WithHTTP(pipeline.PassingExecutor()),
					pipeline.WithAssertion(pipeline.CrashingExecutor()),
				)
			},
			ExpectedSteps: []pipeline.Step{
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("rod", "sod"),
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("rod", "sod", "nod"),
					pipeline.HTTPExecutor,
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("rod", "sod", "nod"),
					pipeline.HTTPExecutor,
					pipeline.FiredPass,
				),
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("rod", "sod"),
					pipeline.FiredPass,
				),
				pipeline.NewScenarioStep(
					specification.NewScenarioSlug("rod", "dod"),
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("rod", "dod", "mod"),
					pipeline.HTTPExecutor,
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("rod", "dod", "mod"),
					pipeline.HTTPExecutor,
					pipeline.FiredPass,
				),
				pipeline.NewThesisStep(
					specification.NewThesisSlug("rod", "dod", "zod"),
					pipeline.AssertionExecutor,
					pipeline.FiredExecute,
				),
				pipeline.NewThesisStepWithErr(
					pipeline.WrapWithTerminatedError(
						errors.New("expected crashing"),
						pipeline.FiredCrash,
					),
					specification.NewThesisSlug("rod", "dod", "zod"),
					pipeline.AssertionExecutor,
					pipeline.FiredCrash,
				),
				pipeline.NewScenarioStepWithErr(
					pipeline.WrapWithTerminatedError(
						errors.New("expected crashing"),
						pipeline.FiredCrash,
					),
					specification.NewScenarioSlug("rod", "dod"),
					pipeline.FiredCrash,
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

			steps, err := c.PipelineFactory().Start(ctx)

			must := func() {
				_ = c.PipelineFactory().MustStart(ctx)
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

func TestOneExecutingAtATime(t *testing.T) {
	t.Parallel()

	pipe := pipeline.Trigger("foo", validSpecification(t))

	_, err := pipe.Start(context.Background())
	require.NoError(t, err, "Start with error")

	_, err = pipe.Start(context.Background())

	require.ErrorIs(
		t,
		err,
		pipeline.ErrAlreadyStarted,
		"Err is not already started error",
	)
}

func TestPipelineStartByStart(t *testing.T) {
	t.Parallel()

	pipe := pipeline.Trigger("bar", validSpecification(t))

	finish := make(chan bool)

	go func() {
		select {
		case <-finish:
		case <-time.After(10 * time.Millisecond):
			require.Fail(t, "Timeout exceeded, test is not finished")
		}
	}()

	s1, err := pipe.Start(context.Background())

	require.NoError(t, err, "First start with error")

	for range s1 {
		// read flow steps of first call
	}

	s2, err := pipe.Start(context.Background())

	require.NoError(t, err, "Second start with error")

	for range s2 {
		// read flow steps of second call
	}

	finish <- true
}

func TestCancelPipelineContext(t *testing.T) {
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
		ExpectedIncludedSteps    []pipeline.Step
	}{
		{
			CancelBeforeStart: true,
			ExpectedIncludedSteps: []pipeline.Step{
				pipeline.NewScenarioStepWithErr(
					pipeline.WrapWithTerminatedError(
						context.Canceled,
						pipeline.FiredCancel,
					),
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredCancel,
				),
			},
		},
		{
			CancelAfterReadFirstStep: true,
			ExpectedIncludedSteps: []pipeline.Step{
				pipeline.NewScenarioStepWithErr(
					pipeline.WrapWithTerminatedError(
						context.Canceled,
						pipeline.FiredCancel,
					),
					specification.NewScenarioSlug("foo", "bar"),
					pipeline.FiredCancel,
				),
			},
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			pipe := pipeline.Trigger(
				"foo",
				spec,
				pipeline.WithHTTP(pipeline.PassingExecutor()),
				pipeline.WithAssertion(pipeline.FailingExecutor()),
			)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if c.CancelBeforeStart {
				cancel()
			}

			steps, err := pipe.Start(ctx)
			require.NoError(t, err)

			if c.CancelAfterReadFirstStep {
				<-steps

				cancel()
			}

			requireStepsContain(t, steps, c.ExpectedIncludedSteps...)

			_, err = pipe.Start(context.Background())

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
			GivenError: pipeline.WrapWithTerminatedError(
				nil,
				pipeline.FiredCrash,
			),
			ExpectedIs: false,
		},
		{
			GivenError: pipeline.WrapWithTerminatedError(
				errors.New("foo"),
				pipeline.FiredFail,
			),
			ExpectedIs: false,
		},
		{
			GivenError: pipeline.WrapWithTerminatedError(
				errTest,
				pipeline.FiredCrash,
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
			GivenError: pipeline.WrapWithTerminatedError(
				nil,
				pipeline.FiredFail,
			),
			ExpectedAs: false,
		},
		{
			GivenError: pipeline.WrapWithTerminatedError(
				errors.New("foo"),
				pipeline.FiredFail,
			),
			ExpectedAs: false,
		},
		{
			GivenError: pipeline.WrapWithTerminatedError(
				testError{},
				pipeline.FiredCrash,
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
		ExpectedEvent     pipeline.Event
		ExpectedUnwrapped error
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:        &pipeline.TerminatedError{},
			ShouldBeWrapped:   true,
			ExpectedEvent:     pipeline.NoEvent,
			ExpectedUnwrapped: nil,
		},
		{
			GivenError: pipeline.WrapWithTerminatedError(
				errors.New("foo"),
				pipeline.FiredCrash,
			),
			ShouldBeWrapped:   true,
			ExpectedEvent:     pipeline.FiredCrash,
			ExpectedUnwrapped: errors.New("foo"),
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var target *pipeline.TerminatedError

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
			GivenError:          &pipeline.TerminatedError{},
			ExpectedErrorString: "pipeline has terminated",
		},
		{
			GivenError: pipeline.WrapWithTerminatedError(
				errors.New("foo"),
				pipeline.NoEvent,
			),
			ExpectedErrorString: "pipeline has terminated: foo",
		},
		{
			GivenError: pipeline.WrapWithTerminatedError(
				errors.New("bar"),
				pipeline.FiredCrash,
			),
			ExpectedErrorString: `pipeline has terminated due to "crash" event: bar`,
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

func TestAsUndefinedExecutorError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError           error
		ShouldBeWrapped      bool
		ExpectedExecutorType pipeline.ExecutorType
	}{
		{
			GivenError:      nil,
			ShouldBeWrapped: false,
		},
		{
			GivenError:           &pipeline.UndefinedExecutorError{},
			ShouldBeWrapped:      true,
			ExpectedExecutorType: pipeline.NoExecutor,
		},
		{
			GivenError:           pipeline.NewUndefinedExecutorError(pipeline.HTTPExecutor),
			ShouldBeWrapped:      true,
			ExpectedExecutorType: pipeline.HTTPExecutor,
		},
	}

	for i := range testCases {
		c := testCases[i]

		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			var rerr *pipeline.UndefinedExecutorError

			if !c.ShouldBeWrapped {
				t.Run("not", func(t *testing.T) {
					require.False(t, errors.As(c.GivenError, &rerr))
				})

				return
			}

			t.Run("as", func(t *testing.T) {
				require.ErrorAs(t, c.GivenError, &rerr)

				t.Run("executor_type", func(t *testing.T) {
					require.Equal(t, c.ExpectedExecutorType, rerr.ExecutorType())
				})
			})
		})
	}
}

func TestFormatUndefinedExecutorError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		GivenError          error
		ExpectedErrorString string
	}{
		{
			GivenError: pipeline.NewUndefinedExecutorError(
				pipeline.UnknownExecutor,
			),
			ExpectedErrorString: "undefined executor with `!` type",
		},
		{
			GivenError:          pipeline.NewUndefinedExecutorError("foo"),
			ExpectedErrorString: "undefined executor with `foo` type",
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

func requireStepsMatch(t *testing.T, expected []pipeline.Step, actual <-chan pipeline.Step) {
	t.Helper()

	require.ElementsMatch(
		t,
		mapStepsSliceToStrings(expected),
		mapStepsChanToStrings(actual, len(expected)),
	)
}

func requireStepsContain(
	t *testing.T,
	steps <-chan pipeline.Step,
	contain ...pipeline.Step,
) {
	t.Helper()

	require.Subset(
		t,
		mapStepsChanToStrings(steps, len(contain)),
		mapStepsSliceToStrings(contain),
	)
}

func mapStepsSliceToStrings(steps []pipeline.Step) []string {
	result := make([]string, 0, len(steps))

	for _, step := range steps {
		result = append(result, step.String())
	}

	return result
}

func mapStepsChanToStrings(steps <-chan pipeline.Step, capacity int) []string {
	result := make([]string, 0, capacity)

	for step := range steps {
		result = append(result, step.String())
	}

	return result
}

type testError struct{}

func (e testError) Error() string {
	return "test"
}
