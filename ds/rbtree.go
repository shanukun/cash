package ds

const (
	RED   bool = true
	BLACK bool = false
)

type Node struct {
	key   string
	Value interface{}
	color bool

	parent *Node
	left   *Node
	right  *Node
}

type RBTree struct {
	Nil  *Node
	Root *Node
}

func (tree *RBTree) leftRotate(x *Node) {
	y := x.right
	x.right = y.left

	if y.left != tree.Nil {
		y.left.parent = x
	}

	y.parent = x.parent

	if x.parent == nil {
		tree.Root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

func (tree *RBTree) rightRotate(x *Node) {
	y := x.left
	x.left = y.right

	if y.right != tree.Nil {
		y.right.parent = x
	}

	y.parent = x.parent
	if x.parent == nil {
		tree.Root = y
	} else if x == x.parent.right {
		x.parent.right = y
	} else {
		x.parent.left = y
	}

	y.right = x
	x.parent = y
}

func (tree *RBTree) Insert(key string, value interface{}) {
	z := &Node{
		key:   key,
		Value: value,
		color: BLACK,
		left:  tree.Nil,
		right: tree.Nil,
	}

	var y *Node
	x := tree.Root

	for x != tree.Nil {
		y = x
		if z.key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}

	z.parent = y
	if y == nil {
		tree.Root = z
	} else if z.key < y.key {
		y.left = z
	} else {
		y.right = z
	}

	tree.insertFixup(z)
}

func (tree *RBTree) insertFixup(z *Node) {
	for z.parent != nil && z.parent.color == RED {
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.right {
					z = z.parent
					tree.leftRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
			}
		} else {
			y := z.parent.parent.left
			if y.color == RED {
				z.parent.color = BLACK
				y.color = BLACK
				z.parent.parent.color = RED
				z = z.parent.parent
			} else {
				if z == z.parent.left {
					z = z.parent
					tree.rightRotate(z)
				}
				z.parent.color = BLACK
				z.parent.parent.color = RED
				tree.leftRotate(z.parent.parent)
			}
		}
		if z == tree.Root {
			break
		}
	}
	tree.Root.color = BLACK
}

func (tree *RBTree) Delete(key string) {
	z := tree.search(key)

	if z == tree.Nil {
		return
	}

	y := z
	y_orig_color := y.color
	var x *Node

	if z.left == tree.Nil {
		// case 1
		x = z.right
		tree.transplant(z, z.right)
	} else if z.right == tree.Nil {
		// case 2
		x = z.left
		tree.transplant(z, z.left)
	} else {
		// case 3
		y = tree.minimum(z.right)
		y_orig_color = y.color
		x = y.right

		if y.parent == z {
			x.parent = y
		} else {
			tree.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}

		tree.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color

	}
	if y_orig_color == BLACK {
		tree.deleteFixup(x)
	}
}

func (tree *RBTree) deleteFixup(x *Node) {
	for x != tree.Root && x.color == BLACK {
		if x == x.parent.right {
			w := x.parent.left

			// type 1
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				tree.leftRotate(x.parent)
				w = x.parent.right
			}

			// type 2
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				// type 3
				if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					tree.rightRotate(w)
					w = x.parent.right
				}

				// type 4
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				tree.leftRotate(x.parent)
				x = tree.Root
			}
		} else {
			w := x.parent.left
			// type 1
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				tree.rightRotate(x.parent)
				w = x.parent.left
			}
			// type 2
			if w.right.color == BLACK && w.left.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				// type 3
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					tree.leftRotate(w)
					w = x.parent.left
				}

				// type 4
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK

				tree.rightRotate(x.parent)
				x = tree.Root
			}
		}
	}
	x.color = BLACK
}

func (tree *RBTree) transplant(u, v *Node) {
	if u.parent == nil {
		tree.Root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}
	v.parent = u.parent
}

func (tree *RBTree) minimum(x *Node) *Node {
	for x.left != tree.Nil {
		x = x.left
	}
	return x
}

func (tree *RBTree) search(key string) *Node {
	x := tree.Root
	for x != tree.Nil && key != x.key {
		if key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}
	return x
}

func (tree *RBTree) Find(key string) (interface{}, bool) {
	item := tree.search(key)
	if item.key == key {
		return item.Value, true
	}

	return nil, false
}
