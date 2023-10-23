package chromedvtools_test

import (
	customDriver "scraper/browserhandler/driver"
	"scraper/browserhandler/driver/chromedvtools"
	"testing"

	"github.com/chromedp/chromedp"
)

func makeNode(selector string) chromedvtools.ChromeDvToolsNode {
	return chromedvtools.ChromeDvToolsNode{
		Driver:   driver,
		Selector: selector,
	}
}

func getFirstNode(t *testing.T, selector string) customDriver.NodeFinder {
	var node = makeNode(selector)
	child, err := node.First()
	if err != nil {
		t.Fatalf("test has error, First got error: %v", err)
	}
	return child
}

func getAllNode(t *testing.T, selector string) []customDriver.NodeFinder {
	var node = makeNode(selector)
	children, err := node.All()
	if err != nil {
		t.Fatalf("test has error, First got error: %v", err)
	}
	return children
}

func TestTextContent(t *testing.T) {
	BeforeEach(t, server.FORM_PAGE)

	tests := []struct {
		sel string
		by  chromedp.QueryOption
		exp string
	}{
		{"#inner-hidden", chromedp.ByID, "this is hidden"},
		{"#hidden", chromedp.ByID, "hidden"},
	}

	for i, test := range tests {
		child := getFirstNode(t, test.sel)
		if text, err := child.TextContent(); err != nil {
			t.Fatalf("test %d All got error: %v", i, err)
		} else if text != test.exp {
			t.Errorf("test %d expected %q, got: %s", i, test.exp, text)
		}
	}

	AfterEach(t)
}

func TestAll(t *testing.T) {
	BeforeEach(t, server.FORM_PAGE)

	tests := []struct {
		sel string
		by  chromedp.QueryOption
		exp int
	}{
		{"input", chromedp.ByID, 4},
	}

	for i, test := range tests {
		if children := getAllNode(t, test.sel); len(children) != test.exp {
			t.Errorf("test %d expected %q, got: %d", i, test.exp, len(children))
		}
	}

	AfterEach(t)
}
