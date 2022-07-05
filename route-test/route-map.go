package main

import (
	"strings"

	"github.com/Gebes/there/v2"
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
	case there.MethodHead:
		return methodHead
	case there.MethodPost:
		return methodPost
	case there.MethodPut:
		return methodPut
	case there.MethodPatch:
		return methodPatch
	case there.MethodDelete:
		return methodDelete
	case there.MethodConnect:
		return methodConnect
	case there.MethodOptions:
		return methodOptions
	case there.MethodTrace:
		return methodTrace
	default:
		return methodGet
	}
}

type node struct {
	name         string
	regularNodes map[string]*node
	paramNode    *node
	wildcardNode *node
	handler      [methods]there.Endpoint
	middlewares  []there.Middleware
}

type handlerResult struct {
	handler     *[methods]there.Endpoint
	middlewares []there.Middleware
	params      map[string]string
	wildcard    string
}

func (n *node) FindHandler(parts []string) *handlerResult {
	return findHandler(n, &parts, 0)
}

func findHandler(currentNode *node, parts *[]string, index int) *handlerResult {
	if currentNode == nil {
		return nil
	}
	if len(*parts)-index == 0 {
		return &handlerResult{
			handler:     &currentNode.handler,
			middlewares: currentNode.middlewares,
			params:      map[string]string{},
		}
	}
	part := (*parts)[index]
	index++

	var result *handlerResult

	findHandler(currentNode.regularNodes[part], parts, index)

	if result == nil {
		result = findHandler(currentNode.paramNode, parts, index)
		result.params[currentNode.paramNode.name] = part
	}

	if result == nil {
		result = findHandler(currentNode.wildcardNode, parts, index)
		if result != nil {
			result.wildcard = strings.Join((*parts)[index:], "/")
		}
	}

	return result
}

func (n *node) AddHandler(method method, path string, endpoint there.Endpoint) {

}

func addHandler(currentNode *node) {}
