// Package kerl implements the Kerl hashing function.
package kerl

import (
	"unsafe"

	"github.com/pkg/errors"

	. "github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/kerl/bigint"
	. "github.com/iotaledger/iota.go/trinary"
)

const (
	// radix used in the conversion
	tryteRadix = 27
	// the middle of the domain described by one tryte
	halfTryte = 1 + 3 + 9
)

// hex representation of the middle of the domain described by 242 trits, i.e. \sum_{k=0}^{241} 3^k
var halfThree = []uint32{
	0xa5ce8964, 0x9f007669, 0x1484504f, 0x3ade00d9, 0x0c24486e, 0x50979d57,
	0x79a4c702, 0x48bbae36, 0xa9f6808b, 0xaa06a805, 0xa87fabdf, 0x5e69ebef,
}

// hex representation of the two's complement of halfThree, i.e. ~halfThree + 1
var negHalfThree = []uint32{
	0x5a31769c, 0x60ff8996, 0xeb7bafb0, 0xc521ff26, 0xf3dbb791, 0xaf6862a8,
	0x865b38fd, 0xb74451c9, 0x56097f74, 0x55f957fa, 0x57805420, 0xa1961410,
}

// hex representation of the last trit, i.e. 3^242
var trit243 = []uint32{
	0x4b9d12c9, 0x3e00ecd3, 0x2908a09f, 0x75bc01b2, 0x184890dc, 0xa12f3aae,
	0xf3498e04, 0x91775c6c, 0x53ed0116, 0x540d500b, 0x50ff57bf, 0xbcd3d7df,
}

// lookup table to convert tryte values into trits
var tryteValueToTritsLUT = [][3]int8{
	{-1, -1, -1}, {0, -1, -1}, {1, -1, -1}, {-1, 0, -1}, {0, 0, -1}, {1, 0, -1},
	{-1, 1, -1}, {0, 1, -1}, {1, 1, -1}, {-1, -1, 0}, {0, -1, 0}, {1, -1, 0},
	{-1, 0, 0}, {0, 0, 0}, {1, 0, 0}, {-1, 1, 0}, {0, 1, 0}, {1, 1, 0},
	{-1, -1, 1}, {0, -1, 1}, {1, -1, 1}, {-1, 0, 1}, {0, 0, 1}, {1, 0, 1},
	{-1, 1, 1}, {0, 1, 1}, {1, 1, 1},
}

// lookup table to convert tryte values into trytes
const tryteValueToTyteLUT = "NOPQRSTUVWXYZ9ABCDEFGHIJKLM"

func tryteValuesToTrytes(vs []int8) Trytes {
	trytes := make([]byte, len(vs))
	for i, v := range vs {
		idx := v - MinTryteValue
		trytes[i] = tryteValueToTyteLUT[idx]
	}
	// convert to string without copying
	return *(*string)(unsafe.Pointer(&trytes))
}

func tryteValuesToTrits(vs []int8) Trits {
	trits := make([]int8, len(vs)*3)
	for i, v := range vs {
		idx := v - MinTryteValue
		trits[i*3], trits[i*3+1], trits[i*3+2] = tryteValueToTritsLUT[idx][0], tryteValueToTritsLUT[idx][1], tryteValueToTritsLUT[idx][2]
	}
	return trits
}

func trytesToTryteValues(trytes Trytes) []int8 {
	vs := make([]int8, len(trytes))
	for i, tryte := range trytes {
		switch {
		case tryte == '9':
			vs[i] = 0
		case tryte >= 'N':
			vs[i] = int8(tryte) - 'N' + MinTryteValue
		default:
			vs[i] = int8(tryte) - 'A' + 1
		}
	}
	return vs
}

func tritsToTryteValues(trits Trits) []int8 {
	vs := make([]int8, len(trits)/3)
	for i := 0; i < len(trits)/3; i++ {
		vs[i] = trits[i*3] + trits[i*3+1]*3 + trits[i*3+2]*9
	}
	return vs
}

// tryteZeroLastTrit takes a tryte value of three trits a+3b+9c and returns a+3b (setting the last trit to zero).
func tryteZeroLastTrit(v int8) int8 {
	if v > 4 {
		return v - 9
	}
	if v < -4 {
		return v + 9
	}
	return v
}

// bigintZeroLastTrit changes the bigint so that the corresponding ternary number has 242th trit set to 0.
func bigintZeroLastTrit(b []uint32) {
	if bigint.IsNegative(b) {
		if bigint.MustCmp(b, negHalfThree) < 0 {
			bigint.MustAdd(b, trit243)
		}
	} else {
		if bigint.MustCmp(b, halfThree) > 0 {
			bigint.MustSub(b, trit243)
		}
	}
}

func tryteValuesToBytes(vs []int8) []byte {
	bytes := make([]byte, HashBytesSize)
	b := (*(*[]uint32)(unsafe.Pointer(&bytes)))[0:IntLength]

	// set the last trit of the last tryte to zero
	v := tryteZeroLastTrit(vs[HashTrytesSize-1])
	// initialize the first part of the bigint with the non-balanced representation of this 2-trit value
	b[0] = uint32(v + 4)

	// initially, all words of the bigint are zero
	nzIndex := 0
	for i := HashTrytesSize - 2; i >= 0; i-- {
		// first, multiply the bigint by the radix
		var carry uint32
		for i := 0; i <= nzIndex; i++ {
			v := tryteRadix*uint64(b[i]) + uint64(carry)
			carry, b[i] = uint32(v>>32), uint32(v)
		}
		if carry > 0 && nzIndex < IntLength-1 {
			nzIndex++
			b[nzIndex] = carry
		}

		// then, add the non-balanced tryte value
		chgIndex := bigint.AddSmall(b, uint32(vs[i]+halfTryte))
		// adapt the non-zero index, if we had an overflow
		if chgIndex > nzIndex {
			nzIndex = chgIndex
		}
	}

	// subtract the middle of the domain to get balanced ternary
	bigint.MustSub(b, halfThree)

	// convert to bytes
	return bigint.Reverse(bytes)
}

func bytesToTryteValues(bytes []byte) []int8 {
	// copy and convert bytes to bigint
	rb := make([]byte, len(bytes))
	copy(rb, bytes)
	bigint.Reverse(rb)
	b := (*(*[]uint32)(unsafe.Pointer(&rb)))[0:IntLength]

	// the two's complement representation is only correct, if the number fits
	// into 48 bytes, i.e. has the 243th trit set to 0
	bigintZeroLastTrit(b)

	// convert to the unsigned bigint representing non-balanced ternary
	bigint.MustAdd(b, halfThree)

	vs := make([]int8, HashTrytesSize)

	// initially, all words of the bigint are non-zero
	nzIndex := IntLength - 1
	for i := 0; i < HashTrytesSize-1; i++ {
		// divide the bigint by the radix
		var rem uint32
		for i := nzIndex; i >= 0; i-- {
			v := (uint64(rem) << 32) | uint64(b[i])
			b[i], rem = uint32(v/tryteRadix), uint32(v%tryteRadix)
		}
		// the tryte value is the remainder converted back to balanced ternary
		vs[i] = int8(rem) - halfTryte

		// decrement index, if the highest considered word of the bigint turned zero
		if nzIndex > 0 && b[nzIndex] == 0 {
			nzIndex--
		}
	}

	// special case for the last tryte, where no further division is necessary
	vs[HashTrytesSize-1] = tryteZeroLastTrit(int8(b[0]) - halfTryte)

	return vs
}

// KerlBytesZeroLastTrit changes a chunk of 48 bytes so that the corresponding ternary number has 242th trit set to 0.
func KerlBytesZeroLastTrit(bytes []byte) {
	// convert to bigint
	bigint.Reverse(bytes)
	b := (*(*[]uint32)(unsafe.Pointer(&bytes)))[0:IntLength]

	bigintZeroLastTrit(b)
	bigint.Reverse(bytes)
}

// KerlTritsToBytes is only defined for hashes, i.e. chunks of trits of length 243. It returns 48 bytes.
func KerlTritsToBytes(trits Trits) ([]byte, error) {
	if !CanBeHash(trits) {
		return nil, errors.Wrapf(ErrInvalidTritsLength, "must be %d in size", HashTrinarySize)
	}

	vs := tritsToTryteValues(trits)
	return tryteValuesToBytes(vs), nil
}

// KerlTrytesToBytes is only defined for hashes, i.e. chunks of trytes of length 81. It returns 48 bytes.
func KerlTrytesToBytes(trytes Trytes) ([]byte, error) {
	if len(trytes) != HashTrytesSize {
		return nil, errors.Wrapf(ErrInvalidTrytesLength, "must be %d in size", HashBytesSize)
	}

	vs := trytesToTryteValues(trytes)
	return tryteValuesToBytes(vs), nil
}

// KerlBytesToTrits is only defined for hashes, i.e. chunks of 48 bytes. It returns 243 trits.
func KerlBytesToTrits(b []byte) (Trits, error) {
	if len(b) != HashBytesSize {
		return nil, errors.Wrapf(ErrInvalidBytesLength, "must be %d in size", HashBytesSize)
	}

	vs := bytesToTryteValues(b)
	return tryteValuesToTrits(vs), nil
}

// KerlBytesToTrytes is only defined for hashes, i.e. chunks of 48 bytes. It returns 81 trytes.
func KerlBytesToTrytes(b []byte) (Trytes, error) {
	if len(b) != HashBytesSize {
		return "", errors.Wrapf(ErrInvalidBytesLength, "must be %d in size", HashBytesSize)
	}

	vs := bytesToTryteValues(b)
	return tryteValuesToTrytes(vs), nil
}
