package chromedvtools

import (
	"context"
	"fmt"
	"log"
	"scraper/browserhandler/driver"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func waitFor(ctx context.Context, eventName string) error {
	ch := make(chan struct{})
	cctx, cancel := context.WithCancel(ctx)
	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			if e.Name == eventName {
				cancel()
				close(ch)
			}
		}
	})
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func navigateAndWaitFor(url string, eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		_, _, _, err := page.Navigate(url).Do(ctx)
		if err != nil {
			return err
		}
		log.Println("Waiting for", eventName)

		return waitFor(ctx, eventName)
	}
}

func (driver *ChromeDvToolsDriver) NavigateAndWaitFor(url string, eventName string) {
	chromedp.Run(driver.Context,
		navigateAndWaitFor(url, eventName),
	)
}

func (driver *ChromeDvToolsDriver) Start(options ...interface{}) {
	fmt.Println("Starting Chromedp Context")

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	var cancelFunc context.CancelFunc
	driver.ExecContext, cancelFunc = chromedp.NewExecAllocator(context.Background(), opts...)
	driver.CancelFuncs = append(driver.CancelFuncs, cancelFunc)
	driver.Context, cancelFunc = chromedp.NewContext(driver.ExecContext)
	driver.CancelFuncs = append(driver.CancelFuncs, cancelFunc)
}

func (driver *ChromeDvToolsDriver) Run() {
	fmt.Println("Running test")

	var err error

	if err != nil {
		log.Fatal(err)
	}
}

func (driver *ChromeDvToolsDriver) Back() {

}

func (driver *ChromeDvToolsDriver) Navigate(url string) {
	err := chromedp.Run(driver.Context,
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (driver *ChromeDvToolsDriver) NodeFinder(selector string) (driver.NodeFinder, error) {
	var cdpNodes []*cdp.Node
	err := chromedp.Run(driver.Context,
		chromedp.Nodes(selector, &cdpNodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return nil, err
	}

	return &ChromeDvToolsNode{
		Driver:   *driver,
		Selector: selector,
	}, nil
}

func (driver *ChromeDvToolsDriver) Fetch() {
	var elems []*cdp.Node
	var example string
	err := chromedp.Run(driver.Context,
		chromedp.Nodes(".athing", &elems, chromedp.ByQueryAll),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for index, node := range elems {
				id, err := dom.QuerySelector(node.NodeID, "td.title > span > a").Do(ctx)
				if err != nil {
					log.Fatal(err)
				}
				chromedp.TextContent([]cdp.NodeID{id}, &example, chromedp.ByNodeID).Do(ctx)
				fmt.Printf("%d: %s\n", index+1, example)
			}
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}
}

func (driver *ChromeDvToolsDriver) Forward() {

}

func (driver *ChromeDvToolsDriver) Cancel() {
	for _, cancelFunc := range driver.CancelFuncs {
		cancelFunc()
	}
}

func (driver *ChromeDvToolsDriver) Close() {
	for _, cancelFunc := range driver.CancelFuncs {
		cancelFunc()
	}
}
