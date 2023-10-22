package pw_test

import (
	"log"
	"os"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

var pwInstance *playwright.Playwright
var browser playwright.Browser
var context playwright.BrowserContext
var page playwright.Page
var expect playwright.PlaywrightAssertions
var server *testServer

var DEFAULT_CONTEXT_OPTIONS = playwright.BrowserNewContextOptions{
	AcceptDownloads: playwright.Bool(true),
	HasTouch:        playwright.Bool(true),
}

func TestMain(m *testing.M) {
	BeforeAll()
	code := m.Run()
	AfterAll()
	os.Exit(code)
}

func BeforeAll() {
	var err error
	pwInstance, err = playwright.Run()
	if err != nil {
		log.Fatalf("could not start Playwright: %v", err)
	}
	browser, err = pwInstance.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		log.Fatalf("could not launch: %v", err)
	}
	expect = playwright.NewPlaywrightAssertions(1000)
	server = newTestServer()
}

func AfterAll() {
	server.testServer.Close()
	if err := pwInstance.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

func BeforeEach(t *testing.T, contextOptions ...playwright.BrowserNewContextOptions) {
	if len(contextOptions) == 1 {
		newContextWithOptions(t, contextOptions[0])
		return
	}
	newContextWithOptions(t, DEFAULT_CONTEXT_OPTIONS)
}

func AfterEach(t *testing.T, closeContext ...bool) {
	if len(closeContext) == 0 {
		if err := context.Close(); err != nil {
			t.Errorf("could not close context: %v", err)
		}
	}
	server.AfterEach()
}

func newContextWithOptions(t *testing.T, contextOptions playwright.BrowserNewContextOptions) {
	var err error
	context, err = browser.NewContext(contextOptions)
	require.NoError(t, err)
	page, err = context.NewPage()
	require.NoError(t, err)
}
