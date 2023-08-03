package stringx

// * js effective?effective:s
func TakeOne(effective string, s string) string {
	if len(s) > 0 {
		return effective
	}
	return s
}
