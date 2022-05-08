package command

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/entity/user"
)

type LoadSpecification struct {
	TestCampaignID string
	LoadedByID     string
	Content        []byte
}

type LoadSpecificationHandler interface {
	Handle(ctx context.Context, cmd LoadSpecification) (string, error)
}

type loadSpecificationHandler struct {
	specRepo          service.SpecificationRepository
	testCampaignRepo  service.TestCampaignRepository
	specParserService service.SpecificationParser
}

func NewLoadSpecificationHandler(
	specRepo service.SpecificationRepository,
	testCampaignRepo service.TestCampaignRepository,
	specParser service.SpecificationParser,
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

	return loadSpecificationHandler{
		specRepo:          specRepo,
		testCampaignRepo:  testCampaignRepo,
		specParserService: specParser,
	}
}

func (h loadSpecificationHandler) Handle(
	ctx context.Context,
	cmd LoadSpecification,
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
		service.WithSpecificationID(specID),
		service.WithSpecificationTestCampaignID(tc.ID()),
		service.WithSpecificationOwnerID(tc.OwnerID()),
		service.WithSpecificationLoadedAt(time.Now().UTC()),
	)
	if err != nil {
		return "", err
	}

	if err := h.specRepo.AddSpecification(ctx, spec); err != nil {
		return "", err
	}

	return specID, nil
}
