package gateway

func selectBestMECServer(aim map[string]string) string {
	for _, IP := range aim {
		return IP
	}
	return ""
}
