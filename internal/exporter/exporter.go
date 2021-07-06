package exporter

import (
	"context"
	"github.com/clambin/solaredge"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func Run(ctx context.Context, client solaredge.API, interval time.Duration) (err error) {
	var sites []int
	sites, err = client.GetSiteIDs(ctx)

	if err == nil {
		for {
			for _, site := range sites {
				var current float64

				_, _, _, _, current, err = client.GetPowerOverview(ctx, site)

				if err == nil {
					currentPower.WithLabelValues(strconv.Itoa(site)).Set(current)
				} else {
					log.WithError(err).Warning("unable to get power statistics")
				}
			}

			time.Sleep(interval)
		}
	}

	return
}
