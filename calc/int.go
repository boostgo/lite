package calc

import "math"

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Abs[T Integer](value T) T {
	if value < 0 {
		return -value
	}

	return value
}

func Neg[T Integer](value T) T {
	if value > 0 {
		return -value
	}

	return value
}

func IsEven[T Integer](value T) bool {
	return value%2 == 0
}

func IsOdd[T Integer](value T) bool {
	return !IsEven(value)
}

func Compare[T1 Integer, T2 Integer](a T1, b T2) int {
	x := int64(a)
	y := int64(b)

	if x > y {
		return 1
	} else if x < y {
		return -1
	}

	return 0
}

func Pow[T Integer](base T, exp int) T {
	var result T = 1
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}

		base *= base
		exp /= 2
	}
	return result
}

func Sum[T Integer](values ...T) T {
	if len(values) == 0 {
		return 0
	}

	var sum T
	for _, value := range values {
		sum += value
	}
	return sum
}

func Factorial[T Integer](value T) T {
	if value < 0 {
		value = Abs(value)
	}

	var result T = 1
	var start T = 2
	for i := start; i <= value; i++ {
		result *= i
	}

	return result
}

func IsPrime[T Integer](value T) bool {
	if value <= 1 {
		return false
	}
	if value <= 3 {
		return true
	}
	if value%2 == 0 || value%3 == 0 {
		return false
	}

	var start T = 5
	for i := start; i <= T(math.Sqrt(float64(value))); i += 6 {
		if value%i == 0 || value%(i+2) == 0 {
			return false
		}
	}

	return true
}

func GCD[T Integer](a, b T) T {
	if b == 0 {
		return Abs(a)
	}

	return GCD[T](b, a%b)
}

func LCM[T Integer](a, b T) T {
	if a == 0 || b == 0 {
		return 0
	}

	return Abs(a*b) / GCD(a, b)
}

func Max[T Integer](nums ...T) T {
	if len(nums) == 0 {
		return 0
	}

	maxValue := nums[0]
	for _, n := range nums[1:] {
		if n > maxValue {
			maxValue = n
		}
	}

	return maxValue
}

func Min[T Integer](nums ...T) T {
	if len(nums) == 0 {
		return 0
	}

	maxValue := nums[0]
	for _, n := range nums[1:] {
		if n < maxValue {
			maxValue = n
		}
	}

	return maxValue
}

func Fibonacci[T Integer](value T) T {
	if value < 0 {
		value = Abs(value)
	}

	if value == 0 {
		return 0
	}

	if value == 1 {
		return 1
	}

	var a T = 0
	var b T = 1
	var start T = 2
	for i := start; i <= value; i++ {
		a, b = b, a+b
	}

	return b
}

func IsPalindrome[T Integer](num T) bool {
	if num < 0 {
		return false
	}

	original := num
	var reversed T = 0

	for num > 0 {
		digit := num % 10
		reversed = reversed*10 + digit
		num /= 10
	}

	return original == reversed
}
