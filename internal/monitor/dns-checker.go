package monitor

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func checkDNS(m models.Monitor) models.CheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result string
	var err error
	start := time.Now()

	r := net.DefaultResolver

	switch m.Type {
	case models.TypeDNS_A:
		result, err = lookupIP(ctx, r, m.Target, "ip4")
	case models.TypeDNS_AAAA:
		result, err = lookupIP(ctx, r, m.Target, "ip6")
	case models.TypeDNS_MX:
		result, err = lookupMX(ctx, r, m.Target)
	case models.TypeDNS_NS:
		result, err = lookupNS(ctx, r, m.Target)
	default:
		return models.CheckResult{MonitorID: m.ID, Status: models.StatusDown, Message: "Invalid DNS Type", CheckedAt: time.Now()}
	}

	latency := time.Since(start)

	if err != nil {
		return models.CheckResult{
			MonitorID: m.ID,
			Status:    models.StatusDown,
			Message:   fmt.Sprintf("Resolution failed: %v", err),
			Latency:   latency,
			CheckedAt: time.Now(),
		}
	}

	status := models.StatusUp
	if m.ExpectedValue != "" && !strings.Contains(result, m.ExpectedValue) {
		status = models.StatusDown
		result = fmt.Sprint("DNS Record has changed\nExpected: %s\nObtained: %s", m.ExpectedValue, result)
	}

	return models.CheckResult{
		MonitorID: m.ID,
		Status:    status,
		Message:   result,
		Latency:   latency,
		CheckedAt: time.Now(),
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

	return strings.Join(results, ", "), nil
}
