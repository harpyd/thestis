package command

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
	"github.com/harpyd/thestis/internal/domain/user"
)

type LoadSpecificationHandler struct {
	specsRepo         specificationsRepository
	testCampaignsRepo testCampaignsRepository
	specParserService specificationParserService
}

func NewLoadSpecificationHandler(
	specsRepo specificationsRepository,
	testCampaignsRepo testCampaignsRepository,
	specParserService specificationParserService,
) LoadSpecificationHandler {
	if specsRepo == nil {
		panic("specifications repository is nil")
	}

	if testCampaignsRepo == nil {
		panic("test campaigns repository is nil")
	}

	if specParserService == nil {
		panic("specification parser service is nil")
	}

	return LoadSpecificationHandler{
		specsRepo:         specsRepo,
		testCampaignsRepo: testCampaignsRepo,
		specParserService: specParserService,
	}
}

func (h LoadSpecificationHandler) Handle(
	ctx context.Context,
	cmd app.LoadSpecificationCommand,
) (specID string, err error) {
	defer func() {
		err = errors.Wrap(err, "specification loading")
	}()

	specID = uuid.New().String()

	if err = h.testCampaignsRepo.UpdateTestCampaign(
		ctx,
		cmd.TestCampaignID,
		h.loadSpecification(specID, cmd),
	); err != nil {
		return "", err
	}

	return
}

func (h LoadSpecificationHandler) loadSpecification(specID string, cmd app.LoadSpecificationCommand) TestCampaignUpdater {
	return func(ctx context.Context, tc *testcampaign.TestCampaign) (*testcampaign.TestCampaign, error) {
		spec, err := h.specParserService.ParseSpecification(
			bytes.NewReader(cmd.Content),
			app.WithSpecificationID(specID),
			app.WithSpecificationOwnerID(tc.OwnerID()),
			app.WithSpecificationLoadedAt(time.Now()),
		)
		if err != nil {
			return nil, err
		}

		if err := user.CanSeeTestCampaign(cmd.LoadedByID, tc); err != nil {
			return nil, err
		}

		if err := h.specsRepo.AddSpecification(ctx, spec); err != nil {
			return nil, err
		}

		tc.SetActiveSpecificationID(spec.ID())

		return tc, nil
	}
}
