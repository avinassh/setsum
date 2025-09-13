package setsum

import "testing"

func TestAddState(t *testing.T) {
	lhs := [SetsumColumns]uint32{1, 2, 3, 4, 5, 6, 7, 8}
	rhs := [SetsumColumns]uint32{2, 4, 6, 8, 10, 12, 14, 16}
	expected := [SetsumColumns]uint32{3, 6, 9, 12, 15, 18, 21, 24}
	returned := addState(lhs, rhs)
	if expected != returned {
		t.Errorf("addState() = %v, want %v", returned, expected)
	}
}

func TestAddStatePrimes(t *testing.T) {
	lhs := [SetsumColumns]uint32{3146800025, 1792545563, 417324692, 3444237760, 2812742746, 1608771649, 1661742866, 3220115897}
	rhs := [SetsumColumns]uint32{1148167266, 2502421716, 3877642539, 850729437, 1482224443, 2686195512, 2633224277, 1074851214}
	expected := [SetsumColumns]uint32{0, 0, 0, 0, 0, 0, 0, 0}
	returned := addState(lhs, rhs)
	if expected != returned {
		t.Errorf("addState() = %v, want %v", returned, expected)
	}
}

func TestInvertStateDesc(t *testing.T) {
	stateIn := [8]uint32{
		0xffffeeee, 0xddddcccc, 0xbbbbaaaa, 0x99998888,
		0x77776666, 0x66665555, 0x44443333, 0x22221111,
	}
	expected := [8]uint32{
		4365, 572666659, 1145328917, 1717991189,
		2290653487, 2576984612, 3149646900, 3722309174,
	}
	returned := invertState(stateIn)
	if returned != expected {
		t.Errorf("invertState() = %v, expected %v", returned, expected)
	}
}

func TestEmptyItemToState(t *testing.T) {
	expected := [8]uint32{
		0xf8c6ffa7, 0x66d71ebf, 0x5647c151, 0x62d661a0,
		0x4dff80f5, 0xfa493be4, 0x4b0ad882, 0x4a43f880,
	}
	returned := itemsToState([][]byte{})
	if returned != expected {
		t.Errorf("itemsToState([]) = %v, expected %v", returned, expected)
	}
}

// the seven values test straight from the rust lib
var sevenValues = [32]byte{
	197, 179, 253, 77, 1, 242, 184, 4, 15, 84, 171, 116, 18, 202, 83, 187,
	252, 153, 14, 39, 42, 64, 173, 209, 196, 206, 186, 107, 47, 228, 114, 213,
}

func TestSetsumInsert7Sorted(t *testing.T) {
	setsum := NewSetsum()
	setsum.Insert([]byte("this is the first value"))
	setsum.Insert([]byte("this is the second value"))
	setsum.Insert([]byte("this is the third value"))
	setsum.Insert([]byte("this is the fourth value"))
	setsum.Insert([]byte("this is the fifth value"))
	setsum.Insert([]byte("this is the sixth value"))
	setsum.Insert([]byte("this is the seventh value"))

	digest := setsum.Digest()
	if digest != sevenValues {
		t.Errorf("Digest() = %v, expected %v", digest, sevenValues)
	}
}

func TestSetsumInsert7Reversed(t *testing.T) {
	setsum := NewSetsum()
	setsum.Insert([]byte("this is the seventh value"))
	setsum.Insert([]byte("this is the sixth value"))
	setsum.Insert([]byte("this is the fifth value"))
	setsum.Insert([]byte("this is the fourth value"))
	setsum.Insert([]byte("this is the third value"))
	setsum.Insert([]byte("this is the second value"))
	setsum.Insert([]byte("this is the first value"))

	digest := setsum.Digest()
	if digest != sevenValues {
		t.Errorf("Digest() = %v, expected %v", digest, sevenValues)
	}
}

func TestSetsumInsert7Random(t *testing.T) {
	setsum := NewSetsum()
	setsum.Insert([]byte("this is the fifth value"))
	setsum.Insert([]byte("this is the fourth value"))
	setsum.Insert([]byte("this is the third value"))
	setsum.Insert([]byte("this is the sixth value"))
	setsum.Insert([]byte("this is the seventh value"))
	setsum.Insert([]byte("this is the second value"))
	setsum.Insert([]byte("this is the first value"))

	digest := setsum.Digest()
	if digest != sevenValues {
		t.Errorf("Digest() = %v, expected %v", digest, sevenValues)
	}
}

func TestSetsumInsertRemove(t *testing.T) {
	setsum := NewSetsum()
	setsum.Insert([]byte("this is the first value"))
	setsum.Insert([]byte("this is the second value"))
	setsum.Insert([]byte("this is the third value"))
	setsum.Insert([]byte("this is the fourth value"))
	setsum.Insert([]byte("this is the fifth value"))
	setsum.Insert([]byte("this is the sixth value"))
	setsum.Insert([]byte("this is the seventh value"))
	setsum.Remove([]byte("this is the seventh value"))
	setsum.Remove([]byte("this is the sixth value"))
	setsum.Remove([]byte("this is the fifth value"))
	setsum.Remove([]byte("this is the fourth value"))
	setsum.Remove([]byte("this is the third value"))
	setsum.Remove([]byte("this is the second value"))
	setsum.Remove([]byte("this is the first value"))

	digest := setsum.Digest()
	defaultDigest := NewSetsum().Digest()
	if digest != defaultDigest {
		t.Errorf("Digest() = %v, expected %v", digest, defaultDigest)
	}
}

func TestSetsumMergeTwoSets(t *testing.T) {
	setsumOne := NewSetsum()
	setsumOne.Insert([]byte("this is the first value"))
	setsumOne.Insert([]byte("this is the second value"))
	setsumOne.Insert([]byte("this is the third value"))
	setsumOne.Insert([]byte("this is the fourth value"))

	setsumTwo := NewSetsum()
	setsumTwo.Insert([]byte("this is the fifth value"))
	setsumTwo.Insert([]byte("this is the sixth value"))
	setsumTwo.Insert([]byte("this is the seventh value"))

	setsumOne.Add(setsumTwo)
	mergedDigest := setsumOne.Digest()
	if mergedDigest != sevenValues {
		t.Errorf("Digest() = %v, expected %v", mergedDigest, sevenValues)
	}
}

func TestSetsumRemoveTwoSets(t *testing.T) {
	baseSetsum := NewSetsum()
	baseSetsum.Insert([]byte("this is the first value"))
	baseSetsum.Insert([]byte("this is the second value"))
	baseSetsum.Insert([]byte("this is the third value"))
	baseSetsum.Insert([]byte("this is the fourth value"))
	baseSetsum.Insert([]byte("this is the fifth value"))
	baseSetsum.Insert([]byte("this is the sixth value"))
	baseSetsum.Insert([]byte("this is the seventh value"))

	setsumOne := NewSetsum()
	setsumOne.Insert([]byte("this is the first value"))
	setsumOne.Insert([]byte("this is the second value"))
	setsumOne.Insert([]byte("this is the third value"))
	setsumOne.Insert([]byte("this is the fourth value"))

	setsumTwo := NewSetsum()
	setsumTwo.Insert([]byte("this is the fifth value"))
	setsumTwo.Insert([]byte("this is the sixth value"))
	setsumTwo.Insert([]byte("this is the seventh value"))

	baseSetsum.Subtract(setsumOne)
	baseSetsum.Subtract(setsumTwo)
	emptyDigest := baseSetsum.Digest()
	defaultDigest := NewSetsum().Digest()
	if emptyDigest != defaultDigest {
		t.Errorf("Digest() = %v, expected %v", emptyDigest, defaultDigest)
	}
}
