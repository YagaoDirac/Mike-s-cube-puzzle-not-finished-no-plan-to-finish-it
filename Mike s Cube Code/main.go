//The question is from this video: https://www.youtube.com/watch?v=g9n0a0644B4
//From block count == 6, it's not guaranteed the first bit of data of a shape is 1, after cleaning.
//The structure of the look up table is like, slice for block count, map for sizeX, map for  sizeY, map for sizeZ, then a slice. Before I use it for searching, I have to sort it.

/*
this proj is not finished. I didn't consider the rotation. If shapes keep their direction, it's easier to do the searching and comparing. But the final result needs to eliminate all the rotating-repeating ones. The report function is not finished, but before that, it should be correct enough.

Basic idea is like, all the shapes are moved to as near to the original point as possible, and try mirroring them. Mirroring a shape doesn't change its dimentional size, so I only have to compare shapes with the same dimentional size. It's possible to do the comparing only with a binary representation. This way keeps different shapes which can be rotated into the same one. Then I try generating all the possible shapes by add the extra block to existing shapes, move to origin, mirror and figure out which is the best, sort the best ones and remove duplicated.
But if I what to figure out the actual count, I have to rotate all the shapes, and mirror them, and then move to origin, sort, remove duplicated. This is beyond my plan, I only want some practice in go, since this is the first proj I've tried in golang. Now, the math part is boring and costs time for nothing. If you want to use this you have to fix the report func.
*/

package main

import (
	"bitset"
	"fmt"
	"sort"
	//"golang.org/x/tools/go/analysis/passes/sortslice"
	//"github.com/petar/GoLLRB"
)

const (
	MaxStep       = 3 //at least when it works for max step == 6, it's possible to be fully correct.
	SAVE_MEM_MODE = false
)

type vec3int struct {
	X, Y, Z int
}

// func (in vec3int) mirror(x, y, z bool, dimSize vec3int) (out vec3int) { //2 b tested
// 	dimSize.X--
// 	dimSize.Y--
// 	dimSize.Z--
// 	if x {
// 		in.X = dimSize.X - in.X
// 	}
// 	if y {
// 		in.Y = dimSize.Y - in.Y
// 	}
// 	if z {
// 		in.Z = dimSize.Z - in.Z
// 	}
// 	return out
// }
// func (in vec3int) mirror_sameDimSize(x, y, z bool, dimSize int) (out vec3int) { //2 b tested
// 	dimSize--
// 	if x {
// 		in.X = dimSize - in.X
// 	}
// 	if y {
// 		in.Y = dimSize - in.Y
// 	}
// 	if z {
// 		in.Z = dimSize - in.Z
// 	}
// 	return out
// }

type shapeData struct {
	Dim  vec3int //the size this shape takes.
	Data bitset.Bitset
}
type rawShape shapeData
type rawShape2 struct {
	blocks           []vec3int
	sorted           bool
	dimSizeShrinkBy1 vec3int
}
type vec3intSortZYX []vec3int

func (a vec3intSortZYX) Len() int      { return len(a) }
func (a vec3intSortZYX) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a vec3intSortZYX) Less(i, j int) bool { ///////////
	if a[i].Z < a[j].Z {
		return true
	}
	if a[i].Z > a[j].Z {
		return false
	}
	if a[i].Y < a[j].Y {
		return true
	}
	if a[i].Y > a[j].Y {
		return false
	}
	return a[i].X < a[j].X
}

type shape shapeData

type shapeSortSafe []shape

func (a shapeSortSafe) Len() int      { return len(a) }
func (a shapeSortSafe) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a shapeSortSafe) Less(i, j int) bool {
	if a[i].Dim != a[j].Dim {
		panic("inside func (a shapeSortSafe) Less(i, j int) bool")
	}
	if len(a[i].Data.Data) != len(a[j].Data.Data) {
		panic("inside func (a shapeSortSafe) Less(i, j int) bool")
	}
	for index, fromI := range a[i].Data.Data {
		var fromJ = a[j].Data.Data[index]
		if fromI < fromJ {
			return true
		}
		if fromI > fromJ {
			return false
		}
	}
	return false
}

type shapeSort []shape            //a bit faster
func (a shapeSort) Len() int      { return len(a) }
func (a shapeSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a shapeSort) Less(i, j int) bool {
	for index, fromI := range a[i].Data.Data {
		var fromJ = a[j].Data.Data[index]
		if fromI < fromJ {
			return true
		}
		if fromI > fromJ {
			return false
		}
	}
	return false
}

func (in vec3int) findNeighbour(size_x, size_y, size_z int) (out []vec3int) { //2 b tested
	var all_x = make([]int, 0, 2)
	if in.X > 0 {
		all_x = append(all_x, in.X-1)
	}
	if in.X < size_x-1 {
		all_x = append(all_x, in.X+1)
	}
	var all_y = make([]int, 0, 2)
	if in.Y > 0 {
		all_y = append(all_y, in.Y-1)
	}
	if in.Y < size_y-1 {
		all_y = append(all_y, in.Y+1)
	}
	var all_z = make([]int, 0, 2)
	if in.Z > 0 {
		all_z = append(all_z, in.Z-1)
	}
	if in.Z < size_z-1 {
		all_z = append(all_z, in.Z+1)
	}
	for _, x := range all_x {
		for _, y := range all_y {
			for _, z := range all_z {
				out = append(out, vec3int{x, y, z})
			}
		}
	}
	return
}

// func (in vec3int) getKey() (out int32) { //2 b tested
// 	out = int32(in.X)
// 	out ^= int32(in.Y << 10)
// 	out ^= int32(in.Z << 20)
// 	return
// }

// func (in *shapeData) mirror(x, y, z bool, blockCount int) (out rawShape) { //2 b tested
//
//		var already_found = 0
//		out.Data = make([]int64, len(in.Data))
//		out.Dim.X = in.Dim.X
//		out.Dim.Y = in.Dim.Y
//		out.Dim.Z = in.Dim.Z
//		var bits = in.Data.FindAllIndex() //Maybe replacing this with findfirstN will speed up a bit.
//		for _, index := range bits {
//			var in_x, in_y, in_z = indexToCoord(index, in)
//			if x {
//				in_x = in.Dim.X - in_x - 1
//			}
//			if y {
//				in_y = in.Dim.Y - in_y - 1
//			}
//			if z {
//				in_z = in.Dim.Z - in_z - 1
//			}
//			index = coordToIndex(in_x, in_y, in_z, in)
//			out.setBitTo1(index)
//		}
//		// for out_i, v := range in.Data {
//		// 	if v != 0 {
//		// 		for i := 0; i < 64; i++ {
//		// 			var mask = int64(1 << i)
//		// 			var masked = mask & v
//		// 			if masked != 0 {
//		// 				already_found++
//		//				var index = out_i*64 + i
//	}
func (in rawShape2) mirror(x, y, z bool) (out rawShape2) { //2 b tested  1
	out = in
	out.blocks = make([]vec3int, 0, len(in.blocks))
	for _, block := range in.blocks {
		if x {
			block.X = in.dimSizeShrinkBy1.X - block.X
		}
		if y {
			block.Y = in.dimSizeShrinkBy1.Y - block.Y
		}
		if z {
			block.Z = in.dimSizeShrinkBy1.Z - block.Z
		}
		out.blocks = append(out.blocks, block)
	}
	out.sort()
	return out
}
func (in shape) mirror(x, y, z bool) (out rawShape2) { //2 b tested
	var bits = in.Data.FindAllIndex()
	var dimSizeShrinkBy1 = in.Dim
	dimSizeShrinkBy1.X--
	dimSizeShrinkBy1.Y--
	dimSizeShrinkBy1.Z--
	for _, index := range bits {
		var block = indexToVec(index, in.Dim)
		if x {
			block.X = dimSizeShrinkBy1.X - block.X
		}
		if y {
			block.Y = dimSizeShrinkBy1.Y - block.Y
		}
		if z {
			block.Z = dimSizeShrinkBy1.Z - block.Z
		}
		out.blocks = append(out.blocks, block)
	}
	out.dimSizeShrinkBy1 = in.Dim
	out.dimSizeShrinkBy1.X--
	out.dimSizeShrinkBy1.Y--
	out.dimSizeShrinkBy1.Z--
	out.sort()
	return out
}
func (in *rawShape) alignDataToDim() {
	in.Data.SetLen(in.Dim.X * in.Dim.Y * in.Dim.Z)
	in.Data.Data = make([]uint64, len(in.Data.Data))
}
func (in *shape) alignDataToDim() {
	in.Data.SetLen(in.Dim.X * in.Dim.Y * in.Dim.Z)
	in.Data.Data = make([]uint64, len(in.Data.Data))
}

func (in shape) toRawShape2() (out rawShape2) { //2 b tested
	var indice = in.Data.FindAllIndex()
	for _, index := range indice {
		out.addBlock(indexToVec(index, in.Dim))
	}
	out.dimSizeShrinkBy1 = in.Dim
	out.dimSizeShrinkBy1.X--
	out.dimSizeShrinkBy1.Y--
	out.dimSizeShrinkBy1.Z--
	out.sort()
	return out
}
func (in *rawShape2) sort() { //2 b tested 1
	if in.sorted {
		return
	}
	sort.Sort(vec3intSortZYX(in.blocks))
	in.sorted = true
}

func (s *rawShape) setBitTo1(index int) { //2 b tested
	s.Data.SetBit(index)
}
func (s *rawShape2) addBlock(pos vec3int) {
	s.blocks = append(s.blocks, pos)
	s.sorted = false
}

func (in *shapeData) copy() (out rawShape) { //2 b tested
	out.Data.Len = in.Data.Len
	out.Data.Data = make([]uint64, len(in.Data.Data))
	copy(out.Data.Data, in.Data.Data)
	out.Dim.X = in.Dim.X
	out.Dim.Y = in.Dim.Y
	out.Dim.Z = in.Dim.Z
	return
}

// func copyShape[copiableShape shapeData | shape | rawShape](in *copiableShape) (out rawShape) { //2 b tested
//
//		out.Data.Len = in.Data.Len
//		out.Data.Data = make([]uint64, len(in.Data.Data))
//		copy(out.Data.Data, in.Data.Data)
//		out.Dim.X = in.Dim.X
//		out.Dim.Y = in.Dim.Y
//		out.Dim.Z = in.Dim.Z
//		return
//	}
func (in rawShape2) copy() (out rawShape2) { //2 b tested
	out = in
	out.blocks = make([]vec3int, len(in.blocks))
	copy(out.blocks, in.blocks)
	return out
}

func coordToIndex(x, y, z int, Dim vec3int) (index int) { //2 b tested
	index = (z*Dim.Y+y)*Dim.X + x
	return index
}
func vecToIndex(vec vec3int, Dim vec3int) (index int) { //2 b tested
	index = (vec.Z*Dim.Y+vec.Y)*Dim.X + vec.X
	return index
}

func indexToCoord(index int, Dim vec3int) (x, y, z int) { //2 b tested
	var xy = Dim.X * Dim.Y
	z = index / xy
	var temp = index - z*xy
	y = temp / Dim.X
	x = temp - Dim.X*y
	return
}
func indexToVec(index int, Dim vec3int) (vec vec3int) { //2 b tested
	var xy = Dim.X * Dim.Y
	vec.Z = index / xy
	var temp = index - vec.Z*xy
	vec.Y = temp / Dim.X
	vec.X = temp - Dim.X*vec.Y
	return vec
}

// func (in *shape) getKey() (key int64) { //2 b tested
// 	Do I really need this func?
// 	var len_of_data = len(in.Data.Data)
// 	key = in.Data.Data[0]
// 	for i := 1; i < len_of_data; i++ {
// 		key ^= in.Data.Data[i]
// 	}
// 	key ^= int64(in.x)
// 	key ^= int64(in.y)
// 	key ^= int64(in.z)
// 	return key
// }

// func (in *shape) valid_(blockCount int) (correct bool) {}
// func (in *shape) getBlockCount_() (out int) {}
func (in *shape) isSame(other shape) bool { //2 b tested
	if len(in.Data.Data) != len(other.Data.Data) {
		return false
	}
	if in.Dim.X != other.Dim.X || in.Dim.Y != other.Dim.Y || in.Dim.Z != other.Dim.Z {
		return false
	}
	for i := 0; i < len(in.Data.Data); i++ {
		if in.Data.Data[i] != other.Data.Data[i] {
			return false
		}
	}
	return true
}

func (in *shape) isSameWithRawShape2_Safe(other rawShape2) bool { //2b tested
	var allTheBits = in.Data.FindAllIndex()
	if len(allTheBits) != len(other.blocks) {
		panic("inside func (in *shape) isSame_Safe(other rawShape2) bool")
	}
	return in.isSameWithRawShape2(other)
}

func (in *shape) isSameWithRawShape2(other rawShape2) bool { //2 b tested
	var indece = in.Data.FindAllIndex()
	for _, index := range indece {
		var vec = indexToVec(index, in.Dim)
		var found = false
		for _, fromOther := range other.blocks {
			if vec == fromOther {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// func (in *shape) ____genBigger(other shape) (out []rawShape) { //2 b tested
// 	var bits = in.Data.FindAllIndex()
// 	//var pos = make([]vec3int, 0)

// 	var searchPos = rawShape{Dim: vec3int{in.Dim.X + 2, in.Dim.Y + 2, in.Dim.Z + 2}}
// 	for _, index := range bits {
// 		var vec = indexToVec(index, in)

// 		var tempVec = vec
// 		tempVec.X++
// 		tempVec.Y++
// 		searchPos.Data.SetBit(vecToIndex(tempVec, in))
// 		tempVec.Z += 2
// 		searchPos.Data.SetBit(vecToIndex(tempVec, in))

// 		tempVec = vec
// 		tempVec.X++
// 		tempVec.Z++
// 		searchPos.Data.SetBit(vecToIndex(tempVec, in))
// 		tempVec.Y += 2
// 		searchPos.Data.SetBit(vecToIndex(tempVec, in))

// 		tempVec = vec
// 		tempVec.Y++
// 		tempVec.Z++
// 		searchPos.Data.SetBit(vecToIndex(tempVec, in))
// 		tempVec.X += 2
// 		searchPos.Data.SetBit(vecToIndex(tempVec, in))
// 	}

// 	var moved = rawShape{Dim: vec3int{in.Dim.X + 2, in.Dim.Y + 2, in.Dim.Z + 2}}
// 	for _, index := range bits {
// 		var vec = indexToVec(index, in)

// 		var tempVec = vec
// 		tempVec.X++
// 		tempVec.Y++
// 		tempVec.Z++
// 		searchPos.Data.ResetBit(vecToIndex(tempVec, in))
// 		moved.Data.SetBit(vecToIndex(tempVec, in))
// 	}

// 	for _, addThisPos := range searchPos.Data.FindAllIndex() {
// 		var temp = copyShape(moved)
// 		temp.Data.SetBit(addThisPos)
// 		out = append(out, temp)
// 	}
// 	return out
// }

func (in *shape) genBiggerIntoRawShape2AndCleanToShape() (out []shape) { //2 b tested correct for 1->2
	var bits = in.Data.FindAllIndex()
	//var outDimShrinkBy1 = in.Dim
	//outDimShrinkBy1.X++
	//outDimShrinkBy1.Y++
	//outDimShrinkBy1.Z++

	var searchPos = rawShape{Dim: vec3int{in.Dim.X + 2, in.Dim.Y + 2, in.Dim.Z + 2}}
	searchPos.alignDataToDim()
	for _, index := range bits {
		var vec = indexToVec(index, in.Dim)

		var tempVec = vec
		tempVec.X++
		tempVec.Y++
		searchPos.Data.SetBit(vecToIndex(tempVec, searchPos.Dim))
		tempVec.Z += 2
		searchPos.Data.SetBit(vecToIndex(tempVec, searchPos.Dim))

		tempVec = vec
		tempVec.X++
		tempVec.Z++
		searchPos.Data.SetBit(vecToIndex(tempVec, searchPos.Dim))
		tempVec.Y += 2
		searchPos.Data.SetBit(vecToIndex(tempVec, searchPos.Dim))

		tempVec = vec
		tempVec.Y++
		tempVec.Z++
		searchPos.Data.SetBit(vecToIndex(tempVec, searchPos.Dim))
		tempVec.X += 2
		searchPos.Data.SetBit(vecToIndex(tempVec, searchPos.Dim))
	}

	var moved = rawShape{Dim: vec3int{in.Dim.X + 2, in.Dim.Y + 2, in.Dim.Z + 2}}
	moved.alignDataToDim()
	for _, index := range bits {
		var vec = indexToVec(index, in.Dim)

		var tempVec = vec
		tempVec.X++
		tempVec.Y++
		tempVec.Z++
		searchPos.Data.ResetBit(vecToIndex(tempVec, moved.Dim))
		moved.Data.SetBit(vecToIndex(tempVec, moved.Dim))
	}
	//now I have searchPos in rawShape as the halo, moved in rawShape as the core
	var movedInRawShape2 rawShape2
	for _, index := range moved.Data.FindAllIndex() {
		movedInRawShape2.addBlock(indexToVec(index, moved.Dim))
	}
	//now I have searchPos in rawShape as the halo, movedInRawShape2 in rawShape2 as the core
	var rawShapes = make([]rawShape2, 0, len(searchPos.Data.Data))
	for _, addThisPos := range searchPos.Data.FindAllIndex() {
		var temp = movedInRawShape2.copy()
		temp.addBlock(indexToVec(addThisPos, moved.Dim))
		//temp.dimSizeShrinkBy1 = outDimShrinkBy1
		temp.sort()
		rawShapes = append(rawShapes, temp)
	}
	//now I have new shapes in rawShape2, I need to clean them. Now they don't have dimention size, it's ok.
	out = make([]shape, 0, len(searchPos.Data.Data))
	for _, rawSh := range rawShapes {
		var sh = rawSh.clean()
		out = append(out, sh)
	}
	return out
}

func (in rawShape2) isBetterThan(other rawShape2) (better bool) { //2b tested
	if len(in.blocks) != len(other.blocks) {
		panic("inside rawShape2.inBetterThan function")
	}
	if in.dimSizeShrinkBy1 != other.dimSizeShrinkBy1 {
		panic("inside rawShape2.inBetterThan function")
	}
	for i, fromIn := range in.blocks {
		var fromOther = other.blocks[i]
		if fromIn.Z < fromOther.Z {
			return true
		}
		if fromIn.Z > fromOther.Z {
			return false
		}
		if fromIn.Y < fromOther.Y {
			return true
		}
		if fromIn.Y > fromOther.Y {
			return false
		}
		if fromIn.X < fromOther.X {
			return true
		}
		if fromIn.X > fromOther.X {
			return false
		}
	}
	return true
}

func (in rawShape2) same(other rawShape2) (same bool) { ///////
	if !in.sorted || !other.sorted {
		panic("inside func(in rawShape2) same ")
	}
	if len(in.blocks) != len(other.blocks) {
		panic("inside func(in rawShape2) same ")
	}
	if in.dimSizeShrinkBy1 != other.dimSizeShrinkBy1 {
		panic("inside rawShape2.inBetterThan function")
	}
	for i, fromIn := range in.blocks {
		var fromOther = other.blocks[i]
		if fromIn != fromOther {
			return false
		}
	}
	return true
}

func (in *rawShape2) moveToOriginAndShrinkDim() {
	//step1, move to origin
	var minX = len(in.blocks)
	var minY = minX
	var minZ = minX
	var maxX = 0
	var maxY = 0
	var maxZ = 0
	for _, v := range in.blocks {
		if v.X < minX {
			minX = v.X
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}

		if v.Z < minZ {
			minZ = v.Z
		}
		if v.Z > maxZ {
			maxZ = v.Z
		}
	}
	if minX > 1 || minY > 1 || minZ > 1 {
		panic("inside clean function.")
	}
	if minX > 0 {
		for i := range in.blocks {
			in.blocks[i].X--
		}
		maxX--
	}
	if minY > 0 {
		for i := range in.blocks {
			in.blocks[i].Y--
		}
		maxY--
	}
	if minZ > 0 {
		for i := range in.blocks {
			in.blocks[i].Z--
		}
		maxZ--
	}
	in.dimSizeShrinkBy1.X = maxX
	in.dimSizeShrinkBy1.Y = maxY
	in.dimSizeShrinkBy1.Z = maxZ
	return
}

func (in rawShape2) clean() (out shape) { //2 b tested  1
	in.moveToOriginAndShrinkDim()
	var candidates = make([]rawShape2, 0, 7)
	//step, try mirror it, comparing in max size representation. Doesn't have to be the final max size
	candidates = append(candidates, in.mirror(true, false, false))
	candidates = append(candidates, in.mirror(true, true, false))
	candidates = append(candidates, in.mirror(true, false, true))
	candidates = append(candidates, in.mirror(true, true, true))
	candidates = append(candidates, in.mirror(false, true, false))
	candidates = append(candidates, in.mirror(false, true, true))
	candidates = append(candidates, in.mirror(false, false, true))

	var best = in
	for _, v := range candidates {
		if v.isBetterThan(best) {
			best = v
		}
	}
	//best one into clean shape

	// var maxX, maxY, maxZ int
	// for _, v := range best.blocks {
	// 	if v.X > maxX {
	// 		maxX = v.X
	// 	}
	// 	if v.Y > maxY {
	// 		maxX = v.Y
	// 	}
	// 	if v.Z > maxZ {
	// 		maxX = v.Z
	// 	}
	// }

	out.Dim = vec3int{best.dimSizeShrinkBy1.X + 1, best.dimSizeShrinkBy1.Y + 1, best.dimSizeShrinkBy1.Z + 1}
	out.alignDataToDim()
	for _, vec3 := range best.blocks {
		out.Data.SetBit(vecToIndex(vec3, out.Dim))
	}
	return out
}

// func (in *shape) toMaxSize(maxSize, blockCount int) (out rawShape) { //2 b tested
// 	out.Dim.X = maxSize
// 	out.Dim.Y = maxSize
// 	out.Dim.Z = maxSize

// 	var already_found = 0
// 	this part should be replaced with new method   for out_i, v := range in.Data.Data {
// 		if v != 0 {
// 			for i := 0; i < 64; i++ {
// 				var mask = int64(1 << i)
// 				var masked = mask & v
// 				if masked != 0 {
// 					already_found++
// 					var index = out_i*64 + i
// 					var in_x, in_y, in_z = indexToCoord(index, in)
// 					index = coordToIndex(in_x, in_y, in_z, in)
// 					out.setBitTo1(index)
// 				}
// 				if already_found >= blockCount {
// 					break
// 				}
// 			}
// 		}
// 		if already_found >= blockCount {
// 			break
// 		}
// 	}
// 	return
// }

//I still need all the containers..

func makeStartPoint() (out shape) {
	out = shape{Dim: vec3int{1, 1, 1}}
	out.Data.SetLen(1)
	out.Data.SetBit(0)
	return out
}

func nextStep(currentBlockCount, currentBlockCount_minus1 *int) { //2b tested
	(*currentBlockCount)++
	(*currentBlockCount_minus1)++
}

type sameDimCont struct {
	data                 []shape
	sortedAndNoDuplicate bool
}

func (c *sameDimCont) init() { ///////////
	c.sortedAndNoDuplicate = false
}
func (c sameDimCont) add(sh *shape) (out sameDimCont) {
	c.data = append(c.data, *sh)
	out.data = c.data
	out.sortedAndNoDuplicate = false
	return out
}
func (c *sameDimCont) sortAndRemoveDuplicate() { ///////////
	sort.Sort(shapeSort(c.data))
	var temp = make([]shape, 1)
	copy(temp, c.data)
	//removes repeating elements.
	var e = c.data[0]
	for i := 1; i < len(c.data); i++ {
		if !c.data[i].isSame(e) {
			e = c.data[i]
			temp = append(temp, e)
		}
	}
	c.data = temp
	c.sortedAndNoDuplicate = true
	return
}

type mainContainer struct {
	data map[int]map[int]map[int]map[int]sameDimCont
	//...blockCount   X       Y       Z  and a container
	currentBlockCount        int
	currentBlockCount_minus1 int
	currentBlockCount_minus2 int
}

func (in *mainContainer) add(sh *shape, blockCount int) {
	var _, ok = in.data[blockCount]
	if !ok {
		in.data[blockCount] = make(map[int]map[int]map[int]sameDimCont, 0)
	}
	_, ok = in.data[blockCount][sh.Dim.X]
	if !ok {
		in.data[blockCount][sh.Dim.X] = make(map[int]map[int]sameDimCont, 0)
	}
	_, ok = in.data[blockCount][sh.Dim.X][sh.Dim.Y]
	if !ok {
		in.data[blockCount][sh.Dim.X][sh.Dim.Y] = make(map[int]sameDimCont, 0)
	}
	_, ok = in.data[blockCount][sh.Dim.X][sh.Dim.Y][sh.Dim.Z]
	if !ok {
		in.data[blockCount][sh.Dim.X][sh.Dim.Y][sh.Dim.Z] = sameDimCont{}
	}
	in.data[blockCount][sh.Dim.X][sh.Dim.Y][sh.Dim.Z] = in.data[blockCount][sh.Dim.X][sh.Dim.Y][sh.Dim.Z].add(sh)
}
func (in *mainContainer) sortAndRemoveDuplicate(blockCount int) { ////////////
	for x, _ := range in.data[blockCount] {
		for y, _ := range in.data[blockCount][x] {
			for z, _ := range in.data[blockCount][x][y] {
				var temp = in.data[blockCount][x][y][z]
				temp.sortAndRemoveDuplicate()
				in.data[blockCount][x][y][z] = temp
			}
		}
	}
}
func (in *mainContainer) init() { ////////////
	var temp = makeStartPoint()
	in.currentBlockCount = 1
	in.currentBlockCount_minus1 = 0
	in.currentBlockCount_minus2 = -1
	in.data = make(map[int]map[int]map[int]map[int]sameDimCont)
	in.add(&temp, 1)
}
func (in *mainContainer) nextStep() { ////////////
	in.currentBlockCount++
	in.currentBlockCount_minus1++
	in.currentBlockCount_minus2++
	in.data[in.currentBlockCount] = make(map[int]map[int]map[int]sameDimCont)

	//var candidates = make(map[int64]shape, 0)
	//step 1, generate new clean shapes,
	for _, YZs := range in.data[in.currentBlockCount_minus1] {
		//var xShrinkBy1 = x - 1
		for _, Zs := range YZs {
			//var yShrinkBy1 = y - 1
			for _, innerCont := range Zs {
				//var zShrinkBy1 = z - 1
				for _, shapeFromBefore := range innerCont.data {
					var newShapes = shapeFromBefore.genBiggerIntoRawShape2AndCleanToShape()
					for _, sh := range newShapes {
						in.add(&sh, in.currentBlockCount)
						//in.data[in.currentBlockCount][x][y][z] = in.data[in.currentBlockCount][x][y][z].add(&sh)
					}
				}
			}
		}
	}
	if SAVE_MEM_MODE {
		in.clearMemForBlockCount(in.currentBlockCount_minus1)
	}
	//step 2, sort the inner cont
	in.sortAndRemoveDuplicate(in.currentBlockCount)
}
func (in *mainContainer) clearMemForBlockCount(blockCount int) { ////////////
	delete(in.data, blockCount)
}

func (in *mainContainer) report() { ////////////
	var shapeCountWithoutSymmetric = 0
	var shapeCountWithSymmetric = 0

	for _, YZs := range in.data[in.currentBlockCount] {
		for _, Zs := range YZs {
			for _, innerCont := range Zs {
				shapeCountWithoutSymmetric += len(innerCont.data)
				for _, sh := range innerCont.data {
					var mirrored = sh.mirror(true, false, false)
					if sh.isSameWithRawShape2_Safe(mirrored) {
						shapeCountWithSymmetric++
						continue
					}
					mirrored = sh.mirror(false, true, false)
					if sh.isSameWithRawShape2_Safe(mirrored) {
						shapeCountWithSymmetric++
						continue
					}
					mirrored = sh.mirror(false, false, true)
					if sh.isSameWithRawShape2_Safe(mirrored) {
						shapeCountWithSymmetric++
						continue
					}
					shapeCountWithSymmetric += 2
				}
			}
		}
	}
	fmt.Printf("Block count:%d, shape count without symmetric:%d, shape count with symmetric:%d", in.currentBlockCount, shapeCountWithoutSymmetric,
		shapeCountWithSymmetric)
	fmt.Println()
}

func main() {
	var data = mainContainer{}
	data.init() //needs no sortation.
	//data.report()

	for data.currentBlockCount <= MaxStep {
		data.nextStep()
		//data.report()the report func was not finished. Im out of patience. Read the comment in the beginning for the idea of how to finish it.
	}

	var a = 321
	fmt.Print(a)
}
