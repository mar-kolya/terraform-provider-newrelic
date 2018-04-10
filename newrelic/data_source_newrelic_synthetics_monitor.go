package newrelic

import (
	"fmt"
	"log"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	"github.com/hashicorp/terraform/helper/schema"
)

const monitorsRequestLimit = 100

func dataSourceNewRelicSyntheticsMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicSyntheticsMonitorRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_KEY", nil),
				Sensitive:   true,
			},
		},
	}
}

func dataSourceNewRelicSyntheticsMonitorRead(d *schema.ResourceData, meta interface{}) error {

	apiKey := d.Get("api_key").(string)

	conf := func(s *synthetics.Client) {
		s.APIKey = apiKey
	}

	syntheticsClient, _ := synthetics.NewClient(conf)

	log.Printf("[INFO] Reading New Relic synthetics monitors")

	name := d.Get("name").(string)
	offset := uint(0)
	for {
		monitors, err := syntheticsClient.GetAllMonitors(offset, monitorsRequestLimit)
		if err != nil {
			return err
		}

		for _, monitor := range monitors.Monitors {
			if monitor.Name == name {
				d.SetId(monitor.ID)
				d.Set("name", monitor.Name)
				d.Set("monitor_id", monitor.ID)
				return nil
			}
		}

		if len(monitors.Monitors) < monitorsRequestLimit {
			break
		}

		offset = offset + 100
	}

	return fmt.Errorf("The name '%s' does not match any New Relic monitors.", name)
}
