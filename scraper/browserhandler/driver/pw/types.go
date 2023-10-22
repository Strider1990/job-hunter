package pw

import (
	"github.com/playwright-community/playwright-go"
)

type PwDriver struct {
	Instance      *playwright.Playwright
	Browser       playwright.Browser
	Pages         []playwright.Page
	ActivePageID  int
	SelectedNodes []*playwright.Locator
}

type PwNode struct {
	Driver   PwDriver
	Selector string
	Node     *playwright.Locator
}
