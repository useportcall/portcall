package downscale

import (
	"fmt"
	"strconv"
	"strings"
)

func NodeMonthlyCost(run RunCmdOutWithEnv, env map[string]string, beforeSize string, beforeCount int, afterSize string, afterCount int) (float64, float64, error) {
	prices, err := DropletPrices(run, env)
	if err != nil {
		return 0, 0, err
	}
	before, okBefore := prices[beforeSize]
	after, okAfter := prices[afterSize]
	if !okBefore || !okAfter {
		return 0, 0, fmt.Errorf("missing droplet price for %s or %s", beforeSize, afterSize)
	}
	return before * float64(beforeCount), after * float64(afterCount), nil
}

func DropletPrices(run RunCmdOutWithEnv, env map[string]string) (map[string]float64, error) {
	out, err := run(env, "doctl", "compute", "size", "list", "--format", "Slug,PriceMonthly", "--no-header")
	if err != nil {
		return nil, err
	}
	prices := map[string]float64{}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		price, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			continue
		}
		prices[fields[0]] = price
	}
	return prices, nil
}
