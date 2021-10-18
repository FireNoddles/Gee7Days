package gee

type node struct {
	pattern  string  //叶子节点的完整路径？
	part     string  //当前节点
	children []*node //当前节点的子节点
	isWild   bool    //是否动态节点
}

//用于insert
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//用于search
func (n *node) matchChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}

func (n *node) insert(pattern string, parts []string, height int) {
	//出口 到头了
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	//查找匹配子节点
	part := parts[height]
	child := n.matchChild(part)
	//没有匹配的孩子节点 需要建一个 并放到当前节点的孩子节点集合中
	if child == nil {
		child := &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)

}

//找叶子节点
func (n *node) search(parts []string, height int) *node {
	//出口 找到底了 返回结点
	if len(parts) == height || n.part[0] == '*' {
		if n.pattern == "" {
			return nil //不懂这个if是干嘛的
		}
		return n
	}

}
