package fp

func Mapper[T, U any](s []T, f func(T) U) []U {
	res := make([]U, len(s))
	for i := range len(s) {
		res[i] = f(s[i])
	}
	return res
}
