package hyperloglog

import (
	"fmt"
	"testing"
)

func TestSketch_Serialization(t *testing.T) {

	var loadedHll *Sketch

	manipulations := []func(sk *Sketch){
		func(sk *Sketch) {
			// Do nothing
		},
		func(sk *Sketch) {
			sk.Insert([]byte("somethung"))
		},
		func(sk *Sketch) {
			sk.Insert([]byte("somethung"))
		},
		func(sk *Sketch) {
			for i := 0; i < 100; i++ {
				sk.Insert([]byte(fmt.Sprintf("Item %d", i)))
			}
		},
	}

	subjects := []struct {
		name   string
		sketch *Sketch
	}{
		{
			name:   "sparse_14",
			sketch: New(),
		},
		{
			name:   "sparse_16",
			sketch: New16(),
		},
		{
			name:   "no_sparse_14",
			sketch: NewNoSparse(),
		},
		{
			name:   "no_sparse_16",
			sketch: New16NoSparse(),
		},
	}

	for _, subject := range subjects {
		sk := subject.sketch
		for i, op := range manipulations {

			t.Run(fmt.Sprintf("Serialization_%d_for_%s", i, subject.name), func(t *testing.T) {
				op(sk)

				data, err := sk.Serialize()

				if err != nil {
					t.Error(err)
				}

				loadedHll, err = DeSerialize(data)

				if err != nil {
					t.Error(err)
				}

				if sk.Estimate() != loadedHll.Estimate() {
					t.Errorf("Estimations not equal; expected %d got %d", sk.Estimate(), loadedHll.Estimate())
				}
			})
		}
	}

}
