package tree

// 节点数据
type NodeData struct {
	// 子节点编号
	Id string `json:"-"`
	// 文本
	Label string `json:"name"`
	// 父节点编号
	Parent string `json:"-"`
	// 是否叶子节点
	IsLeaf bool `json:"expand"`
	// 其他数据
	Data interface{} `json:"data"`
}

// 树形节点
type TreeNode struct {
	// 标识
	id string `json:"id"`
	// 文本
	Label string `json:"name"`
	// 延迟加载
	//Lazy bool `json:"lazy"`
	// 是否叶子节点
	IsLeaf bool `json:"isLeaf"`
	// 子节点
	Children []*TreeNode `json:"children"`
	// 其他数据
	Data interface{} `json:"data"`
}

func (f NodeData) Node() *TreeNode {
	return &TreeNode{
		id:       f.Id,
		Label:     f.Label,
		IsLeaf:   f.IsLeaf,
		Data:     f.Data,
		Children: []*TreeNode{},
	}
}

// 转换为树形
func ParseTree(nodeList []NodeData, nodeFn func(node *TreeNode)) (rootNode *TreeNode) {
	for i, k := range nodeList {
		if k.Id == "" {
			rootNode = k.Node()
			nodeList = append(nodeList[:i], nodeList[i+1:]...)
			break
		}
	}
	if rootNode == nil {
		rootNode = &TreeNode{
			id:       "",
			Label:     "根节点",
			Children: nil}
	}
	walkTree(rootNode, nodeList, nodeFn)
	return rootNode
}

func walkTree(node *TreeNode, nodeList []NodeData, nodeFn func(node *TreeNode)) {
	node.Children = []*TreeNode{}
	for _, v := range nodeList {
		if v.Parent == node.id {
			n := v.Node()
			if nodeFn != nil {
				nodeFn(n)
			}
			node.Children = append(node.Children, n)
			walkTree(n, nodeList, nodeFn)
		}
	}
}
