package local

import(
	"fmt"
	"log"
)

func main() {
	arr := []int{1}
	bTree := Init(arr)
	bTree.PrintTree()
	node := bTree.Query(1)
	fmt.Println(node)
	bTree.Remove(node)
	bTree.PrintTree()
}

func Init(arr []int) *BTree {
	bt := &BTree{}
	for _,a := range arr {
		tn := &TreeNode{Key: a}
		bt.Insert(tn)
	}
	return bt
}

type TreeNode struct {
	Key int
	Red bool
	P, L, R *TreeNode
	H int
}

func (node *TreeNode) LeftRotate(bTree *BTree) {
	if node.R == nil {
		log.Fatal("nil right")
	}
	if node.P != nil {
		if node == node.P.L {
			node.P.L = node.R
		} else {
			node.P.R = node.R
		}
		node.R.P = node.P
	} else {
		node.R.P = nil
	}
	
	nrl := node.R.L
	node.P = node.R
	node.R.L = node
	node.R = nrl
	if nrl != nil {
		nrl.P = node
	}
		
	node.P.H--
	
	CalculateHeight(node.P.R, node.P.H+1)
	node.H++
	CalculateHeight(node.L, node.H+1)
	
	
	if node.P.P == nil {
		bTree.Root = node.P
	}
}

func CalculateHeight(tNode *TreeNode, H int) {
	if tNode == nil {
		return
	}
	tNode.H = H
	if tNode.L != nil {
		CalculateHeight(tNode.L, H+1)
	}
	if tNode.R != nil {
		CalculateHeight(tNode.R, H+1)
	}
}

func (node *TreeNode) RightRotate(bTree *BTree) {
	if node.L == nil {
		log.Fatal("nil left")
	}
	if node.P != nil {
		if node == node.P.L {
			node.P.L = node.L
		} else {
			node.P.R = node.L
		}
		node.L.P = node.P
	} else {
		node.L.P = nil
	}
	
	nlr := node.L.R
	node.P = node.L
	node.L.R = node 
	node.L = nlr
	if nlr != nil {
		nlr.P = node
	}
	node.P.H--
	
	CalculateHeight(node.P.L, node.P.H+1)
	node.H++
	CalculateHeight(node.R, node.H+1)
	
	if node.P.P == nil {
		bTree.Root = node.P
	}
}

type BTree struct {
	Root *TreeNode
}

type BTreeI interface {
	Query(Key interface{}) *TreeNode
	Insert(node TreeNode)
	Remove(node TreeNode)
}

func (bTree *BTree) PrintTree() {
	Root := bTree.Root
	if Root == nil {
		fmt.Println("empty tree")
		return
	}
	PrintNode(Root)
}

func PrintNode(tNode *TreeNode) {
	if tNode.L != nil {
		PrintNode(tNode.L)
	}
	fmt.Printf("%p %v\n", tNode, tNode)
	if tNode.R != nil {
		PrintNode(tNode.R)
	}

}

func (bTree *BTree) Query(Key interface{}) *TreeNode {
	if bTree == nil || bTree.Root == nil {
		log.Fatal("nil tree")
	}
	
	var node *TreeNode
	node = bTree.Root
	
	var val int;
	switch t := Key.(type) {
		case int: val = Key.(int)
		default: log.Fatal("unkown type", t)
	}
	
	for {
		if node == nil {
			return nil
		}
		if node.Key == val {
			return node
		} else if node.Key < val {
			node = node.L
		} else {
			node = node.R
		}
	}
	
	return nil
}

func (bTree *BTree) Insert(tNode *TreeNode) {
	if bTree == nil {
		log.Fatal("nil tree")
	}
	if tNode == nil {
		log.Fatal("insert nil")
	}
	
	tNode.L,tNode.R,tNode.P = nil,nil,nil
	tNode.Red = true
	
	if bTree.Root == nil {
		tNode.H = 0
		tNode.Red = false
		bTree.Root = tNode
		return
	}
	
	
	Key := tNode.Key
	node := bTree.Root
	
	H := 0
	
	for {
		if tNode.P != nil {
			break
		}
		if node.Key == Key {
			log.Fatal("exist same Key")
		}
		if node.Key < Key {
			if node.R == nil {
				node.R = tNode
				tNode.P = node
			} else {
				node = node.R
			}
		} else {
			if node.L == nil {
				node.L = tNode
				tNode.P = node
			} else {
				node = node.L
			}
		}
		H++
	}
	tNode.H = H
	
	if !node.Red {
		return
	} 
	
	for {
		if node.P == nil {
			break
		}
		if node.P.L == node { 
			if node.P.R != nil && node.P.R.Red {
				node.Red = false
				node.P.R.Red = false
				node.P.Red = true
				node = node.P
			} else {
				if node.R == tNode {
					node.LeftRotate(bTree)
					node = tNode
				}
				node.Red = false
				node.P.Red = true
				node.P.RightRotate(bTree)
				break
			}
		} else {
			if node.P.L != nil && node.P.L.Red {
				node.Red = false
				node.P.L.Red = false
				node.P.Red = true
				node = node.P
			} else {
				if node.L == tNode {
					node.RightRotate(bTree)
					node = tNode
				}
				node.Red = false
				node.P.Red = true
				node.P.LeftRotate(bTree)
				break
			}
		}
	}
	node.Red = false
}

func (bTree *BTree) Remove(tNode *TreeNode) {
	if tNode == nil {
		log.Fatal("nil node")
	}
	if bTree == nil {
		log.Fatal("nil tree")
	}
	
	node := tNode
	if node.L != nil && node.R != nil {
		node = node.R
		for {
			if node.L != nil {
				node = node.L
			} else {
				break
			}
		}
	}
	if node.Red {
		return
	}
	if node.P == nil {
		if node.L == nil && node.R == nil {
			bTree.Root = nil
			return
		}
		if node.R != nil {
			node, bTree.Root = node.R, node.R
		} else {
			node, bTree.Root = node.L, node.R
		}
		node.Red = false
		node.H--
		if node.L != nil {
			CalculateHeight(node.L, node.L.H-1)
		}
		if node.R != nil {
			CalculateHeight(node.R, node.R.H-1)
		}
		return
	}
	
	rNode := node
	for {
		if node.P == nil || node.Red {
			break
		}
		if node == node.P.L {
			if node.P.R.Red {
				node.P.Red = true
				node.P.R.Red = false
				node.P.LeftRotate(bTree)
			}
			if node.P.R.L == nil && node.P.R.R == nil {
				node.P.R.Red = true
				node.P.Red = false
				node = node.P
			} else {
				if(node.P.R.L != nil && node.P.R.R == nil) {
					node.P.R.Red = true
					node.P.R.L.Red = false
					node.P.R.RightRotate(bTree)
				}
				node.P.R.Red = node.P.Red
				node.P.Red = false
				node.P.R.R.Red = false
				node.P.LeftRotate(bTree)
				break
			}
		} else {
			if node.P.L.Red {
				node.P.Red = true
				node.P.L.Red = false
				node.P.RightRotate(bTree)
			}
			if node.P.L.L == nil && node.P.L.R == nil {
				node.P.L.Red = true
				node.P.Red = false
				node = node.P
			} else {
				if node.P.L.L != nil && node.P.L.R == nil {
					node.P.L.L.Red = false
					node.P.L.Red = true
					node.P.L.LeftRotate(bTree)
				}
				node.P.L.Red = node.P.Red
				node.P.L.R.Red = false
				node.P.Red = false
				node.P.RightRotate(bTree)
				break
			}
		}
	}
	node.Red = false
	
	if rNode.R != nil {
		rNode.R.P = rNode.P
		if (rNode == rNode.P.L) {
			rNode.P.L = rNode.R
		} else {
			rNode.P.R = rNode.R
		}
		CalculateHeight(rNode.R, rNode.R.H-1)
	} else if rNode.L != nil {
		rNode.L.P = rNode.P
		if rNode == rNode.P.L {
			rNode.P.L = rNode.L
		} else {
			rNode.P.R = rNode.L
		}
		CalculateHeight(rNode.L, rNode.L.H-1)
	} else {
		if rNode.P.L == rNode {
			rNode.P.L = nil
		} else {
			rNode.P.R = nil
		}
	}
	
}



