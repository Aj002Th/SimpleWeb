package sim

import "strings"

/*
*	something not good:
*		没有实现路由模糊匹配和精确匹配之间的优先级
 */

type node struct {
	pattern  string // 记录以这个节点结束的路由的完整路径
	part     string
	children []*node
	isParam  bool
}

// 找到一个匹配的子节点,未匹配返回nil
// 用于插入
func (n *node) matchChild(part string) *node {
	for _, c := range n.children {
		if c.part == part || c.isParam {
			return c
		}
	}
	return nil
}

// 找到所有匹配的子节点,未匹配返回长度为0的切片
// 用于查找
func (n *node) matchChildren(part string) []*node {
	var ret []*node
	for _, c := range n.children {
		if c.part == part || c.isParam {
			ret = append(ret, c)
		}
	}
	return ret
}

// 插入路由规则
// 插入相同的规则会直接进行覆盖，不返回err
func (n *node) insert(pattern string, parts []string, height int) {
	// 结束条件
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 定位或创建子节点
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:     part,
			children: []*node{},
			isParam:  part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	// 递归
	child.insert(pattern, parts, height+1)
}

// 匹配路由规则
func (n *node) search(parts []string, height int) *node {
	// 结束条件
	// 注意遇到*都是能够匹配的
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, c := range children {
		ret := c.search(parts, height+1)
		if ret != nil {
			return ret
		}
	}
	return nil
}
