package randomlist

import "testing"

func TestZeroLength(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(`Accessing an element of a zero length list did not panic`)
		}
	}()

	x := New[byte](nil)
	if x.Len() != 0 {
		t.Fatal(`Zero elements slice does not have zero length`)
	}

	x.RandomElement()
}

func TestInt(t *testing.T) {
	for testLen := 1; testLen <= 100; testLen++ {
		testSlice := make([]int, testLen)
		for i := 0; i < testLen; i++ {
			testSlice[i] = i
		}
		testRandomList := New(testSlice)
		count := 0
		for i := 0; i < testLen; i++ {
			n := testRandomList.RandomElement()
			if n == i {
				count++
			}
		}
		if testLen > 4 && count == testLen {
			t.Fatal(`Random list does not have any randomness`)
		}
	}
}

func TestFloat64(t *testing.T) {
	for testLen := 1; testLen <= 100; testLen++ {
		testSlice := make([]float64, testLen)
		for i := 0; i < testLen; i++ {
			testSlice[i] = float64(i)
		}
		testRandomList := New(testSlice)
		count := 0
		for i := 0; i < testLen; i++ {
			n := testRandomList.RandomElement()
			if n == float64(i) {
				count++
			}
		}
		if testLen > 4 && count == testLen {
			t.Fatal(`Random list does not have any randomness`)
		}
	}
}
