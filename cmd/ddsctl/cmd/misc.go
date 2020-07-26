package cmd

const (
	Logo = `       __    __          __  __
  ____/ /___/ /_________/ /_/ /
 / __  / __  / ___/ ___/ __/ /
/ /_/ / /_/ (__  ) /__/ /_/ /
\__,_/\__,_/____/\___/\__/_/`

	Template = Logo +
		"\n\n" +
		"{{with .Long}}{{. | trimTrailingWhitespaces}}\n\n" +
		"{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}"
)
