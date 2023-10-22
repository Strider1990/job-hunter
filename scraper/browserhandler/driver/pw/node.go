package pw

import (
	"scraper/browserhandler/driver"

	"github.com/playwright-community/playwright-go"
)

func makeNode(node *PwNode, locator playwright.Locator, selector *string) *PwNode {
	var newSelector string
	if selector != nil {
		newSelector = *selector
	} else {
		newSelector = node.Selector
	}
	return &PwNode{
		Driver:   node.Driver,
		Selector: newSelector,
		Node:     &locator,
	}
}

func (node *PwNode) NodeFinder(selector string) (driver.NodeFinder, error) {
	locator := getActivePage(&node.Driver).Locator(selector)

	return makeNode(node, locator, &selector), nil
}

func (node *PwNode) All() ([]driver.NodeFinder, error) {
	entries, err := (*node.Node).All()

	if err != nil {
		return nil, err
	}

	var output = make([]driver.NodeFinder, 0)
	for _, entry := range entries {
		output = append(output, makeNode(node, entry, nil))
	}
	return output, nil
}

func (node *PwNode) First() (driver.NodeFinder, error) {
	entry := (*node.Node).First()
	return makeNode(node, entry, nil), nil
}

func (node *PwNode) TextContent() (string, error) {
	return (*node.Node).TextContent()
}
