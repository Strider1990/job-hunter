package chromedvtools

import (
	"errors"
	"fmt"
	"scraper/browserhandler/driver"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func (node *ChromeDvToolsNode) TextContent() (string, error) {
	var output string
	err := chromedp.Run(
		node.Driver.Context,
		chromedp.TextContent([]cdp.NodeID{node.Node.NodeID}, &output, chromedp.ByNodeID),
	)
	return output, err
}

func (node *ChromeDvToolsNode) NodeFinder(selector string) (driver.NodeFinder, error) {
	node.Selector = selector
	return node, nil
}

func (node *ChromeDvToolsNode) All() ([]driver.NodeFinder, error) {
	var cdpNodes []*cdp.Node
	fmt.Println("Running ALL")
	err := chromedp.Run(node.Driver.Context,
		chromedp.Nodes(node.Selector, &cdpNodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, err
	}
	var output = make([]driver.NodeFinder, 0)
	for _, cdpNode := range cdpNodes {
		output = append(output, &ChromeDvToolsNode{
			Driver:   node.Driver,
			Selector: node.Selector,
			NodeID:   cdpNode.NodeID.Int64(),
			Node:     cdpNode,
		})
	}
	return output, nil
}

func (node *ChromeDvToolsNode) First() (driver.NodeFinder, error) {
	var cdpNodes []*cdp.Node
	err := chromedp.Run(node.Driver.Context,
		chromedp.Nodes(node.Selector, &cdpNodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, err
	}
	if len(cdpNodes) > 0 {
		return &ChromeDvToolsNode{
			Driver:   node.Driver,
			Selector: node.Selector,
			NodeID:   cdpNodes[0].NodeID.Int64(),
			Node:     cdpNodes[0],
		}, nil
	}
	return nil, errors.New("No nodes found by that selector")
}
