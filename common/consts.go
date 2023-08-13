package common

const (
	// flags:
	AboutFlag      = "about"
	AllFlag        = "all"
	AmFlag         = "am"
	AscendingFlag  = "asc"
	ClientFlag     = "client"
	DescendingFlag = "desc"
	DirFlag        = "dir"
	HoursFlag      = "hr"
	IdFlag         = "id"
	MessageFlag    = "message"
	MinutesFlag    = "min"
	PmFlag         = "pm"
	PortFlag       = "port"
	PostponeFlag   = "postpone"
	SecondsFlag    = "sec"
	ServerFlag     = "server"
	SortFlag       = "sort"
	TimeFlag       = "time"

	// time format:
	DateTimeFormatWithoutTimeZone = "2006-01-02 15:04:05"
	TimeFormat12AmPmHours         = "03:04 PM"
	TimeFormat24Hours             = "15:04:05"

	// OS:
	WindowsOS = "windows"
	LinuxOS   = "linux"
	MacOS     = "darwin"

	// Shell:
	BashShell = "bash"
	ZshShell  = "zsh"
	FishShell = "fish"

	// configs:
	AdminConfigsFileName  = "remindme_admin_configs.yaml"
	ClientLogsFileName    = "remindme_client_logs.log"
	DefaultHttpServerPort = 15555
	ServerPortEnvVar      = "REMINDME_SERVER_PORT"
	ServerLogsFileName    = "remindme_server_logs.log"
)
