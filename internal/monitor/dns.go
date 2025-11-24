package monitor

import (
	"context"
	"log"
	"net"
	"sort"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartDNSMonitoring(ctx context.Context, db *pgxpool.Pool) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		monitors, err := database.GetAllDNSMonitors(ctx, db)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch DNS monitors: %v", err)
		} else {
			for _, dnsMonitor := range monitors {
				go verifyDNS(ctx, db, dnsMonitor)
			}
		}
		select {
			case<-ctx.Done():
				return
			case<-ticker.C:
				continue
		}
	}
}

func verifyDNS(ctx context.Context, db *pgxpool.Pool, dnsMonitor *models.DNSMonitor) {
	a, aaaa, mx, ns := getDNSRecords(dnsMonitor.Domain)

	changed := isDifferent(dnsMonitor.LastA, a) || isDifferent(dnsMonitor.LastAAAA, aaaa) || isDifferent(dnsMonitor.LastMX, mx) || isDifferent(dnsMonitor.LastNS, ns)

	if changed {
		log.Printf("[INFO] DNS records changed for domain %s", dnsMonitor.Domain)

		err := database.UpdateDNSMonitorRecords(ctx, db, dnsMonitor.ID, a, aaaa, mx, ns)
		if err != nil {
			log.Printf("[ERROR] Failed to update DNS monitor records for domain %s: %v", dnsMonitor.Domain, err)
			return
		}

		userEmail, err := database.GetUserEmail(ctx, db, dnsMonitor.UserID)
		if err != nil {
			log.Printf("[ERROR] Failed to get user email for DNS monitor ID %d: %v", dnsMonitor.ID, err)
			return
		}
		go func(userEmail, domain, subject, message string) {
			notification.SendEmailNotification(userEmail, subject, domain, message)
		}(userEmail, dnsMonitor.Domain, dnsMonitor.Domain + "DNS records update", "The DNS records for "+dnsMonitor.Domain+" have changed.")

	}
}

func getDNSRecords(domain string) (a, aaaa, mx, ns []string) {
	ips, _ := net.LookupIP(domain)
	for _, ip := range ips {
		if ip.To4() != nil {
			a = append(a, ip.String())
		} else {
			aaaa = append(aaaa, ip.String())
		}
	}

	mxRecords, _ := net.LookupMX(domain)
	for _, mxRecord := range mxRecords {
		mx = append(mx, mxRecord.Host)
	}

	nsRecords, _ := net.LookupNS(domain)
	for _, nsRecord := range nsRecords {
		ns = append(ns, nsRecord.Host)
	}

	sort.Strings(a)
	sort.Strings(aaaa)
	sort.Strings(mx)
	sort.Strings(ns)

	return
}

func isDifferent(old, newRecords []string) bool{
	if len(old) != len(newRecords) {
		return true
	}
	for i := range old {
		if old[i] != newRecords[i] {
			return true
		}
	}
	return false
}
