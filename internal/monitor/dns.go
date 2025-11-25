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

type DNSMonitorControl struct {
	Cancel context.CancelFunc
	Data  models.DNSMonitor
}

func StartDNSMonitoring(ctx context.Context, db *pgxpool.Pool) {
	monitoringMap := make(map[int]DNSMonitorControl)

	for {
		monitors, err := database.GetAllDNSMonitors(ctx, db)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch DNS monitors: %v", err)
		} else {
			validsIds := make(map[int]bool)

			for _, dnsMonitor := range monitors {
				validsIds[dnsMonitor.ID] = true

				existingMonitor, exists := monitoringMap[dnsMonitor.ID]

				if exists && hasDNSConfigChanged(existingMonitor.Data, *dnsMonitor) {
					existingMonitor.Cancel()
					delete(monitoringMap, dnsMonitor.ID)
					log.Printf("[INFO] Configuration changed for domain: %s", dnsMonitor.Domain)
					exists = false
				}

				if !exists {
					monitorCtx, cancel := context.WithCancel(ctx)

					monitoringMap[dnsMonitor.ID] = DNSMonitorControl{
						Cancel: cancel,
						Data:  *dnsMonitor,
					}
					go verifyDNS(monitorCtx, db, dnsMonitor)
					log.Printf("[INFO] Started DNS monitoring for domain: %s", dnsMonitor.Domain)
				}
			}

			for id, control := range monitoringMap {
				if !validsIds[id] {
					control.Cancel()
					delete(monitoringMap, id)
					log.Printf("[INFO] Stopped DNS monitoring for domain: %s", control.Data.Domain)
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func verifyDNS(ctx context.Context, db *pgxpool.Pool, dnsMonitor *models.DNSMonitor) {
	ticker := time.NewTicker(dnsMonitor.Interval)
	defer ticker.Stop()

	for {
	a, aaaa, mx, ns := getDNSRecords(dnsMonitor.Domain)

	changed := isDifferent(dnsMonitor.LastA, a) || isDifferent(dnsMonitor.LastAAAA, aaaa) || isDifferent(dnsMonitor.LastMX, mx) || isDifferent(dnsMonitor.LastNS, ns)

	if changed {
		log.Printf("[INFO] DNS records changed for domain %s", dnsMonitor.Domain)

		err := database.UpdateDNSMonitorRecords(ctx, db, dnsMonitor.ID, a, aaaa, mx, ns)
		if err != nil {
			log.Printf("[ERROR] Failed to update DNS monitor records for domain %s: %v", dnsMonitor.Domain, err)
			return
		} else {
			dnsMonitor.LastA = a
			dnsMonitor.LastAAAA = aaaa
			dnsMonitor.LastMX = mx
			dnsMonitor.LastNS = ns
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

	select {
		case <-ctx.Done():
		return
		case <-ticker.C:
		continue
	}
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

func hasDNSConfigChanged(old, new models.DNSMonitor) bool {
	if old.Domain != new.Domain || old.Interval != new.Interval {
		return true
	}
	return false
}
