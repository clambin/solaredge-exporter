package exporter

import (
	"github.com/clambin/solaredge-exporter/pkg/solaredge"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func Run(apiKey string, interval time.Duration) (err error) {
	client := solaredge.NewClient(apiKey, nil)

	var sites []int
	sites, err = client.GetSiteIDs()

	if err != nil {
		return err
	}

	for {
		for _, site := range sites {
			var current float64

			_, _, _, _, current, err = client.GetPowerOverview(site)

			if err == nil {
				currentPower.WithLabelValues(strconv.Itoa(site)).Set(current)
			} else {
				log.WithError(err).Warning("unable to get power stats")
			}
		}

		time.Sleep(interval)
	}
}
