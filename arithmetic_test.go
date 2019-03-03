package intexact

import (
	"strconv"
	"testing"
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
			if err == nil {
				if err == test.integerOverflow {
					if result != test.result {
						t.Errorf("EXPECTED %d\nACTUAL %d\n",
							test.result, result)
					} else {}
				} else {
					t.Errorf("EXPECTED overflow\nACTUAL no overflow (got %d)",
						result)
				}
			} else {
				if err != test.integerOverflow {
					t.Errorf("EXPECTED %d, %s\nACTUAL %d, %s\n",
						test.result, test.integerOverflow, result, err)
				} else {}
			}
		})
	}
}