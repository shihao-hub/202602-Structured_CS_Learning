package tree

import "fmt"

// TreeNode 二叉树节点
type TreeNode struct {
	Value int
	Left  *TreeNode
	Right *TreeNode
}

// NewTreeNode 创建新节点
func NewTreeNode(value int) *TreeNode {
	return &TreeNode{
		Value: value,
		Left:  nil,
		Right: nil,
	}
}

// BinarySearchTree 二叉搜索树
type BinarySearchTree struct {
	Root *TreeNode
	Size int
}

// NewBinarySearchTree 创建二叉搜索树
func NewBinarySearchTree() *BinarySearchTree {
	return &BinarySearchTree{
		Root: nil,
		Size: 0,
	}
}

// Insert 插入节点
func (bst *BinarySearchTree) Insert(value int) {
	bst.Root = bst.insertNode(bst.Root, value)
	bst.Size++
}

func (bst *BinarySearchTree) insertNode(node *TreeNode, value int) *TreeNode {
	if node == nil {
		return NewTreeNode(value)
	}

	if value < node.Value {
		node.Left = bst.insertNode(node.Left, value)
	} else if value > node.Value {
		node.Right = bst.insertNode(node.Right, value)
	}
	// 如果值相等，不插入（BST通常不允许重复值）

	return node
}

// Search 查找节点
func (bst *BinarySearchTree) Search(value int) *TreeNode {
	return bst.searchNode(bst.Root, value)
}

func (bst *BinarySearchTree) searchNode(node *TreeNode, value int) *TreeNode {
	if node == nil || node.Value == value {
		return node
	}

	if value < node.Value {
		return bst.searchNode(node.Left, value)
	}
	return bst.searchNode(node.Right, value)
}

// Contains 检查是否包含某值
func (bst *BinarySearchTree) Contains(value int) bool {
	return bst.Search(value) != nil
}

// Delete 删除节点
func (bst *BinarySearchTree) Delete(value int) bool {
	if !bst.Contains(value) {
		return false
	}
	bst.Root = bst.deleteNode(bst.Root, value)
	bst.Size--
	return true
}

func (bst *BinarySearchTree) deleteNode(node *TreeNode, value int) *TreeNode {
	if node == nil {
		return nil
	}

	if value < node.Value {
		node.Left = bst.deleteNode(node.Left, value)
	} else if value > node.Value {
		node.Right = bst.deleteNode(node.Right, value)
	} else {
		// 找到要删除的节点
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		}

		// 有两个子节点，找到右子树的最小值
		minNode := bst.findMin(node.Right)
		node.Value = minNode.Value
		node.Right = bst.deleteNode(node.Right, minNode.Value)
	}

	return node
}

// FindMin 找到最小值
func (bst *BinarySearchTree) FindMin() *TreeNode {
	if bst.Root == nil {
		return nil
	}
	return bst.findMin(bst.Root)
}

func (bst *BinarySearchTree) findMin(node *TreeNode) *TreeNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// FindMax 找到最大值
func (bst *BinarySearchTree) FindMax() *TreeNode {
	if bst.Root == nil {
		return nil
	}
	return bst.findMax(bst.Root)
}

func (bst *BinarySearchTree) findMax(node *TreeNode) *TreeNode {
	current := node
	for current.Right != nil {
		current = current.Right
	}
	return current
}

// InorderTraversal 中序遍历（升序）
func (bst *BinarySearchTree) InorderTraversal() []int {
	result := make([]int, 0, bst.Size)
	bst.inorderHelper(bst.Root, &result)
	return result
}

func (bst *BinarySearchTree) inorderHelper(node *TreeNode, result *[]int) {
	if node != nil {
		bst.inorderHelper(node.Left, result)
		*result = append(*result, node.Value)
		bst.inorderHelper(node.Right, result)
	}
}

// PreorderTraversal 前序遍历
func (bst *BinarySearchTree) PreorderTraversal() []int {
	result := make([]int, 0, bst.Size)
	bst.preorderHelper(bst.Root, &result)
	return result
}

func (bst *BinarySearchTree) preorderHelper(node *TreeNode, result *[]int) {
	if node != nil {
		*result = append(*result, node.Value)
		bst.preorderHelper(node.Left, result)
		bst.preorderHelper(node.Right, result)
	}
}

// PostorderTraversal 后序遍历
func (bst *BinarySearchTree) PostorderTraversal() []int {
	result := make([]int, 0, bst.Size)
	bst.postorderHelper(bst.Root, &result)
	return result
}

func (bst *BinarySearchTree) postorderHelper(node *TreeNode, result *[]int) {
	if node != nil {
		bst.postorderHelper(node.Left, result)
		bst.postorderHelper(node.Right, result)
		*result = append(*result, node.Value)
	}
}

// LevelOrderTraversal 层序遍历（广度优先）
func (bst *BinarySearchTree) LevelOrderTraversal() []int {
	if bst.Root == nil {
		return []int{}
	}

	result := make([]int, 0, bst.Size)
	queue := []*TreeNode{bst.Root}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, node.Value)

		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}

	return result
}

// Height 获取树的高度
func (bst *BinarySearchTree) Height() int {
	return bst.heightHelper(bst.Root)
}

func (bst *BinarySearchTree) heightHelper(node *TreeNode) int {
	if node == nil {
		return -1
	}

	leftHeight := bst.heightHelper(node.Left)
	rightHeight := bst.heightHelper(node.Right)

	if leftHeight > rightHeight {
		return leftHeight + 1
	}
	return rightHeight + 1
}

// IsEmpty 检查树是否为空
func (bst *BinarySearchTree) IsEmpty() bool {
	return bst.Root == nil
}

// IsBalanced 检查树是否平衡
func (bst *BinarySearchTree) IsBalanced() bool {
	return bst.isBalancedHelper(bst.Root)
}

func (bst *BinarySearchTree) isBalancedHelper(node *TreeNode) bool {
	if node == nil {
		return true
	}

	leftHeight := bst.heightHelper(node.Left)
	rightHeight := bst.heightHelper(node.Right)

	diff := leftHeight - rightHeight
	if diff < -1 || diff > 1 {
		return false
	}

	return bst.isBalancedHelper(node.Left) && bst.isBalancedHelper(node.Right)
}

// PrintTree 打印树结构
func (bst *BinarySearchTree) PrintTree() {
	fmt.Printf("BST (size=%d, height=%d):\n", bst.Size, bst.Height())
	bst.printTreeHelper(bst.Root, "", true)
}

func (bst *BinarySearchTree) printTreeHelper(node *TreeNode, prefix string, isTail bool) {
	if node == nil {
		return
	}

	connector := "├── "
	if isTail {
		connector = "└── "
	}
	fmt.Println(prefix + connector + fmt.Sprintf("%d", node.Value))

	childPrefix := prefix
	if isTail {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	if node.Left != nil || node.Right != nil {
		if node.Right != nil {
			bst.printTreeHelper(node.Right, childPrefix, false)
		}
		if node.Left != nil {
			bst.printTreeHelper(node.Left, childPrefix, true)
		}
	}
}

// BSTExample 二叉搜索树示例
func BSTExample() {
	fmt.Println("=== 二叉搜索树 (Binary Search Tree) 示例 ===")

	bst := NewBinarySearchTree()

	// 插入节点
	fmt.Println("\n1. 插入节点:")
	values := []int{50, 30, 70, 20, 40, 60, 80, 10, 25}
	fmt.Printf("插入顺序: %v\n", values)
	for _, v := range values {
		bst.Insert(v)
	}
	bst.PrintTree()

	// 查找节点
	fmt.Println("\n2. 查找节点:")
	searchValues := []int{40, 35, 80}
	for _, v := range searchValues {
		found := bst.Contains(v)
		fmt.Printf("查找 %d: %t\n", v, found)
	}

	// 遍历
	fmt.Println("\n3. 树的遍历:")
	fmt.Printf("中序遍历 (升序): %v\n", bst.InorderTraversal())
	fmt.Printf("前序遍历: %v\n", bst.PreorderTraversal())
	fmt.Printf("后序遍历: %v\n", bst.PostorderTraversal())
	fmt.Printf("层序遍历: %v\n", bst.LevelOrderTraversal())

	// 最小值和最大值
	fmt.Println("\n4. 最值查找:")
	fmt.Printf("最小值: %d\n", bst.FindMin().Value)
	fmt.Printf("最大值: %d\n", bst.FindMax().Value)

	// 树的属性
	fmt.Println("\n5. 树的属性:")
	fmt.Printf("节点数: %d\n", bst.Size)
	fmt.Printf("树高度: %d\n", bst.Height())
	fmt.Printf("是否平衡: %t\n", bst.IsBalanced())

	// 删除节点
	fmt.Println("\n6. 删除节点:")
	deleteValues := []int{20, 30, 50}
	for _, v := range deleteValues {
		fmt.Printf("删除 %d:\n", v)
		bst.Delete(v)
		fmt.Printf("中序遍历: %v\n", bst.InorderTraversal())
	}

	bst.PrintTree()
	fmt.Println()
}
