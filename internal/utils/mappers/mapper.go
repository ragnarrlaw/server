package mappers

/**
  For type conversions a mapper function can be used
  Most common mapper types goes here
*/

type MapFunc[T any, U any] func(v T) U

/**
Mapper a basic mapping between the domain of T and the taget domain V
	Example usage :=
	var stringToInt MapFunc[string, int] = MapFunc[string, int](func (v string) int {
	  // implement the conversion here
	  return 0
	})
*/

func (m MapFunc[T, U]) Map(v T) U {
	return m(v)
}

/**
For mapping a collection of values from the domain of T the domain of U
*/

func (m MapFunc[T, U]) MapEach(v []T) []U {
	results := make([]U, len(v))
	for index, item := range v {
		results[index] = m(item)
	}
	return results
}

/**
For mapping with cases where the step prior to the conversion could fail
	convert only if an error didn't occur - it returns a base value
*/

func (m MapFunc[T, U]) MapErr(v T, err error) U {
	if err != nil {
		var zero U
		return zero
	}
	return m(v)
}

/**
For mapping a collection of values from the domain of T the domain of U
*/

func (m MapFunc[T, U]) MapEachErr(v []T, err error) []U {
	if err != nil {
		return make([]U, 0)
	} else {
		results := make([]U, len(v))
		for index, item := range v {
			results[index] = m(item)
		}
		return results
	}
}
