package hyperloglog

import (
	"bytes"
	"encoding/gob"
)

type SketchOps struct {
	P          uint8
	B          uint8
	M          uint32
	Alpha      float64
	TmpSet     map[uint32]struct{}
	SparseList *struct {
		Count uint32
		Last  uint32
		B     []uint8
	}
	Regs *struct {
		T  []reg
		NZ uint32
	}
}

func (sk *Sketch) Serialize() ([]byte, error) {

	opts := SketchOps{
		P:      sk.p,
		B:      sk.b,
		M:      sk.m,
		Alpha:  sk.alpha,
		TmpSet: sk.tmpSet,
	}

	if sk.regs != nil {
		opts.Regs = &struct {
			T  []reg
			NZ uint32
		}{
			T:  sk.regs.tailcuts,
			NZ: sk.regs.nz,
		}
	}

	if sk.sparseList != nil {
		opts.SparseList = &struct {
			Count uint32
			Last  uint32
			B     []uint8
		}{
			Count: sk.sparseList.count,
			Last:  sk.sparseList.last,
			B:     sk.sparseList.b,
		}
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(opts)

	return buf.Bytes(), err
}

func DeSerialize(data []byte) (*Sketch, error) {
	var ops SketchOps
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&ops)

	if err != nil {
		return nil, err
	}

	sk := &Sketch{
		p:      ops.P,
		b:      ops.B,
		m:      ops.M,
		alpha:  ops.Alpha,
		tmpSet: ops.TmpSet,
	}

	if ops.SparseList != nil {
		sk.sparseList = &compressedList{
			count: ops.SparseList.Count,
			last:  ops.SparseList.Last,
			b:     ops.SparseList.B,
		}
	}

	if ops.Regs != nil {
		sk.regs = &registers{
			tailcuts: ops.Regs.T,
			nz:       ops.Regs.NZ,
		}
	}

	return sk, nil
}
