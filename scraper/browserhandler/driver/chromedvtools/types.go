package chromedvtools

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
)

type ChromeDvToolsDriver struct {
	ExecContext   context.Context
	Context       context.Context
	CancelFuncs   []context.CancelFunc
	SelectedNodes []*cdp.Node
}

type ChromeDvToolsNode struct {
	Driver   ChromeDvToolsDriver
	Selector string
	NodeID   int64
	Node     *cdp.Node
}
