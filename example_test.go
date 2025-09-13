package setsum

import (
	"fmt"
)

func ExampleSetsum() {
	s1 := NewSetsum()

	// add some items
	s1.Insert([]byte("hello"))
	s1.Insert([]byte("world"))

	// add same items but in different order
	s2 := NewSetsum()
	s2.Insert([]byte("world"))
	s2.Insert([]byte("hello"))

	fmt.Printf("Digests equal: %t\n", s1.Digest() == s2.Digest())
	// Output: Digests equal: true
}

func ExampleSetsum_Insert() {
	s := NewSetsum()
	s.Insert([]byte("hello"))
	s.Insert([]byte("hello"))

	// remove the dupe
	s.Remove([]byte("hello"))

	c := NewSetsum()
	c.Insert([]byte("hello"))
	fmt.Printf("`hello hello` after removing one 'hello' equals digest of 'hello': %t\n", s.Digest() == c.Digest())
	// Output: `hello hello` after removing one 'hello' equals digest of 'hello': true
}

func ExampleSetsum_Remove() {
	s := NewSetsum()
	s.Insert([]byte("hello"))
	s.Insert([]byte("world"))

	// remove in different order
	s.Remove([]byte("hello"))

	{
		world := NewSetsum()
		world.Insert([]byte("world"))
		fmt.Printf("`hello world` after removing 'hello' equals digest of 'world': %t\n", s.Digest() == world.Digest())
	}

	s.Remove([]byte("world"))

	empty := NewSetsum()
	fmt.Printf("Digest after removing all items equals empty digest: %t\n", s.Digest() == empty.Digest())
	// Output:
	// `hello world` after removing 'hello' equals digest of 'world': true
	// Digest after removing all items equals empty digest: true
}
