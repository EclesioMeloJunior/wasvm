package collections

type Iterator[T any] interface {
	Next() (T, error)
}
