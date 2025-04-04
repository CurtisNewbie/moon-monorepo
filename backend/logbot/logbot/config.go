package logbot

import "github.com/curtisnewbie/miso/miso"

// misoconfig-section: Logbot Configuration
const (
	// misoconfig-prop: logbot node name | default
	PropLogbotNode = "logbot.node"

	// misoconfig-prop: (`slice of watch object`) logbot watch configuration |
	PropLogbotWatch = "logbot.watch"

	// misoconfig-prop: (`watch object`) app name |
	PropLogbotWatchApp = "logbot.watch.app"

	// misoconfig-prop: (`watch object`) path of the log file |
	PropLogbotWatchFile = "logbot.watch.file "

	// misoconfig-prop: (`watch object`) type of log pattern `[ 'go', 'java' ]` |
	PropLogbotWatchType = "logbot.watch.type"

	// misoconfig-prop: enable task to remove error logs reported 7 days ago | false
	PropEnableRemoveErrLogTask = "logbot.remove-history-error-log"

	// misoconfig-prop: (`slice of string`) log pattern supported (regexp) |
	PropPattern = "log.pattern"

	// misoconfig-prop: merged log filename |
	PropMergedLogFilename = "log.merged-file-name"
)

// misoconfig-default-start
func init() {
	miso.SetDefProp(PropLogbotNode, "default")
	miso.SetDefProp(PropEnableRemoveErrLogTask, false)
}

// misoconfig-default-end

type WatchConfig struct {
	App         string
	File        string
	Type        string
	ReportError bool `mapstructure:"report-error"`
}

type LogBotConfig struct {
	NodeName     string        `mapstructure:"node"`
	WatchConfigs []WatchConfig `mapstructure:"watch"`
}

type Config struct {
	Config LogBotConfig `mapstructure:"logbot"`
}

func LoadLogBotConfig() Config {
	var conf Config
	miso.UnmarshalFromProp(&conf)
	return conf
}

func IsRmErrorLogTaskEnabled() bool {
	return miso.GetPropBool(PropEnableRemoveErrLogTask)
}
