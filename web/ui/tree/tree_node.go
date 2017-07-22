package tree

type TreeNode struct {
	//子节点编号
	ID int64 `json:"id"`
	//父节点编号
	Pid      int64       `json:"pid"`
	Text     string      `json:"text"`
	Value    string      `json:"value"`
	Url      string      `json:"url"`
	Icon     string      `json:"icon"`
	Open     bool        `json:"open"`
	Children []*TreeNode `json:"childs"`
}
