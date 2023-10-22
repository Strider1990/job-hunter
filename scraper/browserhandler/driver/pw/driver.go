package pw

import (
	"log"
	"scraper/browserhandler/driver"

	"github.com/playwright-community/playwright-go"
)

func getActivePage(driver *PwDriver) playwright.Page {
	return driver.Pages[driver.ActivePageID]
}

func (driver *PwDriver) Start(options ...interface{}) {
	var err error
	driver.Instance, err = playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	driver.Browser, err = driver.Instance.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := driver.Browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	driver.Pages = append(driver.Pages, page)
	driver.ActivePageID = len(driver.Pages) - 1
}

func (driver *PwDriver) Navigate(url string) {
	page := getActivePage(driver)
	if _, err := page.Goto(url); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
}

func (driver *PwDriver) NavigateAndWaitFor(url string, eventName string) {
	page := getActivePage(driver)
	if _, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	}); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
}

func (driver *PwDriver) NodeFinder(selector string) (driver.NodeFinder, error) {
	locator := getActivePage(driver).Locator(selector)

	return &PwNode{
		Driver:   *driver,
		Selector: selector,
		Node:     &locator,
	}, nil
}
func (driver *PwDriver) Back()    {}
func (driver *PwDriver) Forward() {}
func (driver *PwDriver) Cancel() {
}
func (driver *PwDriver) Close() {
	if err := driver.Browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err := driver.Instance.Stop(); err != nil {
		log.Fatalf("could not close playwright: %v", err)
	}
}
