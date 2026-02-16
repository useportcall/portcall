package deploy

import "fmt"

type deployApp struct {
	Name       string
	Image      string
	Dockerfile string
	Context    string
	ChartOnly  bool
}

var deployApps = []deployApp{
	{"api", "registry.digitalocean.com/portcall-registry/portcall-api", "apps/api/Dockerfile", ".", false},
	{"admin", "registry.digitalocean.com/portcall-registry/portcall-admin", "apps/admin/Dockerfile", ".", false},
	{"dashboard", "registry.digitalocean.com/portcall-registry/portcall-dashboard", "apps/dashboard/Dockerfile", ".", false},
	{"checkout", "registry.digitalocean.com/portcall-registry/portcall-checkout", "apps/checkout/Dockerfile", ".", false},
	{"file", "registry.digitalocean.com/portcall-registry/portcall-file", "apps/file/Dockerfile", ".", false},
	{"quote", "registry.digitalocean.com/portcall-registry/portcall-quote", "apps/quote/Dockerfile", ".", false},
	{"billing", "registry.digitalocean.com/portcall-registry/portcall-billing-worker", "apps/billing/Dockerfile", ".", false},
	{"email", "registry.digitalocean.com/portcall-registry/portcall-email-worker", "apps/email/Dockerfile", ".", false},
	{"keycloak", "registry.digitalocean.com/portcall-registry/portcall-keycloak", "docker-compose/keycloak/Dockerfile", "docker-compose/keycloak", false},
	{"observability", "", "", "", true},
}

func listDeployApps() {
	info("Available apps:")
	for i, app := range deployApps {
		fmt.Printf("  %d) %s\n", i+1, app.Name)
	}
}

func findDeployApp(name string) (deployApp, bool) {
	for _, app := range deployApps {
		if app.Name == name {
			return app, true
		}
	}
	return deployApp{}, false
}
