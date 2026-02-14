package app

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/discordx"
)

func sendAccountSignupNotification(email, projectName string, apps []*models.App) {
	testID, liveID := "", ""
	for _, app := range apps {
		if app.IsLive {
			liveID = app.PublicID
			continue
		}
		testID = app.PublicID
	}

	discordx.SendFromEnvAsync(
		"DISCORD_WEBHOOK_URL_SIGNUP",
		fmt.Sprintf(
			"New account signed up on Dashboard: %s created project %q (test app: %s, live app: %s)",
			email,
			projectName,
			testID,
			liveID,
		),
	)
}
