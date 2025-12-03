package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func checkDNS(m models.Monitor) models.CheckResult {
	var config models.DNSConfig
	if err := json.Unmarshal(m.Config, &config); err != nil {
		return models.CheckResult{Status: models.StatusDown, Message: "[ERROR] DNS configuration error.", CheckedAt: time.Now()}
	}

	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return d.DialContext(ctx, network, "8.8.8.8:53")
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var resultString string
	var err error

	switch config.RecordType {
	case "A":
		resultString, err = lookupIP(ctx, r, m.Target, "ip4")
	case "AAAA":
		resultString, err = lookupIP(ctx, r, m.Target, "ip6")
	case "MX":
		resultString, err = lookupMX(ctx, r, m.Target)
	case "NS":
		resultString, err = lookupNS(ctx, r, m.Target)
	default:
		return models.CheckResult{Status: models.StatusDown, Message: "Invalid DNS record type."}
	}

	if err != nil {
		return models.CheckResult{
			MonitorID: m.ID,
			Status:    models.StatusDown,
			Message:   fmt.Sprintf("Could not resolve the specified %s record for target %s", config.RecordType, err.Error()),
			CheckedAt: time.Now(),
		}
	}

	currentValue := strings.TrimSpace(resultString)
	expectedParts := strings.Split(config.ExpectedValue, ",")

	for i := range expectedParts {
		expectedParts[i] = strings.TrimSpace(expectedParts[i])
	}

	sort.Strings(expectedParts)

	expectedValue := strings.Join(expectedParts, ", ")

	if expectedValue == "" {
		return models.CheckResult{
			MonitorID:   m.ID,
			Status:      models.StatusUp,
			ResultValue: currentValue,
			Message:     "Valores DNS detectados automaticamente.",
			CheckedAt:   time.Now(),
		}
	}

	if currentValue != expectedValue {
		return models.CheckResult{
			MonitorID:   m.ID,
			Status:      models.StatusDown,
			Message:     fmt.Sprintf("DNS record mismatch. Expected '%s', Found '%s'", expectedValue, currentValue),
			ResultValue: currentValue,
			CheckedAt:   time.Now(),
		}
	}

	return models.CheckResult{
		MonitorID:   m.ID,
		Status:      models.StatusUp,
		ResultValue: resultString,
		CheckedAt:   time.Now(),
	}
}

func lookupIP(ctx context.Context, r *net.Resolver, host string, network string) (string, error) {
	ips, err := r.LookupIP(ctx, network, host)
	if err != nil {
		return "", err
	}

	var results []string
	for _, ip := range ips {
		results = append(results, ip.String())
	}

	sort.Strings(results)

	return strings.Join(results, ", "), nil
}

func lookupMX(ctx context.Context, r *net.Resolver, host string) (string, error) {
	mxs, err := r.LookupMX(ctx, host)
	if err != nil {
		return "", err
	}

	var results []string
	for _, mx := range mxs {
		results = append(results, fmt.Sprintf("%s (%d)", mx.Host, mx.Pref))
	}

	sort.Strings(results)

	return strings.Join(results, ", "), nil
}

func lookupNS(ctx context.Context, r *net.Resolver, host string) (string, error) {
	nss, err := r.LookupNS(ctx, host)
	if err != nil {
		return "", nil
	}

	var results []string
	for _, ns := range nss {
		results = append(results, ns.Host)
	}

	sort.Strings(results)

	return strings.Join(results, ", "), nil
}
