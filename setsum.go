package setsum

import (
	"crypto/sha3"
	"encoding/binary"
	"encoding/hex"
)

const SetsumBytes = 32
const SetsumBytesPerColumn = 4
const SetsumColumns = SetsumBytes / SetsumBytesPerColumn

var SetsumPrimes = [...]uint32{
	4294967291, 4294967279, 4294967231, 4294967197, 4294967189, 4294967161, 4294967143, 4294967111}

type Setsum struct {
	state [SetsumColumns]uint32
}

func NewSetsum() *Setsum {
	return &Setsum{
		state: [SetsumColumns]uint32{},
	}
}

func (s *Setsum) Insert(item []byte) {
	s.InsertMany([][]byte{item})
}

func (s *Setsum) InsertMany(items [][]byte) {
	s.state = addState(s.state, itemsToState(items))
}

func (s *Setsum) Remove(item []byte) {
	s.RemoveMany([][]byte{item})
}

func (s *Setsum) RemoveMany(items [][]byte) {
	itemState := itemsToState(items)
	invertedState := invertState(itemState)
	s.state = addState(s.state, invertedState)
}

func (s *Setsum) Add(other *Setsum) {
	s.state = addState(s.state, other.state)
}

func (s *Setsum) Subtract(other *Setsum) {
	invertedState := invertState(other.state)
	s.state = addState(s.state, invertedState)
}

func (s *Setsum) Digest() [SetsumBytes]byte {
	itemHash := [SetsumBytes]byte{}
	for i := 0; i < SetsumColumns; i++ {
		binary.LittleEndian.PutUint32(itemHash[i*SetsumBytesPerColumn:(i+1)*SetsumBytesPerColumn], s.state[i])
	}
	return itemHash
}

func (s *Setsum) HexDigest() string {
	digest := s.Digest()
	return hex.EncodeToString(digest[:])
}

func addState(lhs [SetsumColumns]uint32, rhs [SetsumColumns]uint32) [SetsumColumns]uint32 {
	result := [SetsumColumns]uint32{}
	for i := 0; i < SetsumColumns; i++ {
		sum := uint64(lhs[i]) + uint64(rhs[i])
		if sum >= uint64(SetsumPrimes[i]) {
			sum -= uint64(SetsumPrimes[i])
		}
		result[i] = uint32(sum)
	}
	return result
}

func invertState(state [SetsumColumns]uint32) [SetsumColumns]uint32 {
	result := [SetsumColumns]uint32{}
	for i := 0; i < SetsumColumns; i++ {
		result[i] = SetsumPrimes[i] - state[i]
	}
	return result
}

func hashToState(hash *[SetsumBytes]byte) [SetsumColumns]uint32 {
	itemState := [SetsumColumns]uint32{}
	for i := 0; i < SetsumColumns; i++ {
		idx := i * SetsumBytesPerColumn
		num := binary.LittleEndian.Uint32(hash[idx : idx+4])
		if num >= SetsumPrimes[i] {
			num -= SetsumPrimes[i]
		}
		itemState[i] = num
	}
	return itemState
}

func itemsToState(item [][]byte) [SetsumColumns]uint32 {
	hasher := sha3.New256()
	for _, piece := range item {
		hasher.Write(piece)
	}
	var hashBytes [SetsumBytes]byte
	copy(hashBytes[:], hasher.Sum(nil))
	return hashToState(&hashBytes)
}
