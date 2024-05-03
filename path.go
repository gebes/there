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

func methodToInt(method string) method {
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

func methodToString(method method) string {
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
	isDynamic := path != "" && path[0] == ':'
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

	if path == "" {
		return n, params
	}

	// skip initial slash
	if path[0] == '/' {
		if len(path) == 1 {
			path = ""
		} else {
			pathIndex++
		}
	}

	// current segment of the path /segment1/segment2/segment3
	var segment string

	for {
		// check if we reached the end of the route
		if len(path)-pathIndex == 0 {
			return n, params
		}

		// get next path segment
		// optimized way of split(path, "/")[index]
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

		// check if we got a child
		next, ok := n.children[segment]
		if !ok {
			// if not check, if a wildcard child exists
			if n.paramNode == nil {
				return nil, nil
			}
			// if so, initialize the params map now
			// lazy initialization to save some time
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

	if path == "" {
		return n, nil
	}

	if path[0] == '/' {
		pathIndex++
	}

	// current segment of the path /segment1/segment2/segment3
	var segment string

	for {
		// check if we reached the end of the route
		if len(path)-pathIndex == 0 {
			return n, nil
		}

		// get next path segment
		// optimized way of split(path, "/")[index]
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

		// check if we got a child
		next, ok := n.children[segment]
		if !ok {
			isParamSegment := segment != "" && segment[0] == ':'
			if isParamSegment {
				// create node with current name.
				// name needs to be stored later with request.RouteParams
				stripped := strings.TrimPrefix(segment, ":")
				if n.paramNode == nil {
					n.paramNode = &node{name: stripped}
				}

				// cover edge case when user tries to have the following routes
				//  /user/:id/test
				//  /user/:name/test2
				// the data structure does not comply with that, since we only
				// store the name of one wildcard child
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
