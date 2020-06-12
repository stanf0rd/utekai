package main

func filterIntArray(source []int, removing []int) []int {
	var result []int

	for _, value := range source {
		if !intInSlice(value, removing) {
			result = append(result, value)
		}
	}

	return result
}

func intInSlice(desired int, list []int) bool {
	for _, value := range list {
		if desired == value {
			return true
		}
	}
	return false
}

func stringInSlice(desired string, list []string) bool {
	for _, value := range list {
		if value == desired {
			return true
		}
	}
	return false
}
