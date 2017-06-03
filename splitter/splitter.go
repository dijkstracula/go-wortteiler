package splitter

// partition is a tuple of strings.
type partition struct {
	prefix string
	suffix string
}

// SplitNode is a partitioning of a dictionary word into two smaller ones.
type SplitNode struct {
	word   string
	prefix *SplitNode
	suffix *SplitNode
}

// offsets produces [0..n] starting with n/2 and working outwards.
func offsets(n int) []int {
	var ret []int

	if n < 1 {
		return ret
	}

	ret = make([]int, 0, n)

	mid := n / 2
	down := mid - 1
	up := mid

	for i := 0; i < mid; i++ {
		ret = append(ret, up)
		ret = append(ret, down)

		down--
		up++
	}

	if n%2 == 1 {
		ret = append(ret, up)
	}

	return ret
}

// partitions produces all partitions of a given string, starting
// with splits on the end and working inward.  This
// ordering is done to favour "more even" word split choices
// as we iterate further through the string.
func partitions(str string) []partition {
	var ret []partition

	//TODO: obviously this is silly and offsets should just produce
	//the offsets in the reversed order.
	offs := offsets(len(str) - 1)
	for i := 0; i < len(offs)/2; i++ {
		t := offs[i]
		offs[i] = offs[len(offs)-1-i]
		offs[len(offs)-1-i] = t
	}

	if len(offs) == 0 {
		return ret
	}

	ret = make([]partition, 0, len(offs))

	for _, off := range offs {
		ret = append(ret, partition{str[:off+1], str[off+1:]})
	}

	return ret
}

// memoize produces a memoized copy of `outer`.
// Parametric polymorphism?  We don't need no stinkin' parametric
// polymorphism *revs up 180 horsepower copy/paste engine*
func memoize(outer func(string) *SplitNode) func(string) *SplitNode {
	cache := make(map[string]*SplitNode)
	return func(str string) *SplitNode {
		if node, ok := cache[str]; ok {
			return node
		}
		node := outer(str)
		cache[str] = node
		return node
	}
}

// Splitter produces a function that uses the supplied function to
// generate a tree of word splits.
func Splitter(valid func(string) bool) func(string) *SplitNode {
	var fn func(str string) *SplitNode

	fn = func(str string) *SplitNode {
		var node *SplitNode

		// Case 0: If the node is too short to partition, just
		// return a leaf node.
		if len(str) <= 1 {
			if valid(str) {
				return &SplitNode{str, nil, nil}
			}
			return nil
		}

		// Look at all partitions from least desirable to most
		// (that is, how similar in length the partitions are)
		// and choose the best.
		for _, partition := range partitions(str) {
			prefTree := fn(partition.prefix)
			suffTree := fn(partition.suffix)

			// Case 1: This is a valid leaf node (the string is valid, but
			// this particular split doesn't yield two valid substrings)
			if valid(str) && (prefTree == nil || suffTree == nil) {
				node = &SplitNode{str, nil, nil}
			}

			// Case 2: This is a valid split: this will certainly be a valid
			// frond node in the tree; if the supplied string isn't a valid
			// word, however, don't include it in the node.
			if prefTree != nil && suffTree != nil {
				node = &SplitNode{str, prefTree, suffTree}
				if !valid(str) {
					node.word = ""
				}
			}
		}

		return node
	}

	return memoize(fn)
}
