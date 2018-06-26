package pkg

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/yiqinguo/armyant/pkg/models"
)

func ParseArgs() models.Job {
	var job models.Job
	kingpin.Flag("id", "job id").Short('i').Required().StringVar(&job.Id)
	kingpin.Flag("hid", "exec history id").Short('I').Required().StringVar(&job.HistoryId)
	kingpin.Flag("url", "url").Required().Short('u').StringVar(&job.Url)
	kingpin.Flag("conn", "Maximum number of concurrent connections").
		Short('c').Default("125").Int64Var(&job.Connections)
	kingpin.Flag("requests", "[pos. int.]  Number of requests").
		Short('n').Int64Var(&job.Requests)
	kingpin.Flag("duration", "Duration of test").Short('d').StringVar(&job.Duration)
	kingpin.Flag("method", "Request method").Default("GET").Short('m').StringVar(&job.Method)
	kingpin.Flag("body", "Request body").Short('b').StringVar(&job.Body)
	kingpin.Flag("master", "master grpc address").Short('M').StringVar(&job.Master)
	kingpin.Flag("header", "HTTP headers to use(can be repeated)").
		Short('H').StringMapVar(&job.Headers)

	kingpin.Parse()
	return job
}
