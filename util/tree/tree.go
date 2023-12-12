package tree

// 节点数据
type NodeData struct {
	// 子节点编号
	Id string `json:"id"`
	// 文本
	Label string `json:"label"`
	// 父节点编号
	ParentId string `json:"parentId"`
	// 是否叶子节点
	IsLeaf bool `json:"isLeaf"`
	// 其他数据
	Data interface{} `json:"data"`
}

// 树形节点
type TreeNode struct {
	// 标识
	Id string `json:"id"`
	// 文本
	Label string `json:"label"`
	// 是否叶子节点
	IsLeaf bool `json:"isLeaf"`
	// 子节点
	Children []*TreeNode `json:"children"`
	// 其他数据
	Data interface{} `json:"data"`
}

func (f NodeData) Node() *TreeNode {
	return &TreeNode{
		Id:       f.Id,
		Label:    f.Label,
		IsLeaf:   f.IsLeaf,
		Data:     f.Data,
		Children: []*TreeNode{},
	}
}

// 转换为树形
func ParseTree(nodeList []NodeData, nodeFn func(node *TreeNode)) []*TreeNode {
	root := &TreeNode{}
	for i, k := range nodeList {
		if k.Id == "" {
			root = k.Node()
			nodeList = append(nodeList[:i], nodeList[i+1:]...)
			break
		}
	}
	walkTree(root, nodeList, nodeFn)
	return root.Children
}

func walkTree(node *TreeNode, nodeList []NodeData, nodeFn func(node *TreeNode)) {
	node.Children = []*TreeNode{}
	for _, v := range nodeList {
		if v.ParentId == node.Id {
			n := v.Node()
			if nodeFn != nil {
				nodeFn(n)
			}
			node.Children = append(node.Children, n)
			walkTree(n, nodeList, nodeFn)
		}
	}
}
