package mpnasne

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ayumu83s/go-nasne/nasne"
	mp "github.com/mackerelio/go-mackerel-plugin"
)

// NasnePlugin is mackerel plugin
type NasnePlugin struct {
	prefix      string
	nasneClient *nasne.Client
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p *NasnePlugin) MetricKeyPrefix() string {
	if p.prefix == "" {
		p.prefix = "nasne"
	}
	return p.prefix
}

// GraphDefinition interface for mackerelplugin
// ここでグラフの定義
func (p *NasnePlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.prefix)
	return map[string]mp.Graphs{
		"recorded_num": {
			Label: labelPrefix + " Recorded Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_count", Label: "Total Count"},
			},
		},
		"record_fail_num": {
			Label: labelPrefix + " Record Fail Num",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "total_count", Label: "Total Count", Diff: true},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
// ここでグラフに渡すデータのfetch
func (p *NasnePlugin) FetchMetrics() (map[string]float64, error) {
	ret := make(map[string]float64)

	// 録画件数
	ret["total_count"] = 0
	recordedCount, err := p.getRecordedCount()
	if err != nil {
		return nil, err
	}
	ret["total_count"] = recordedCount

	// 空き容量

	// 録画失敗の件数
	ret["record_fail_num"] = 0
	recordFailNum, err := p.getRecordFailNum()
	if err != nil {
		return nil, err
	}
	ret["record_fail_num"] = recordFailNum

	return ret, nil
}

func (p *NasnePlugin) getRecordedCount() (float64, error) {
	args := &nasne.RecordedTitleListArgs{
		RequestedCount: 1,
	}
	titleList, err := p.nasneClient.Recorded.TitleListGet(nil, args)
	if err != nil {
		// エラーだったどうするのが正解なんだろう
		fmt.Errorf("fail to TitleListGet: %s")
		return 0, err
	}
	return float64(titleList.TotalMatches), nil
}

func (p *NasnePlugin) getRecordFailNum() (float64, error) {
	recNgList, err := p.nasneClient.Status.RecNgListGet(nil)
	if err != nil {
		fmt.Errorf("fail to RecNgListGet: %s")
		return 0, err
	}
	return float64(recNgList.Number), nil
}

// Do the plugin
func Do() {
	var (
		optPrefix = flag.String("metric-key-prefix", "", "Metric key prefix")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION] NASNE IP \n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	nasneClient, err := nasne.NewClient(flag.Args()[0], nil)
	if err != nil {
		fmt.Errorf("fail to nasne client: %s")
		os.Exit(1)
	}

	mp.NewMackerelPlugin(&NasnePlugin{
		prefix:      *optPrefix,
		nasneClient: nasneClient,
	}).Run()
}
