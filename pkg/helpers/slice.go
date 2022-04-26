package helpers

func IndexOf(values []string) func(value string) int {
	return func(value string) int {
		for i, element := range values {
			if value == element {
				return i
			}
		}
		return -1
	}
}
