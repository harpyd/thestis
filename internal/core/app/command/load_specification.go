package command

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type LoadSpecificationHandler struct {
	specRepo          app.SpecificationRepository
	testCampaignRepo  app.TestCampaignRepository
	specParserService app.SpecificationParser
}

func NewLoadSpecificationHandler(
	specRepo app.SpecificationRepository,
	testCampaignRepo app.TestCampaignRepository,
	specParser app.SpecificationParser,
) LoadSpecificationHandler {
	if specRepo == nil {
		panic("specification repository is nil")
	}

	if testCampaignRepo == nil {
		panic("test campaign repository is nil")
	}

	if specParser == nil {
		panic("specification parser is nil")
	}

	return LoadSpecificationHandler{
		specRepo:          specRepo,
		testCampaignRepo:  testCampaignRepo,
		specParserService: specParser,
	}
}

func (h LoadSpecificationHandler) Handle(
	ctx context.Context,
	cmd app.LoadSpecificationCommand,
) (specID string, err error) {
	defer func() {
		err = errors.Wrap(err, "specification loading")
	}()

	tc, err := h.testCampaignRepo.GetTestCampaign(ctx, cmd.TestCampaignID)
	if err != nil {
		return "", err
	}

	if err := user.CanAccessTestCampaign(cmd.LoadedByID, tc, user.Read); err != nil {
		return "", err
	}

	specID = uuid.New().String()

	spec, err := h.specParserService.ParseSpecification(
		bytes.NewReader(cmd.Content),
		app.WithSpecificationID(specID),
		app.WithSpecificationTestCampaignID(tc.ID()),
		app.WithSpecificationOwnerID(tc.OwnerID()),
		app.WithSpecificationLoadedAt(time.Now().UTC()),
	)
	if err != nil {
		return "", err
	}

	if err := h.specRepo.AddSpecification(ctx, spec); err != nil {
		return "", err
	}

	return specID, nil
}
