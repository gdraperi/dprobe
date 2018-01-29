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

// gfP12 implements the field of size p¹² as a quadratic extension of gfP6
// where ω²=τ.
type gfP12 struct ***REMOVED***
	x, y *gfP6 // value is xω + y
***REMOVED***

func newGFp12(pool *bnPool) *gfP12 ***REMOVED***
	return &gfP12***REMOVED***newGFp6(pool), newGFp6(pool)***REMOVED***
***REMOVED***

func (e *gfP12) String() string ***REMOVED***
	return "(" + e.x.String() + "," + e.y.String() + ")"
***REMOVED***

func (e *gfP12) Put(pool *bnPool) ***REMOVED***
	e.x.Put(pool)
	e.y.Put(pool)
***REMOVED***

func (e *gfP12) Set(a *gfP12) *gfP12 ***REMOVED***
	e.x.Set(a.x)
	e.y.Set(a.y)
	return e
***REMOVED***

func (e *gfP12) SetZero() *gfP12 ***REMOVED***
	e.x.SetZero()
	e.y.SetZero()
	return e
***REMOVED***

func (e *gfP12) SetOne() *gfP12 ***REMOVED***
	e.x.SetZero()
	e.y.SetOne()
	return e
***REMOVED***

func (e *gfP12) Minimal() ***REMOVED***
	e.x.Minimal()
	e.y.Minimal()
***REMOVED***

func (e *gfP12) IsZero() bool ***REMOVED***
	e.Minimal()
	return e.x.IsZero() && e.y.IsZero()
***REMOVED***

func (e *gfP12) IsOne() bool ***REMOVED***
	e.Minimal()
	return e.x.IsZero() && e.y.IsOne()
***REMOVED***

func (e *gfP12) Conjugate(a *gfP12) *gfP12 ***REMOVED***
	e.x.Negative(a.x)
	e.y.Set(a.y)
	return a
***REMOVED***

func (e *gfP12) Negative(a *gfP12) *gfP12 ***REMOVED***
	e.x.Negative(a.x)
	e.y.Negative(a.y)
	return e
***REMOVED***

// Frobenius computes (xω+y)^p = x^p ω·ξ^((p-1)/6) + y^p
func (e *gfP12) Frobenius(a *gfP12, pool *bnPool) *gfP12 ***REMOVED***
	e.x.Frobenius(a.x, pool)
	e.y.Frobenius(a.y, pool)
	e.x.MulScalar(e.x, xiToPMinus1Over6, pool)
	return e
***REMOVED***

// FrobeniusP2 computes (xω+y)^p² = x^p² ω·ξ^((p²-1)/6) + y^p²
func (e *gfP12) FrobeniusP2(a *gfP12, pool *bnPool) *gfP12 ***REMOVED***
	e.x.FrobeniusP2(a.x)
	e.x.MulGFP(e.x, xiToPSquaredMinus1Over6)
	e.y.FrobeniusP2(a.y)
	return e
***REMOVED***

func (e *gfP12) Add(a, b *gfP12) *gfP12 ***REMOVED***
	e.x.Add(a.x, b.x)
	e.y.Add(a.y, b.y)
	return e
***REMOVED***

func (e *gfP12) Sub(a, b *gfP12) *gfP12 ***REMOVED***
	e.x.Sub(a.x, b.x)
	e.y.Sub(a.y, b.y)
	return e
***REMOVED***

func (e *gfP12) Mul(a, b *gfP12, pool *bnPool) *gfP12 ***REMOVED***
	tx := newGFp6(pool)
	tx.Mul(a.x, b.y, pool)
	t := newGFp6(pool)
	t.Mul(b.x, a.y, pool)
	tx.Add(tx, t)

	ty := newGFp6(pool)
	ty.Mul(a.y, b.y, pool)
	t.Mul(a.x, b.x, pool)
	t.MulTau(t, pool)
	e.y.Add(ty, t)
	e.x.Set(tx)

	tx.Put(pool)
	ty.Put(pool)
	t.Put(pool)
	return e
***REMOVED***

func (e *gfP12) MulScalar(a *gfP12, b *gfP6, pool *bnPool) *gfP12 ***REMOVED***
	e.x.Mul(e.x, b, pool)
	e.y.Mul(e.y, b, pool)
	return e
***REMOVED***

func (c *gfP12) Exp(a *gfP12, power *big.Int, pool *bnPool) *gfP12 ***REMOVED***
	sum := newGFp12(pool)
	sum.SetOne()
	t := newGFp12(pool)

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

func (e *gfP12) Square(a *gfP12, pool *bnPool) *gfP12 ***REMOVED***
	// Complex squaring algorithm
	v0 := newGFp6(pool)
	v0.Mul(a.x, a.y, pool)

	t := newGFp6(pool)
	t.MulTau(a.x, pool)
	t.Add(a.y, t)
	ty := newGFp6(pool)
	ty.Add(a.x, a.y)
	ty.Mul(ty, t, pool)
	ty.Sub(ty, v0)
	t.MulTau(v0, pool)
	ty.Sub(ty, t)

	e.y.Set(ty)
	e.x.Double(v0)

	v0.Put(pool)
	t.Put(pool)
	ty.Put(pool)

	return e
***REMOVED***

func (e *gfP12) Invert(a *gfP12, pool *bnPool) *gfP12 ***REMOVED***
	// See "Implementing cryptographic pairings", M. Scott, section 3.2.
	// ftp://136.206.11.249/pub/crypto/pairings.pdf
	t1 := newGFp6(pool)
	t2 := newGFp6(pool)

	t1.Square(a.x, pool)
	t2.Square(a.y, pool)
	t1.MulTau(t1, pool)
	t1.Sub(t2, t1)
	t2.Invert(t1, pool)

	e.x.Negative(a.x)
	e.y.Set(a.y)
	e.MulScalar(e, t2, pool)

	t1.Put(pool)
	t2.Put(pool)

	return e
***REMOVED***
