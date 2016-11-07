package blockchain

func arrayContains(function string, list []string) bool {
	for _, v := range list {
		if v == function {
			return true
		}
	}
	return false
}
