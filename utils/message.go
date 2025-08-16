package utils


func CanonDirectKey(uid1, uid2 string) string {
	a, b := uid1, uid2
	if a > b { // lexicographic compare
		a, b = b, a
	}
	return a + ":" + b
}

