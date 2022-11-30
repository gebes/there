package there

import (
	"fmt"
	"strings"
)

type method uint8

const (
	methodGet method = iota
	methodHead
	methodPost
	methodPut
	methodPatch
	methodDelete
	methodConnect
	methodOptions
	methodTrace
	methods
)

func methodToInt(method Method) method {
	switch method {
	case MethodHead:
		return methodHead
	case MethodPost:
		return methodPost
	case MethodPut:
		return methodPut
	case MethodPatch:
		return methodPatch
	case MethodDelete:
		return methodDelete
	case MethodConnect:
		return methodConnect
	case MethodOptions:
		return methodOptions
	case MethodTrace:
		return methodTrace
	default:
		return methodGet
	}
}

func methodToString(method method) Method {
	switch method {
	case methodHead:
		return MethodHead
	case methodPost:
		return MethodPost
	case methodPut:
		return MethodPut
	case methodPatch:
		return MethodPatch
	case methodDelete:
		return MethodDelete
	case methodConnect:
		return MethodConnect
	case methodOptions:
		return MethodOptions
	case methodTrace:
		return MethodTrace
	default:
		return MethodGet
	}
}

type (
	matcher struct {
		root   *node
		static map[string]*node
	}

	node struct {
		name        string
		children    map[string]*node
		paramNode   *node
		handler     [methods]Endpoint
		middlewares [methods][]Middleware
	}
)

func (m *matcher) ensureNodeExists(path string) (*node, error) {
	isDynamic := len(path) != 0 && path[0] == ':'
	if !isDynamic {
		for i := range path {
			if path[i] == '/' && i != len(path)-1 && path[i+1] == ':' {
				isDynamic = true
				break
			}
		}
	}
	if isDynamic {
		return m.root.ensureNodeExists(path)
	} else {
		n, ok := m.static[path]
		if ok {
			return n, nil
		}
		n = &node{}
		m.static[path] = n
		return n, nil
	}
}

func (m *matcher) findNode(path string) (*node, map[string][]string) {
	n, ok := m.static[path]
	if ok {
		return n, nil
	}

	pathIndex := 0

	var params map[string][]string

	n = m.root

	if len(path) == 0 {
		return n, params
	}

	if path[0] == '/' {
		if len(path) == 1 {
			path = ""
		} else {
			pathIndex++
		}
	}

	var segment string

	for {
		if len(path)-pathIndex == 0 {
			return n, params
		}

		for i := pathIndex; i < len(path); i++ {
			if path[i] == '/' {
				segment = path[pathIndex:i]
				pathIndex = i + 1
				break
			} else if i == len(path)-1 {
				i++
				segment = path[pathIndex:i]
				pathIndex = i
				break
			}
		}

		next, ok := n.children[segment]
		if !ok {
			if n.paramNode == nil {
				return nil, nil
			}
			if params == nil {
				params = map[string][]string{}
			}
			params[n.paramNode.name] = append(params[n.paramNode.name], segment)
			n = n.paramNode
		} else {
			n = next
		}
	}
}

func (currentNode *node) ensureNodeExists(path string) (*node, error) {
	pathIndex := 0

	n := currentNode

	if len(path) == 0 {
		return n, nil
	}

	if path[0] == '/' {
		pathIndex++
	}

	var segment string

	for {
		if len(path)-pathIndex == 0 {
			return n, nil
		}

		for i := pathIndex; i < len(path); i++ {
			if path[i] == '/' {
				segment = path[pathIndex:i]
				pathIndex = i + 1
				break
			} else if i == len(path)-1 {
				i++
				segment = path[pathIndex:i]
				pathIndex = i
				break
			}
		}
		next, ok := n.children[segment]
		if !ok {
			isParamSegment := len(segment) != 0 && segment[0] == ':'
			if isParamSegment {
				stripped := strings.TrimPrefix(segment, ":")
				if n.paramNode == nil {
					n.paramNode = &node{name: stripped}
				}
				if n.paramNode.name != stripped {
					return nil, fmt.Errorf("path variable \"%s\" for path \"%s\" needs to equal \"%s\" as in all other routes", stripped, path, n.paramNode.name)
				}
				n = n.paramNode
			} else {
				if n.children == nil {
					n.children = map[string]*node{}
				}
				next = &node{name: segment}
				n.children[segment] = next
				n = next
			}
		} else {
			n = next
		}
	}
}
