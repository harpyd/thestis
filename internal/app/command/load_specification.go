package command

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/domain/user"
)

type LoadSpecificationHandler struct {
	specsRepo         app.SpecificationsRepository
	testCampaignsRepo app.TestCampaignsRepository
	specParserService app.SpecificationParser
}

func NewLoadSpecificationHandler(
	specsRepo app.SpecificationsRepository,
	testCampaignsRepo app.TestCampaignsRepository,
	specParserService app.SpecificationParser,
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

	tc, err := h.testCampaignsRepo.GetTestCampaign(ctx, cmd.TestCampaignID)
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

	if err := h.specsRepo.AddSpecification(ctx, spec); err != nil {
		return "", err
	}

	return specID, nil
}
