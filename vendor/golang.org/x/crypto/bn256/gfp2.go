// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bn256

// For details of the algorithms used, see "Multiplication and Squaring on
// Pairing-Friendly Fields, Devegili et al.
// http://eprint.iacr.org/2006/471.pdf.

import (
	"math/big"
)

// gfP2 implements a field of size p² as a quadratic extension of the base
// field where i²=-1.
type gfP2 struct ***REMOVED***
	x, y *big.Int // value is xi+y.
***REMOVED***

func newGFp2(pool *bnPool) *gfP2 ***REMOVED***
	return &gfP2***REMOVED***pool.Get(), pool.Get()***REMOVED***
***REMOVED***

func (e *gfP2) String() string ***REMOVED***
	x := new(big.Int).Mod(e.x, p)
	y := new(big.Int).Mod(e.y, p)
	return "(" + x.String() + "," + y.String() + ")"
***REMOVED***

func (e *gfP2) Put(pool *bnPool) ***REMOVED***
	pool.Put(e.x)
	pool.Put(e.y)
***REMOVED***

func (e *gfP2) Set(a *gfP2) *gfP2 ***REMOVED***
	e.x.Set(a.x)
	e.y.Set(a.y)
	return e
***REMOVED***

func (e *gfP2) SetZero() *gfP2 ***REMOVED***
	e.x.SetInt64(0)
	e.y.SetInt64(0)
	return e
***REMOVED***

func (e *gfP2) SetOne() *gfP2 ***REMOVED***
	e.x.SetInt64(0)
	e.y.SetInt64(1)
	return e
***REMOVED***

func (e *gfP2) Minimal() ***REMOVED***
	if e.x.Sign() < 0 || e.x.Cmp(p) >= 0 ***REMOVED***
		e.x.Mod(e.x, p)
	***REMOVED***
	if e.y.Sign() < 0 || e.y.Cmp(p) >= 0 ***REMOVED***
		e.y.Mod(e.y, p)
	***REMOVED***
***REMOVED***

func (e *gfP2) IsZero() bool ***REMOVED***
	return e.x.Sign() == 0 && e.y.Sign() == 0
***REMOVED***

func (e *gfP2) IsOne() bool ***REMOVED***
	if e.x.Sign() != 0 ***REMOVED***
		return false
	***REMOVED***
	words := e.y.Bits()
	return len(words) == 1 && words[0] == 1
***REMOVED***

func (e *gfP2) Conjugate(a *gfP2) *gfP2 ***REMOVED***
	e.y.Set(a.y)
	e.x.Neg(a.x)
	return e
***REMOVED***

func (e *gfP2) Negative(a *gfP2) *gfP2 ***REMOVED***
	e.x.Neg(a.x)
	e.y.Neg(a.y)
	return e
***REMOVED***

func (e *gfP2) Add(a, b *gfP2) *gfP2 ***REMOVED***
	e.x.Add(a.x, b.x)
	e.y.Add(a.y, b.y)
	return e
***REMOVED***

func (e *gfP2) Sub(a, b *gfP2) *gfP2 ***REMOVED***
	e.x.Sub(a.x, b.x)
	e.y.Sub(a.y, b.y)
	return e
***REMOVED***

func (e *gfP2) Double(a *gfP2) *gfP2 ***REMOVED***
	e.x.Lsh(a.x, 1)
	e.y.Lsh(a.y, 1)
	return e
***REMOVED***

func (c *gfP2) Exp(a *gfP2, power *big.Int, pool *bnPool) *gfP2 ***REMOVED***
	sum := newGFp2(pool)
	sum.SetOne()
	t := newGFp2(pool)

	for i := power.BitLen() - 1; i >= 0; i-- ***REMOVED***
		t.Square(sum, pool)
		if power.Bit(i) != 0 ***REMOVED***
			sum.Mul(t, a, pool)
		***REMOVED*** else ***REMOVED***
			sum.Set(t)
		***REMOVED***
	***REMOVED***

	c.Set(sum)

	sum.Put(pool)
	t.Put(pool)

	return c
***REMOVED***

// See "Multiplication and Squaring in Pairing-Friendly Fields",
// http://eprint.iacr.org/2006/471.pdf
func (e *gfP2) Mul(a, b *gfP2, pool *bnPool) *gfP2 ***REMOVED***
	tx := pool.Get().Mul(a.x, b.y)
	t := pool.Get().Mul(b.x, a.y)
	tx.Add(tx, t)
	tx.Mod(tx, p)

	ty := pool.Get().Mul(a.y, b.y)
	t.Mul(a.x, b.x)
	ty.Sub(ty, t)
	e.y.Mod(ty, p)
	e.x.Set(tx)

	pool.Put(tx)
	pool.Put(ty)
	pool.Put(t)

	return e
***REMOVED***

func (e *gfP2) MulScalar(a *gfP2, b *big.Int) *gfP2 ***REMOVED***
	e.x.Mul(a.x, b)
	e.y.Mul(a.y, b)
	return e
***REMOVED***

// MulXi sets e=ξa where ξ=i+3 and then returns e.
func (e *gfP2) MulXi(a *gfP2, pool *bnPool) *gfP2 ***REMOVED***
	// (xi+y)(i+3) = (3x+y)i+(3y-x)
	tx := pool.Get().Lsh(a.x, 1)
	tx.Add(tx, a.x)
	tx.Add(tx, a.y)

	ty := pool.Get().Lsh(a.y, 1)
	ty.Add(ty, a.y)
	ty.Sub(ty, a.x)

	e.x.Set(tx)
	e.y.Set(ty)

	pool.Put(tx)
	pool.Put(ty)

	return e
***REMOVED***

func (e *gfP2) Square(a *gfP2, pool *bnPool) *gfP2 ***REMOVED***
	// Complex squaring algorithm:
	// (xi+b)² = (x+y)(y-x) + 2*i*x*y
	t1 := pool.Get().Sub(a.y, a.x)
	t2 := pool.Get().Add(a.x, a.y)
	ty := pool.Get().Mul(t1, t2)
	ty.Mod(ty, p)

	t1.Mul(a.x, a.y)
	t1.Lsh(t1, 1)

	e.x.Mod(t1, p)
	e.y.Set(ty)

	pool.Put(t1)
	pool.Put(t2)
	pool.Put(ty)

	return e
***REMOVED***

func (e *gfP2) Invert(a *gfP2, pool *bnPool) *gfP2 ***REMOVED***
	// See "Implementing cryptographic pairings", M. Scott, section 3.2.
	// ftp://136.206.11.249/pub/crypto/pairings.pdf
	t := pool.Get()
	t.Mul(a.y, a.y)
	t2 := pool.Get()
	t2.Mul(a.x, a.x)
	t.Add(t, t2)

	inv := pool.Get()
	inv.ModInverse(t, p)

	e.x.Neg(a.x)
	e.x.Mul(e.x, inv)
	e.x.Mod(e.x, p)

	e.y.Mul(a.y, inv)
	e.y.Mod(e.y, p)

	pool.Put(t)
	pool.Put(t2)
	pool.Put(inv)

	return e
***REMOVED***
