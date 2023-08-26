package sliceutil

func Mapper[I, O any](list []I, f func(I) O) []O {
	var results []O

	for _, item := range list {
		results = append(results, f(item))
	}

	return results
}

func Filter[I any](list []I, f func(I) bool) []I {
	var results []I

	for _, item := range list {
		if f(item) {
			results = append(results, item)
		}
	}

	return results
}

func Reduce[I, O any](list []I, fn func(O, I, int) O, initial O) O {
	result := initial

	for i, item := range list {
		result = fn(result, item, i)
	}

	return result
}
