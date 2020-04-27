package hyperloglog

import "encoding/json"

type SketchOps struct {
	P          uint8
	B          uint8
	M          uint32
	Alpha      float64
	TmpSet     set
	SparseList struct {
		Count uint32
		Last  uint32
		B     []uint8
	}
	Regs struct {
		t  []reg
		nz uint32
	}
}

func (sk *Sketch) Serialize() ([]byte, error) {

	opts := SketchOps{
		P:      sk.p,
		B:      sk.b,
		M:      sk.m,
		Alpha:  sk.alpha,
		TmpSet: sk.tmpSet,
		SparseList: struct {
			Count uint32
			Last  uint32
			B     []uint8
		}{
			Count: sk.sparseList.count,
			Last:  sk.sparseList.last,
			B:     sk.sparseList.b,
		},
	}

	if sk.regs != nil {
		opts.Regs = struct {
			t  []reg
			nz uint32
		}{
			t:  sk.regs.tailcuts,
			nz: sk.regs.nz,
		}
	}

	if sk.sparseList != nil {
		opts.SparseList = struct {
			Count uint32
			Last  uint32
			B     []uint8
		}{
			Count: sk.sparseList.count,
			Last:  sk.sparseList.last,
			B:     sk.sparseList.b,
		}
	}

	return json.Marshal(opts)
}

func DeSerialize(data []byte) (*Sketch, error) {
	var ops SketchOps
	err := json.Unmarshal(data, &ops)

	if err != nil {
		return nil, err
	}

	return &Sketch{
		p:      ops.P,
		b:      ops.B,
		m:      ops.M,
		alpha:  ops.Alpha,
		tmpSet: ops.TmpSet,
		sparseList: &compressedList{
			count: ops.SparseList.Count,
			last:  ops.SparseList.Last,
			b:     ops.SparseList.B,
		},
		regs: &registers{
			tailcuts: ops.Regs.t,
			nz:       ops.Regs.nz,
		},
	}, nil

}
