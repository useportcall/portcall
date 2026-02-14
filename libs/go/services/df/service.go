package df

import (
	"context"
	"log"
	"os"

	sdk "github.com/useportcall/portcall/sdks/go"
)

// Service is the dogfood service interface.
type Service interface {
	Increment(ctx context.Context, input *IncrementInput) error
	Decrement(ctx context.Context, input *DecrementInput) error
}

type service struct{}

// NewService creates a new dogfood service.
func NewService() Service {
	return &service{}
}

func (s *service) Increment(ctx context.Context, input *IncrementInput) error {
	secret := dfSecret(input.IsTest)
	client := sdk.New(sdk.Config{APIKey: secret})
	if err := client.MeterEvents.Increment(ctx, input.UserID, input.Feature); err != nil {
		log.Printf("Error incrementing feature %s for user %s: %v",
			input.Feature, input.UserID, err)
		return err
	}
	return nil
}

func (s *service) Decrement(ctx context.Context, input *DecrementInput) error {
	secret := dfSecret(input.IsTest)
	client := sdk.New(sdk.Config{APIKey: secret})
	if err := client.MeterEvents.Record(ctx, input.UserID, input.Feature, -1); err != nil {
		log.Printf("Error decrementing feature %s for user %s: %v",
			input.Feature, input.UserID, err)
		return err
	}
	return nil
}

func dfSecret(isTest bool) string {
	if isTest {
		return os.Getenv("DF_TEST_SECRET")
	}
	return os.Getenv("DF_SECRET")
}
