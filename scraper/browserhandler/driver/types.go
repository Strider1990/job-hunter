package driver

type NodeFinder interface {
	NodeFinder(selector string) (NodeFinder, error)
	All() ([]NodeFinder, error)
	First() (NodeFinder, error)
	TextContent() (string, error)
}
