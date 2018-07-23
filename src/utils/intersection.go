package utils

// Intersection will do intersection
func Intersection(a, b []string) (ret []string) {
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			ret = append(ret, item)
		}
	}
	return
}

// Intersections will do many intersection
func Intersections(input [][]string) (ret []string) {
	if len(input) == 0 {
		return
	}

	ret = input[0]
	for i := 1; i < len(input); i++ {
		ret = Intersection(ret, input[i])
	}

	return
}
