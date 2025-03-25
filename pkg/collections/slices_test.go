package collections

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/pkg/profile"
)

func TestFilter(t *testing.T) {

	numbers := intSlice(100_000_000)

	filtered := Filter(Iter(numbers), func(i int) bool {
		return i%2 == 0
	})

	s := slices.Collect(filtered)

	t.Log(len(s))

}

func TestFilter1(t *testing.T) {

	numbers := intSlice(100_000_000)

	filtered, cnt := FilterCount(Iter(numbers), func(i int) bool {
		return i%2 == 0
	})

	out := make([]int, 0, cnt)
	for i := range filtered {
		out = append(out, i)
	}

	t.Log(len(out))

}

func TestFilter2(t *testing.T) {

	numbers := intSlice(100_000_000)

	var out []int
	for _, n := range numbers {
		if n%2 == 0 {
			out = append(out, n)
		}
	}

	t.Log(len(out))

}

func intSlice(n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i
	}
	return out
}

func randomIntSlice(n int) []int {
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		slice[i] = rand.Intn(n / 2) // This will ensure some duplicates
	}
	return slice
}

func BenchmarkUnique(b *testing.B) {
	defer profile.Start(
		profile.MemProfile,
		profile.ProfilePath("./mem.out"),
		profile.MemProfileRate(1),
		profile.NoShutdownHook,
	).Stop()
	sizes := []int{
		100000,
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("unique_int_size_%d", size), func(b *testing.B) {
			slice := randomIntSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique(slice)
			}
		})

	}
}
