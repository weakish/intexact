package intexact

import (
	"strconv"
	"testing"
	"testing/quick"
)

var incTests = []struct {
	in int
	result int
	saturated Saturated
}{
	{0, 1, false},
	{1, 2, false},
	{-1, 0, false},
	{128, 129, false},
	{-256, -255, false},
	{MinInt, MinInt+1, false},
	{MaxInt-1, MaxInt, false},
	{MaxInt, MaxInt, true},
}

func TestSaturatedInc(t *testing.T) {
	for i, test := range incTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var result int
			var saturated Saturated
			result, saturated = SaturatedInc(test.in)
			if result != test.result || saturated != test.saturated {
				t.Errorf("EXPECTED %d, %t\nACTUAL %d, %t\n",
					test.result, test.saturated, result, saturated)
			} else {}
		})
	}
}

var decTests = []struct {
	in int
	result int
	saturated Saturated
}{
	{0, -1, false},
	{1, 0, false},
	{-1, -2, false},
	{128, 127, false},
	{-256, -257, false},
	{MinInt, MinInt, true},
	{MinInt+1, MinInt, false},
	{MaxInt, MaxInt-1, false},
}

func TestSaturatedDec(t *testing.T) {
	for i, test := range decTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var result int
			var saturated Saturated
			result, saturated = SaturatedDec(test.in)
			if result != test.result || saturated != test.saturated {
				t.Errorf("EXPECTED %d, %t\nACTUAL %d, %t\n",
					test.result, test.saturated, result, saturated)
			} else {}
		})
	}
}

var negTests = []struct {
	in int
	result int
	integerOverflow error
}{
	{0, 0, nil},
	{1, -1, nil},
	{-1, 1, nil},
	{64, -64, nil},
	{MaxInt, -MaxInt, nil},
	{MinInt, 0, IntegerOverflow},
}

func TestNeg(t *testing.T) {
	for i, test := range negTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var result int
			var err error
			result, err = Neg(test.in)
			checkOverflow(test.result, test.integerOverflow, result, err, t)
		})
	}
}

func checkOverflow(expectedResult int, expectedError  error,
			  actualResult int, actualError error,
			  t *testing.T) {
	if actualError == nil {
		if expectedError == nil {
			if actualResult != expectedResult {
				t.Errorf("EXPECTED %d\nACTUAL %d\n",
					expectedResult, actualResult)
			} else {
			}
		} else {
			t.Errorf("EXPECTED overflow\nACTUAL no overflow (got %d)",
				actualResult)
		}
	} else {
		if actualError != expectedError {
			t.Errorf("EXPECTED %d, %s\nACTUAL %d, %s\n",
				expectedResult, expectedError, actualResult, actualError)
		} else {
		}
	}
}

type arithmeticTest struct {
	x int
	y int
	r int
	e error
}

func testArithmetic(tests []arithmeticTest, operation func(int, int) (int, error), t *testing.T) {
	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var result int
			var err error
			result, err = operation(test.x, test.y)
			checkOverflow(test.r, test.e, result, err, t)
		})
	}
}

var addTests = []arithmeticTest{
	{0, 0, 0, nil},
	{1, 1, 2, nil},
	{-1, -1, -2, nil},
	{MinInt, MaxInt, -1, nil},
	{MaxInt, MaxInt, 0, IntegerOverflow},
	{MinInt, MinInt, 0, IntegerOverflow},
}

func TestAdd(t *testing.T) {
	testArithmetic(addTests, Add, t)
}

var subTests = []arithmeticTest{
	{0, 0, 0, nil},
	{1, 1, 0, nil},
	{-1, -1, 0, nil},
	{2, 1, 1, nil},
	{MinInt, MaxInt, 0, IntegerOverflow},
	{MaxInt, MaxInt, 0, nil},
	{MinInt, MinInt, 0, nil},
}

func TestSub(t *testing.T) {
	testArithmetic(subTests, Sub, t)
}

type True bool

func TestSaturatedDecInversesSaturatedInc(t *testing.T) {
	saturatedDecInversesSaturatedInc := func(n int) True {
		var result int
		var saturated Saturated
		result, saturated = SaturatedInc(n)

		var r int
		var s Saturated
		r, s = SaturatedDec(result)

		if saturated {
			return s == false && r == MaxInt - 1
		} else {
			return r == n
		}
	}
	var err error = quick.Check(saturatedDecInversesSaturatedInc, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestSaturatedIncInversesSaturatedDec(t *testing.T) {
	saturatedIncInversesSaturatedDec := func(n int) True {
		var result int
		var saturated Saturated
		result, saturated = SaturatedDec(n)

		var r int
		var s Saturated
		r, s = SaturatedInc(result)

		if saturated {
			return s == false && r == MinInt + 1
		} else {
			return r == n
		}
	}
	var err error = quick.Check(saturatedIncInversesSaturatedDec, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestNegNeg(t *testing.T) {
	negNeg := func(n int) True {
		var result int
		var integerOverflow error
		result, integerOverflow = Neg(n)

		if integerOverflow == IntegerOverflow {
			return true
		} else {
			var ret int
			var err error
			ret, err = Neg(result)
			return err == nil && ret == n
		}

	}
	var err error = quick.Check(negNeg, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestIncViaAdd(t *testing.T) {
	incViaAdd := func(n int) (int, error) {
		return Add(n, 1)
	}
	var err error = quick.CheckEqual(Inc, incViaAdd, nil)
	if err != nil {
		t.Error(err)
	}
}

func TestDecViaSub(t *testing.T) {
	decViaSub := func(n int) (int, error) {
		return Sub(n, 1)
	}
	var err error = quick.CheckEqual(Dec, decViaSub, nil)
	if err != nil {
		t.Error(err)
	}
}