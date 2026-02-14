package main

import "fmt"

func testObservability() bool {
	labels := []string{"app=loki", "app=grafana", "app=promtail"}
	fmt.Print("  Testing observability (pods running)... ")
	for _, l := range labels {
		out, err := runCmdOut("kubectl", "get", "pods", "-n", "portcall", "-l", l, "--field-selector=status.phase=Running")
		if err != nil || out == "" {
			fail("FAILED")
			return false
		}
	}
	ok("OK")
	return true
}
