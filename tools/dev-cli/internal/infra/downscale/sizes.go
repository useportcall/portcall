package downscale

import (
	"regexp"
	"strconv"
)

var sizeSlugPattern = regexp.MustCompile(`(\d+)vcpu-(\d+)gb`)

func IsDowngradeSize(current string, target string) bool {
	curScore, curOK := SizeScore(current)
	tgtScore, tgtOK := SizeScore(target)
	if !curOK || !tgtOK {
		return false
	}
	return tgtScore < curScore
}

func SizeScore(slug string) (int, bool) {
	parts := sizeSlugPattern.FindStringSubmatch(slug)
	if len(parts) != 3 {
		return 0, false
	}
	cpu, errCPU := strconv.Atoi(parts[1])
	mem, errMem := strconv.Atoi(parts[2])
	if errCPU != nil || errMem != nil {
		return 0, false
	}
	return cpu*1000 + mem, true
}
