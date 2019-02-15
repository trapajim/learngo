package main

import (
	"net/http"
	"net/url"
	"strings"
)

// node represents a node in the tree
type node struct {
	children     []*node
	component    string
	isNamedParam bool
	methods      map[string]Handle
}

// Handle is like net/http handlers but takes params
type Handle func(http.ResponseWriter, *http.Request, url.Values)

// Router our router
type Router struct {
	tree        *node
	rootHandler Handle
}

// Method represents http methods
type Method string

// New creates a new Router and automatically adds a Root Route
func New(rootHandler Handle) *Router {
	node := node{component: "/", isNamedParam: false, methods: make(map[string]Handle)}
	return &Router{tree: &node, rootHandler: rootHandler}
}

// Handle takes an http method, route pattern and handler for the route
func (r *Router) Handle(method, path string, handler Handle) {
	if path[0] != '/' {
		panic("Path has to start with a /")
	}
	r.tree.addNode(method, path, handler)
}

// GET is a shortcut to for router.Handle("GET", path, handler)
func (r *Router) GET(path string, handler Handle) {
	r.Handle("GET", path, handler)
}

// POST is a shortcut to for router.Handle("POST", path, handler)
func (r *Router) POST(path string, handler Handle) {
	r.Handle("POST", path, handler)
}

// PUT is a shortcut to for router.Handle("PUT", path, handler)
func (r *Router) PUT(path string, handler Handle) {
	r.Handle("PUT", path, handler)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) PATCH(path string, handle Handle) {
	r.Handle("PATCH", path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) DELETE(path string, handle Handle) {
	r.Handle("DELETE", path, handle)
}

// implemenation of net/http
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	params := req.Form
	node, _ := r.tree.traverse(strings.Split(req.URL.Path, "/")[1:], params)
	if handler := node.methods[req.Method]; handler != nil {
		handler(w, req, params)
	} else {
		r.rootHandler(w, req, params)
	}
}

// add node ads a node to the tree, it will add multiple components if the
// path can be broken up. Those nodes won't have any handler
// and fall through to the default
func (n *node) addNode(method, path string, handler Handle) {
	components := strings.Split(path, "/")[1:]
	count := len(components)

	for {
		aNode, component := n.traverse(components, nil)
		if aNode.component == component && count == 1 { // update an existing node.
			aNode.methods[method] = handler
			return
		}
		newNode := node{component: component, isNamedParam: false, methods: make(map[string]Handle)}

		if len(component) > 0 && component[0] == ':' { // check if it is a named param.
			newNode.isNamedParam = true
		}
		if count == 1 { // this is the last component of the url resource, so it gets the handler.
			newNode.methods[method] = handler
		}
		aNode.children = append(aNode.children, &newNode)
		count--
		if count == 0 {
			break
		}
	}
}

// traveres follows the tree and adds named parameters
func (n *node) traverse(components []string, params url.Values) (*node, string) {
	component := components[0]
	if len(n.children) > 0 { // no children, then bail out.
		for _, child := range n.children {
			if component == child.component || child.isNamedParam {
				if child.isNamedParam && params != nil {
					params.Add(child.component[1:], component)
				}
				next := components[1:]
				if len(next) > 0 { // http://xkcd.com/1270/
					return child.traverse(next, params) // tail recursion is it's own reward.
				} else {
					return child, component
				}
			}
		}
	}
	return n, component
}
