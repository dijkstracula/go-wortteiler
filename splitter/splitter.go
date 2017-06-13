package splitter

import "fmt"

// partition is a tuple of strings.
type partition struct {
	prefix string
	suffix string
}

// Node is a partitioning of a dictionary word into two smaller ones.
type Node struct {
	Word   string `json:"defn,omitempty"`
	Prefix *Node  `json:"prefix,omitempty"`
	Suffix *Node  `json:"suffix,omitempty"`
}

// String makes an string.
func (n *Node) String() string {
	if n == nil {
		return "<nil>"
	}

	return fmt.Sprintf("{%#v %s %s}", n.Word, n.Prefix.String(), n.Suffix.String())
}

// MakeNode produces a new Node.
func MakeNode(word string, prefix, suffix *Node) *Node {
	var n Node

	n.Word = word
	n.Prefix = prefix
	n.Suffix = suffix

	return &n
}

// MakeLeaf produces a leaf node (e.g. with both children nil)
func MakeLeaf(word string) *Node {
	var n Node
	n.Word = word

	return &n
}

// ForEach invokes the function `f` as an inorder traversal on the Tree.
func (n *Node) ForEach(f func(*Node)) {
	if n == nil {
		return
	}

	n.Prefix.ForEach(f)
	f(n)
	n.Suffix.ForEach(f)
}

// Score scores a Node.
// I have no idea if this score function makes sense, but, two things seem to be true:
// 1) We want trees with a higher proportion of valid words to be weighted higher,
// 2) A node with a longer valid word maybe should be weighted higher too?
func (n *Node) Score() (int, int) {
	num := 0
	den := 1
	if n == nil {
		return num, den
	}

	if n.Word != "" {
		num++
	}

	pn, pd := n.Prefix.Score()
	sn, sd := n.Suffix.Score()

	num += pn + sn
	den += pd + sd

	return num, den
}

// SplitFunc consumes a string and produces a tree.
type SplitFunc func(string) *Node

// partitions produces all partitions of a given string, starting
// with left-most splits and working right.  This
// ordering is done to favour the best possible right-most split,
// so the left half of the final splitted string includes any possible
// suffix.
func partitions(str string) []partition {
	var ret []partition

	n := len(str) - 1
	if n <= 0 {
		return ret
	}

	ret = make([]partition, 0, n)

	//for _, off := range offs {
	for off := 0; off < n; off++ {
		ret = append(ret, partition{str[:off+1], str[off+1:]})
	}

	return ret
}

// memoize produces a memoized copy of `outer`.
// Parametric polymorphism?  We don't need no stinkin' parametric
// polymorphism *revs up 180 horsepower copy/paste engine*
func memoize(outer func(string) *Node) SplitFunc {
	cache := make(map[string]*Node)
	return func(str string) *Node {
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
func Splitter(valid func(string) bool) SplitFunc {
	var fn SplitFunc

	fn = memoize(func(str string) *Node {
		var tree *Node
		treeN := 0
		treeD := 1

		// Case 0: If the node is too short to partition, just
		// return a leaf node.
		if len(str) <= 1 {
			if valid(str) {
				return MakeNode(str, nil, nil)
			}
			return nil
		}

		// Look at all partitions from least desirable to most
		// (that is, how similar in length the partitions are)
		// and choose the best.
		for _, partition := range partitions(str) {
			prefTree := fn(partition.prefix)
			suffTree := fn(partition.suffix)

			if valid(str) && (prefTree == nil || suffTree == nil) {
				// Case 1: This is a valid leaf node (the string is valid, but
				// this particular split doesn't yield two valid substrings)
				if tree == nil || (tree.Prefix == nil && tree.Suffix == nil) {
					// Never fall back from case 2 to case 1!
					tree = MakeNode(str, nil, nil)
				}
			} else if prefTree != nil && suffTree != nil {
				// Case 2: This is a valid split: this will certainly be a valid
				// frond node in the tree; if the supplied string isn't a valid
				// word, however, don't include it in the node.
				newTree := MakeNode(str, prefTree, suffTree)
				if !valid(str) {
					newTree.Word = ""
				}

				if n, d := newTree.Score(); float64(n)/float64(d) > float64(treeN)/float64(treeD) {
					tree = newTree
					treeN = n
					treeD = d
				}
			}
		}

		return tree
	})

	return fn
}
