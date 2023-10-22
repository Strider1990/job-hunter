package browserhandler

import "scraper/browserhandler/driver"

type DriverTypes string

const (
	Chromedp   DriverTypes = "chromedp"
	Playwright DriverTypes = "playwright"
)

type BrowserDriver interface {
	Start(options ...interface{})
	Navigate(url string)
	NavigateAndWaitFor(url string, eventName string)
	NodeFinder(selector string) (driver.NodeFinder, error)
	Back()
	Forward()
	Cancel()
	Close()
}
