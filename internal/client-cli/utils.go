package clientcli

func truncateWithEllipsis(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return s[:maxLength]
	}
	return s[:maxLength-3] + "..."
}

func truncateProjectDescription(desc string) string {
	return truncateWithEllipsis(desc, 30)
}

func truncateProjectName(name string) string {
	return truncateWithEllipsis(name, 20)
}
