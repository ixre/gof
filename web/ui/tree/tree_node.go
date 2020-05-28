package tree

type FlatNode struct {
	// 子节点编号
	ID int `json:"-"`
	// 父节点编号
	Pid int `json:"-"`
	// 文本
	Title string `json:"title"`
	// 值
	Value string `json:"value"`
	// 地址
	Url string `json:"url"`
	// 图标
	Icon string `json:"icon"`
	// 是否展开
	Expand bool `json:"expand"`
	// 延迟加载
	Lazy bool `json:"lazy"`
	// 其他数据
	Data map[string]string `json:"data"`
}

type TreeNode struct {
	// 子节点编号
	id int `json:"-"`
	// 文本
	Title string `json:"title"`
	// 值
	Value string `json:"value"`
	// 地址
	Url string `json:"url"`
	// 图标,icon与JS树形控件冲突
	Icon string `json:"icon"`
	// 是否展开
	Expand bool `json:"expanded"`
	// 是否目录，通常Children有元素,则为true
	Folder bool `json:"folder"`
	// 延迟加载
	Lazy bool `json:"lazy"`
	// 其他数据
	Data map[string]string `json:"data"`
	// 子节点
	Children []*TreeNode `json:"children"`
}

func (f FlatNode) Node() *TreeNode {
	return &TreeNode{
		id:       f.ID,
		Title:    f.Title,
		Value:    f.Value,
		Url:      f.Url,
		Icon:     f.Icon,
		Expand:   f.Expand,
		Lazy:     f.Lazy,
		Data:     f.Data,
		Children: []*TreeNode{},
	}
}
