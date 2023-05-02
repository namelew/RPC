package procedures

func Add(a ...int64) int64 {
	var sum int64 = 0

	for i := range a {
		sum += a[i]
	}

	return sum
}

func Sub(a ...int64) int64 {
	var result int64 = a[0]

	for i := 1; i < len(a); i++ {
		result -= a[i]
	}

	return result
}
