package ast

func GetSlice[T any](prog *Program) []T {
	var values []T

	for _, v := range prog.Nodes {
		if t, ok := v.(T); ok {
			values = append(values, t)
		}
	}

	return values
}
