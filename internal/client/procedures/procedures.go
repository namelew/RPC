package procedures

type Number interface {
	float32 | float64
}

func Add[T Number](a ...T) T {
	var sum T = 0

	for i := range a {
		sum += a[i]
	}

	return sum
}

func Sub[T Number](a ...T) T {
	var result T = a[0]

	for i := 1; i < len(a); i++ {
		result -= a[i]
	}

	return result
}
