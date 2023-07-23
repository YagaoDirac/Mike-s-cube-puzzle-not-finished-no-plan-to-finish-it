//A super small bitset tool. Author: github.com/yagaodirac. Twitter: yagaodirac
//You mean what is the licence? What licence? Driving licence? Why do I need a driving licence to write code?

package bitset

import "fmt"

type Bitset struct {
	Len  int
	Data []uint64
}

func (bs Bitset) FindAllIndex() (out []int) {
	for i, v := range bs.Data {
		for inner_i := 0; inner_i < 64; inner_i++ {
			if v&(1<<inner_i) != 0 {
				out = append(out, i*64+inner_i)
			}
		}
	}
	return out
}

func (bs Bitset) FindFirstN(n int) (out []int) {
	if n < 0 {
		panic("n should be at least 0")
	}
	if 0 == n {
		return
	}
	for i, v := range bs.Data {
		for inner_i := 0; inner_i < 64; inner_i++ {
			if v&(1<<inner_i) != 0 {
				out = append(out, i*64+inner_i)
				if len(out) >= n {
					return out
				}
			}
		}
	}
	return out
}

func (bs *Bitset) SetLen(newLength int) {
	if newLength < 0 {
		panic("newLength must be at least 0")
	}
	var oldDataLen = len(bs.Data)
	var newDataLen = (newLength + 63) / 64
	if newDataLen <= oldDataLen {
		if (newDataLen > 26) && (oldDataLen/2 > newDataLen) { //I didn't test the condition, change it if you know what you are doing. Or if you don't care the extreme performance, don't touch it.
			var temp = make([]uint64, newDataLen)
			copy(temp, bs.Data)
			bs.Data = temp
		} else {
			//Otherwise no need to shrink, only need to set to 0s.
			for i := newDataLen; i < oldDataLen; i++ {
				bs.Data[i] = 0
			}
		}
		if newDataLen > 0 {

			//var boundary = (newLength - 1 + 64) & 0b111111
			var boundary = (uint(newLength) - 1) % 64
			boundary = 63 - boundary
			bs.Data[newDataLen-1] <<= boundary
			bs.Data[newDataLen-1] >>= boundary
		}

	} else { //if newDataLen>oldDataLen{
		bs.Data = append(bs.Data, make([]uint64, newDataLen-oldDataLen)...)
	}
	bs.Len = newLength
}
func (bs *Bitset) SetBit(pos int) {
	if pos < 0 || pos >= bs.Len {
		panic("pos must be >= 0 and < len of this Bitset")
	}
	var posOver64 = pos / 64
	bs.Data[posOver64] |= 1 << (pos - posOver64*64)
}
func (bs *Bitset) ResetBit(pos int) {
	if pos < 0 {
		panic("newLength must be at least 0")
	}
	var posOver64 = pos / 64
	bs.Data[posOver64] &= ^(1 << (pos - posOver64*64))
}
func (bs *Bitset) Print_() {
	fmt.Print("Len: ", bs.Len, ", bits at: ")
	var in = bs.FindAllIndex()
	for _, v := range in {
		fmt.Print(v, ", ")
	}
	fmt.Println()
}
