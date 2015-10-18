package formatter

func SelectFormatter(s string) Formatter {
	switch s {
	case "tap":
		return TapFormatter{}
	default:
		return DefaultFormatter{}
	}
}
