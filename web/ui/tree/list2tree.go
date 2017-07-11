package tree

import (
	_ "fmt"
)

func List2Tree(nodeList []TreeNode) (rootNode *TreeNode) {
	for i, k := range nodeList {
		if k.Id == 0 {
			rootNode = &k
			nodeList = append(nodeList[:i], nodeList[i+1:]...)
			break
		}
	}

	if rootNode == nil {
		rootNode = &TreeNode{
			Id:       0,
			Pid:      0,
			Text:     "根节点",
			Value:    "",
			Url:      "",
			Icon:     "",
			Open:     true,
			Children: nil}
	}
	iterTree(rootNode, nodeList)
	return rootNode
}
func iterTree(node *TreeNode, nodeList []TreeNode) {
	node.Children = []*TreeNode{}
	for _, _cnode := range nodeList {
		cnode := _cnode //必须要新建变量，否则都会引用到最后一个元素
		if cnode.Pid == node.Id {
			node.Children = append(node.Children, &cnode)
			iterTree(&cnode, nodeList)
		}
	}
}
