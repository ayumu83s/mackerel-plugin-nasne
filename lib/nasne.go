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

// グラフ定義で取得した値を保持して使いまわす
var hddInfoList [](*nasne.HDDInfo)

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

	err := p.fetchHddInfoList()
	if err != nil {
		fmt.Errorf("fail to fetchHddList: %s")
		return nil
	}

	graphs := 2 + len(hddInfoList)
	ret := make(map[string]mp.Graphs, graphs)

	for _, hddInfo := range hddInfoList {
		graphsKey := fmt.Sprintf("HDD_%d", hddInfo.Hdd.ID)
		usedMetricsKey := fmt.Sprintf("used_%d", hddInfo.Hdd.ID)
		freeMetricsKey := fmt.Sprintf("free_%d", hddInfo.Hdd.ID)
		totalMetricsKey := fmt.Sprintf("total_%d", hddInfo.Hdd.ID)

		ret[graphsKey] = mp.Graphs{
			Label: fmt.Sprintf("%s HDD %d", labelPrefix, hddInfo.Hdd.ID),
			Unit:  "bytes",
			Metrics: []mp.Metrics{
				{Name: usedMetricsKey, Label: "used", Stacked: true},
				{Name: freeMetricsKey, Label: "free", Stacked: true},
				{Name: totalMetricsKey, Label: "total", Stacked: true},
			},
		}
	}

	ret["recorded_num"] = mp.Graphs{
		Label: labelPrefix + " Recorded Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "recorded_count", Label: "Recorded Count"},
		},
	}
	ret["record_fail_num"] = mp.Graphs{
		Label: labelPrefix + " Record Fail Num",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "record_fail_count", Label: "Record Fail Count", Diff: true},
		},
	}
	return ret
}

// FetchMetrics interface for mackerelplugin
// ここでグラフに渡すデータのfetch
func (p *NasnePlugin) FetchMetrics() (map[string]float64, error) {
	ret := make(map[string]float64)

	// 空き容量
	err := p.fetchHddInfoList()
	if err != nil {
		return nil, err
	}

	for _, hddInfo := range hddInfoList {
		usedMetricsKey := fmt.Sprintf("used_%d", hddInfo.Hdd.ID)
		freeMetricsKey := fmt.Sprintf("free_%d", hddInfo.Hdd.ID)
		totalMetricsKey := fmt.Sprintf("total_%d", hddInfo.Hdd.ID)

		ret[usedMetricsKey] = float64(hddInfo.Hdd.UsedVolumeSize)
		ret[freeMetricsKey] = float64(hddInfo.Hdd.FreeVolumeSize)
		ret[totalMetricsKey] = float64(hddInfo.Hdd.TotalVolumeSize)
	}

	// 録画件数
	ret["recorded_count"] = 0
	recordedCount, err := p.getRecordedCount()
	if err != nil {
		return nil, err
	}
	ret["recorded_count"] = recordedCount

	// 録画失敗の件数
	ret["record_fail_count"] = 0
	recordFailNum, err := p.getRecordFailNum()
	if err != nil {
		return nil, err
	}
	ret["record_fail_count"] = recordFailNum

	return ret, nil
}

func (p *NasnePlugin) fetchHddInfoList() error {
	// すでにfetch済みならreturn
	if hddInfoList != nil {
		return nil
	}

	// 一覧を取得
	hddList, err := p.nasneClient.Status.HDDListGet(nil)
	if err != nil {
		fmt.Errorf("fail to HDDListGet: %s")
		return err
	}

	// 各詳細を取得
	hddCount := hddList.Number
	hddInfoList = make([](*nasne.HDDInfo), hddCount)
	for i, hdd := range hddList.Hdd {
		fmt.Println(i)
		hddInfoList[i], err = p.nasneClient.Status.HDDInfoGet(nil, hdd.ID)
		if err != nil {
			fmt.Errorf("fail to HDDInfoGet(%d): %s", hdd.ID)
			return err
		}
	}
	return nil
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
