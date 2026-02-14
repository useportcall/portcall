package billing

import (
	"github.com/useportcall/portcall/libs/go/routerx"
	"github.com/useportcall/portcall/sdks/go/resources"
)

// CheckFeatureEntitlement checks entitlement for the current app and feature.
func CheckFeatureEntitlement(ctx *routerx.Context, featureID string) (*resources.Entitlement, error) {
	client, err := getPortcall(ctx)
	if err != nil {
		return nil, err
	}

	return client.Entitlements.Get(ctx, ctx.PublicAppID(), featureID)
}
