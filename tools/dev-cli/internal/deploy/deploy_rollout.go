package deploy

import (
	"fmt"
	"strings"
	"time"
)

func clearStuckRelease() {
	script := `latest=$(kubectl get secrets -n portcall -o json 2>/dev/null | jq -r '.items[]|select(.type=="helm.sh/release.v1")|select(.metadata.name|startswith("sh.helm.release.v1.portcall.v"))|.metadata.name' | sed 's/sh.helm.release.v1.portcall.v//' | sort -rn | head -1)
[ -z "$latest" ] && exit 0
status=$(kubectl get secret "sh.helm.release.v1.portcall.v${latest}" -n portcall -o json 2>/dev/null | jq -r '.data.release' | base64 -d | base64 -d | gunzip 2>/dev/null | jq -r '.info.status // "unknown"')
[[ "$status" =~ ^pending|^failed ]] && kubectl delete secret "sh.helm.release.v1.portcall.v${latest}" -n portcall >/dev/null 2>&1 || true`
	_, _ = runShellOut(script)
}

func watchRollout(deploy string, timeout time.Duration) bool {
	start := time.Now()
	for time.Since(start) < timeout {
		ready, _ := runCmdOut("kubectl", "get", "deployment", deploy, "-n", "portcall", "-o", "jsonpath={.status.readyReplicas}")
		want, _ := runCmdOut("kubectl", "get", "deployment", deploy, "-n", "portcall", "-o", "jsonpath={.spec.replicas}")
		upd, _ := runCmdOut("kubectl", "get", "deployment", deploy, "-n", "portcall", "-o", "jsonpath={.status.updatedReplicas}")
		if strings.TrimSpace(want) == "" {
			want = "1"
		}
		fmt.Printf("\r  [%ds] Ready: %s/%s Updated: %s", int(time.Since(start).Seconds()), zero(ready), want, zero(upd))
		if zero(ready) == want && zero(upd) == want && zero(ready) != "0" {
			fmt.Print("\n")
			return true
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Print("\n")
	return false
}

func zero(v string) string {
	if strings.TrimSpace(v) == "" {
		return "0"
	}
	return strings.TrimSpace(v)
}
