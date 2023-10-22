package browserhandler

import (
	"scraper/browserhandler/driver/chromedvtools"
	"scraper/browserhandler/driver/pw"
)

func NewBrowser(driverType DriverTypes) BrowserDriver {
	switch driverType {
	case Chromedp:
		return &chromedvtools.ChromeDvToolsDriver{}
	case Playwright:
		return &pw.PwDriver{}
	}
	return &pw.PwDriver{}
}
