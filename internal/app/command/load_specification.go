package command

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/specification"
	"github.com/harpyd/thestis/internal/domain/testcampaign"
)

type LoadSpecificationHandler struct {
	specRepo          specificationRepository
	testCampaignsRepo testCampaignsRepository
	specParserService specificationParserService
}

func NewLoadSpecificationHandler(
	specRepo specificationRepository,
	testCampaignsRepo testCampaignsRepository,
	specParserService specificationParserService,
) LoadSpecificationHandler {
	if specRepo == nil {
		panic("spec repository is nil")
	}

	if testCampaignsRepo == nil {
		panic("test campaigns repository is nil")
	}

	if specParserService == nil {
		panic("spec parser service")
	}

	return LoadSpecificationHandler{
		specRepo:          specRepo,
		testCampaignsRepo: testCampaignsRepo,
		specParserService: specParserService,
	}
}

func (h LoadSpecificationHandler) Handle(
	ctx context.Context,
	cmd app.LoadSpecificationCommand,
) (specID string, err error) {
	defer func() {
		err = errors.Wrap(err, "spec loading")
	}()

	specID = uuid.New().String()

	spec, err := h.specParserService.ParseSpecification(specID, strings.NewReader(cmd.Content))
	if err != nil {
		return "", err
	}

	if err = h.testCampaignsRepo.UpdateTestCampaign(
		ctx,
		cmd.TestCampaignID,
		h.loadSpecification(spec),
	); err != nil {
		return "", err
	}

	return
}

func (h LoadSpecificationHandler) loadSpecification(spec *specification.Specification) TestCampaignUpdater {
	return func(ctx context.Context, tc *testcampaign.TestCampaign) (*testcampaign.TestCampaign, error) {
		if err := h.specRepo.AddSpecification(ctx, spec); err != nil {
			return nil, err
		}

		tc.SetActiveSpecificationID(spec.ID())

		return tc, nil
	}
}
