package geom

import "testing"

func TestInitCoordSeq(t *testing.T) {
	seq, _ := initCoordSeq(2, 3)
	defer seq.Destroy()

	seq.SetX(0, 10)
	seq.SetY(0, 30)
	seq.SetZ(0, 50)

	seq.SetX(1, 20)
	seq.SetY(1, 50)
	seq.SetZ(1, 70)

	size := seq.GetSize()

	if size != 2 {
		t.Errorf("Error: initCoordSeq()")
	}
}

func TestSetX(t *testing.T) {
	seq, _ := initCoordSeq(1, 2)
	defer seq.Destroy()

	seq.SetX(0, 10)

	x := seq.GetX(0)

	if x != 10 {
		t.Errorf("Error: GetX()")

	}
}

func TestSetY(t *testing.T) {
	seq, _ := initCoordSeq(1, 2)
	defer seq.Destroy()

	seq.SetY(0, 30)

	y := seq.GetY(0)

	if y != 30 {
		t.Errorf("Error: GetY()")

	}
}

func TestSetZ(t *testing.T) {
	seq, _ := initCoordSeq(1, 2)
	defer seq.Destroy()

	seq.SetZ(0, 50)

	y := seq.GetZ(0)

	if y != 50 {
		t.Errorf("Error: GetZ()")

	}
}
