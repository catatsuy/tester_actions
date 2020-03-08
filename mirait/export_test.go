package mirait

func SetTargetURL(u string) (resetFunc func()) {
	var tmp string
	tmp, targetURL = targetURL, u
	return func() {
		targetURL = tmp
	}
}
