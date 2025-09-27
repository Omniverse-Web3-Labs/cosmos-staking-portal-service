package utils

import "math/big"

type BigFloat struct {
	v *big.Float
}

func NewBigFloat() *BigFloat {
	return &BigFloat{v: big.NewFloat(0)}
}

func NewBigFloatFromFloat64(f float64) *BigFloat {
	return &BigFloat{v: big.NewFloat(f)}
}

func NewBigFloatFromBigInt(i *big.Int) *BigFloat {
	return &BigFloat{v: big.NewFloat(0).SetInt(i)}
}

func NewBigFloatFromBigFloat(f *big.Float) *BigFloat {
	return &BigFloat{v: big.NewFloat(0).Set(f)}
}

func Complement(f *big.Float) *big.Float {
	return big.NewFloat(0).Sub(big.NewFloat(1), f)
}

func (f *BigFloat) Mul(x *big.Float) *BigFloat {
	f.v.Mul(f.v, x)
	return f
}

func (f *BigFloat) Div(x *big.Float) *BigFloat {
	f.v.Quo(f.v, x)
	return f
}

func (f *BigFloat) Add(x *big.Float) *BigFloat {
	f.v.Add(f.v, x)
	return f
}

func (f *BigFloat) Sub(x *big.Float) *BigFloat {
	f.v.Sub(f.v, x)
	return f
}

func (f *BigFloat) Val() *big.Float {
	return f.v
}
