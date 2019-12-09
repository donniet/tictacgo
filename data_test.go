package tictacgo

import "testing"

func TestPosition(t *testing.T) {
	pGood := []string{
		"XXO\n XO\n X \n",
		"   \n   \n   \n",
		"XXX\nXXX\nXXX\n",
		"OOO\nOOO\nOOO\n",
	}
	pBad := []string{
		"",
		"\n",
		"X\n",
		"XXXX",
		"XOX\nOXO\nXOX\nX",
	}

	for _, pos := range pGood {
		_, err := FromString(pos)

		if err != nil {
			t.Errorf("expected valid string: '%s' got error '%s'", pos, err)
		}
	}
	for _, pos := range pBad {
		_, err := FromString(pos)

		if err == nil {
			t.Errorf("expected invalid string: '%s'", pos)
		}
	}

	p := Position{}
	p.Set(0, 1, X)
	if p.Get(0, 1) != X {
		t.Error("Set didn't work")
	}
	q := p
	p.Set(-1, 0, X)
	if p != q {
		t.Error("out of bounds should do nothing")
	}
	if p.Get(4, 0) != Empty {
		t.Error("out of bounds should be empty")
	}

	p.Clear()
	if p.Get(0, 1) != Empty {
		t.Error("clear didn't work")
	}

	p.Set(0, 1, X)

	q = p.FlipX()
	if q.Get(0, 1) != Empty || q.Get(2, 1) != X {
		t.Log(p, q)
		t.Error("flipX doesn't work")
	}

	p.Set(1, 0, O)
	q = p.FlipY()
	if q.Get(1, 0) != Empty || q.Get(1, 2) != O {
		t.Error("flipY doesn't work")
	}

	// p should be:
	//  |O|
	// X| |
	//  | |

	// these should all be the same
	q = p.Rotate(1)
	if r, _ := FromString(" X \n  O\n   \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}
	q = p.Rotate(-3)
	if r, _ := FromString(" X \n  O\n   \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}
	q = p.Rotate(5)
	if r, _ := FromString(" X \n  O\n   \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}

	// and these
	q = p.Rotate(2)
	if r, _ := FromString("   \n  X\n O \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}
	q = p.Rotate(-2)
	if r, _ := FromString("   \n  X\n O \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}
	q = p.Rotate(6)
	if r, _ := FromString("   \n  X\n O \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}

	// and these
	q = p.Rotate(3)
	if r, _ := FromString("   \nO  \n X \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}
	q = p.Rotate(-1)
	if r, _ := FromString("   \nO  \n X \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}
	q = p.Rotate(7)
	if r, _ := FromString("   \nO  \n X \n"); q != r {
		t.Error("Rotate 1 doesn't work")
	}

	// idempotent
	q = p.Rotate(0)

	if p != q {
		t.Error("rotate of zero should be idempotent")
	}
}

func TestWinners(t *testing.T) {
	winners := map[Square][]string{
		Empty: []string{
			"   \n   \n   \n",
			"X  \nX  \nO  \n",
		},
		X: []string{
			"XXX\n   \n   \n",
			"X  \nX  \nX  \n",
			"X  \n X \n  X\n",
		},
		O: []string{
			"OOO\n   \n   \n",
			"O  \nO  \nO  \n",
			"O  \n O \n  O\n",
		},
	}

	for s, wins := range winners {
		for _, w := range wins {
			p, _ := FromString(w)

			if r := p.IsWin(); s != r {
				t.Errorf("position: '%s' expected '%s' for position but got '%s'", w, s, r)
			}
		}
	}
}

func TestNext(t *testing.T) {
	expected := map[string]map[string]bool{
		"   \n   \n   \n": map[string]bool{
			"X  \n   \n   \n": true,
			" X \n   \n   \n": true,
			"  X\n   \n   \n": true,
			"   \nX  \n   \n": true,
			"   \n X \n   \n": true,
			"   \n  X\n   \n": true,
			"   \n   \nX  \n": true,
			"   \n   \n X \n": true,
			"   \n   \n  X\n": true,
		},
	}

	for a, next := range expected {
		p, _ := FromString(a)

		n := p.Next()
		if len(n) != len(next) {
			t.Errorf("different lengths for next positions of: '%s'", a)
			continue
		}

		for _, w := range n {
			if !next[w.String()] {
				t.Errorf("next position not found: '%s'", w)
			} else {
				delete(next, w.String())
			}
		}
	}
}

func TestTurn(t *testing.T) {
	expected := map[string]Square{
		"   \n   \n   \n": X,
		"X  \n   \n   \n": O,
		"XO \n   \n   \n": X,
	}

	for p, s := range expected {
		pp, _ := FromString(p)

		if s != pp.Turn() {
			t.Errorf("\n%s expected '%s' got '%s'", p, s, pp.Turn())
		}
	}
}
