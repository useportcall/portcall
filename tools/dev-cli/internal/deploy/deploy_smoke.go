package deploy

import (
	"fmt"
	"strings"
)

func smoke(cluster string, apps []string) {
	api, dash := "http://localhost:8080", "http://localhost:8082"
	if cluster == "digitalocean" {
		api = "https://api.useportcall.com"
		dash = "https://dashboard.useportcall.com"
	}
	pass, failN := 0, 0
	for _, app := range apps {
		okTest := testApp(app, api, dash)
		if okTest {
			pass++
		} else {
			failN++
		}
	}
	plain("\nPassed: %d | Failed: %d", pass, failN)
	if failN > 0 {
		warn("Some smoke tests failed")
	}
}

func testApp(app, api, dash string) bool {
	tests := map[string][]string{
		"api":       {"curl", "-sf", "--max-time", "10", api + "/ping"},
		"dashboard": {"curl", "-sf", "--max-time", "10", dash + "/api/config"},
		"checkout":  {"curl", "-sf", "--max-time", "10", "https://checkout.useportcall.com"},
		"quote":     {"curl", "-sf", "--max-time", "10", "https://quote.useportcall.com/ping"},
		"file":      {"curl", "-sf", "--max-time", "10", "https://file.useportcall.com/ping"},
	}
	if app == "dashboard" {
		fmt.Printf("  Testing %s... ", app)
		if runCmd(tests["dashboard"][0], tests["dashboard"][1:]...) != nil {
			fail("FAILED")
			return false
		}
		if runCmd("curl", "-sf", "--max-time", "10", "https://webhook.useportcall.com/ping") != nil {
			fail("FAILED")
			return false
		}
		ok("OK")
		return true
	}
	if app == "admin" {
		fmt.Print("  Testing admin (internal pod health)... ")
		err := runCmd("kubectl", "exec", "-n", "portcall", "deploy/admin", "--", "curl", "-sf", "http://localhost:8081/health")
		if err == nil {
			ok("OK")
			return true
		}
		warn("SKIPPED (internal only)")
		return true
	}
	if app == "observability" {
		return testObservability()
	}
	if c, okT := tests[app]; okT {
		fmt.Printf("  Testing %s... ", app)
		if runCmd(c[0], c[1:]...) == nil {
			ok("OK")
			return true
		}
		fail("FAILED")
		return false
	}
	deploy := deployName(app)
	fmt.Printf("  Testing %s (pod running)... ", app)
	out, err := runCmdOut("kubectl", "get", "pods", "-n", "portcall", "-l", "app="+deploy, "--field-selector=status.phase=Running")
	if err == nil && strings.Contains(out, "Running") {
		ok("OK")
		return true
	}
	fail("FAILED")
	return false
}
