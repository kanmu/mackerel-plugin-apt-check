package apt_check

import (
	"bytes"
	"flag"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

const scriptPathDefault = "/usr/lib/update-notifier/apt-check"

var logger = logging.GetLogger("metrics.plugin.apt-check")

// AptCheckPlugin stores the parameters for apt-check Mackerel plugin.
type AptCheckPlugin struct {
	ScriptPath string // path to apt-check script
	Prefix     string
}

// aptCheckResult stores the result of invocation of apt-check script.
type aptCheckResult struct {
	NumOfUpdates         int
	NumOfSecurityUpdates int
}

// MetricKeyPrefix returns the metrics key prefix
func (p AptCheckPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "apt-check"
	}
	return p.Prefix
}

func (q AptCheckPlugin) invokeAptCheck() (res aptCheckResult, err error) {
	cmd := exec.Command(q.ScriptPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	tokens := bytes.Split(out, []byte{';'})
	if len(tokens) != 2 {
		err = fmt.Errorf("invalid output: %s", string(out))
	}

	res.NumOfUpdates, err = strconv.Atoi(string(tokens[0]))
	if err != nil {
		return
	}

	res.NumOfSecurityUpdates, err = strconv.Atoi(string(tokens[1]))
	if err != nil {
		return
	}

	return
}

// FetchMetrics interface for mackerelplugin
func (q AptCheckPlugin) FetchMetrics() (map[string]interface{}, error) {
	res, err := q.invokeAptCheck()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"updates":          uint64(res.NumOfUpdates),
		"security_updates": uint64(res.NumOfSecurityUpdates),
	}, nil
}

// GraphDefinition interface for mackerelplugin
func (p AptCheckPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.MetricKeyPrefix())

	var graphdef = map[string]mp.Graphs{
		"updates": {
			Label: (labelPrefix + " Updates"),
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "updates", Label: "Available updates", Type: "uint64", Diff: false, Stacked: false},
				{Name: "security_updates", Label: "Available security updates", Type: "uint64", Diff: false, Stacked: false},
			},
		},
	}

	return graphdef
}

// Do the plugin
func Do() {
	optScriptPath := flag.String("script", scriptPathDefault, "Path to apt-check")
	optPrefix := flag.String("metric-key-prefix", "apt-check", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	var a AptCheckPlugin
	a.ScriptPath = *optScriptPath
	a.Prefix = *optPrefix

	helper := mp.NewMackerelPlugin(a)

	helper.Tempfile = *optTempfile
	helper.Run()
}
