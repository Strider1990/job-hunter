package chromedvtools

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func getTextAction(selector string, output *string, node *cdp.Node) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		var nodeText string
		chromedp.TextContent([]cdp.NodeID{node.NodeID}, &nodeText, chromedp.ByNodeID).Do(ctx)
		*output = nodeText
		return nil
	}
}
