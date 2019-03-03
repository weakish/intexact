// Package intexact provides checked basic arithmetic operations for int.

package intexact

import "errors"

const maxUint = ^uint(0)
const MaxInt = int(maxUint >> 1)
const MinInt = -MaxInt - 1

var IntegerOverflow error = errors.New("message")
type Saturated bool

func Inc(n int) (int, error) {
	if n == MaxInt {
		return MaxInt, IntegerOverflow
	} else {
		return n + 1, nil
	}
}

func SaturatedInc(n int) (int, Saturated) {
	var result int
	result, err := Inc(n)
	if err == IntegerOverflow {
		return result, true
	} else {
		return result, false
	}
}

func Dec(n int) (int, error) {
	if n == MinInt {
		return MinInt, IntegerOverflow
	} else {
		return n - 1, nil
	}
}

func SaturatedDec(n int) (int, Saturated) {
	var result int
	result, err := Dec(n)
	if err == IntegerOverflow {
		return result, true
	} else {
		return result, false
	}
}

func Neg(n int) (int, error) {
	if n == MinInt {
		return MaxInt, IntegerOverflow
	} else {
		return -n, nil
	}
}


func Add(x int, y int) (int, error) {
	var r int = x + y
	// x and y have the opposite sign of the result
	if ((x ^ r) & (y ^ r)) < 0 {
		return 0, IntegerOverflow
	} else {
		return r, nil
	}
}

func Sub(x int, y int) (int, error) {
	var r int = x - y
	// x and y have different signs and x and result have different sign
	if ((x ^ y) & (x ^ r)) < 0 {
		return 0, IntegerOverflow
	} else {
		return r, nil
	}
}

func Mul(x int, y int) (int, error) {
	// inspired by Rob Pike
	// https://groups.google.com/forum/#!msg/golang-nuts/h5oSN5t3Au4/KaNQREhZh0QJ
	if x == 0 || y == 0 {
		return 0, nil
	} else if x == 1 {
		return y, nil
	} else if y == 1 {
		return x, nil
	} else if x == MinInt || y == MinInt {
		return 0, IntegerOverflow
	} else {
		var r int = x * y
		if r / y == x {
			return r, nil
		} else {
			return 0, IntegerOverflow
		}
	}
}