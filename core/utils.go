package core

func Contains(S []string, E string) bool {
	for _, s := range S {
		if s == E {
			return true
		}
	}
	return false
}
