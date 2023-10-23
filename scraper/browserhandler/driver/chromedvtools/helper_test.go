package chromedvtools_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"scraper/browserhandler"
	"scraper/browserhandler/driver/chromedvtools"
	"sync"
	"testing"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

var (
	// these are set up in init
	execPath  string
	allocOpts = chromedp.DefaultExecAllocatorOptions[:]

	// allocCtx is initialised in TestMain, to cancel before exiting.
	allocCtx context.Context

	// browserCtx is initialised with allocateOnce
	browserCtx   context.Context
	allocateOnce sync.Once
	allocTempDir string
	browserOpts  []chromedp.ContextOption
	server       *browserhandler.TestServer
)

func BeforeAll() {

}

func AfterAll() {
	server.TestServer.Close()
}

func init() {
	server = browserhandler.NewTestServer("../tests/assets")
	var err error

	allocTempDir, err = os.MkdirTemp("", "chromedp-test")
	if err != nil {
		panic(fmt.Sprintf("could not create temp directory: %v", err))
	}

	// Disabling the GPU helps portability with some systems like Travis,
	// and can slightly speed up the tests on other systems.
	allocOpts = append(allocOpts, chromedp.DisableGPU)

	if noHeadless := os.Getenv("CHROMEDP_NO_HEADLESS"); noHeadless != "" && noHeadless != "false" {
		allocOpts = append(allocOpts, chromedp.Flag("headless", false))
	}

	// Find the exec path once at startup.
	execPath = os.Getenv("CHROMEDP_TEST_RUNNER")
	allocOpts = append(allocOpts, chromedp.ExecPath(execPath))

	// Not explicitly needed to be set, as this speeds up the tests
	if noSandbox := os.Getenv("CHROMEDP_NO_SANDBOX"); noSandbox != "false" {
		allocOpts = append(allocOpts, chromedp.NoSandbox)
	}
}

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	allocCtx, cancel = chromedp.NewExecAllocator(context.Background(), allocOpts...)

	if debug := os.Getenv("CHROMEDP_DEBUG"); debug != "" && debug != "false" {
		browserOpts = append(browserOpts, chromedp.WithDebugf(log.Printf))
	}

	code := m.Run()
	cancel()

	if infos, _ := os.ReadDir(allocTempDir); len(infos) > 0 {
		os.RemoveAll(allocTempDir)
		panic(fmt.Sprintf("leaked %d temporary dirs under %s", len(infos), allocTempDir))
	} else {
		os.Remove(allocTempDir)
	}

	os.Exit(code)
}

func testAllocate(tb testing.TB, url string) (context.Context, context.CancelFunc) {
	// Start the browser exactly once, as needed.
	allocateOnce.Do(func() { browserCtx, _ = testAllocateSeparate(tb) })

	if browserCtx == nil {
		// allocateOnce.Do failed; continuing would result in panics.
		tb.FailNow()
	}

	// Same browser, new tab; not needing to start new chrome browsers for
	// each test gives a huge speed-up.
	ctx, _ := chromedp.NewContext(browserCtx)

	// Only navigate if we want an HTML file name, otherwise leave the blank page.
	if name != "" {
		if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
			tb.Fatal(err)
		}
	}

	cancel := func() {
		if err := chromedp.Cancel(ctx); err != nil {
			tb.Error(err)
		}
	}
	return ctx, cancel
}

func testAllocateSeparate(tb testing.TB) (context.Context, context.CancelFunc) {
	// Entirely new browser, unlike testAllocate.
	ctx, _ := chromedp.NewContext(allocCtx, browserOpts...)
	if err := chromedp.Run(ctx); err != nil {
		tb.Fatal(err)
	}
	chromedp.ListenBrowser(ctx, func(ev interface{}) {
		if ev, ok := ev.(*runtime.EventExceptionThrown); ok {
			tb.Errorf("%+v\n", ev.ExceptionDetails)
		}
	})
	cancel := func() {
		if err := chromedp.Cancel(ctx); err != nil {
			tb.Error(err)
		}
	}
	return ctx, cancel
}

// ===== Setup driver helper =====

var driver chromedvtools.ChromeDvToolsDriver
var cancel context.CancelFunc

func BeforeEach(t *testing.T, url string) {
	t.Parallel()

	var ctx context.Context
	ctx, cancel = testAllocate(t, url)
	driver = chromedvtools.ChromeDvToolsDriver{
		ExecContext: allocCtx,
		Context:     ctx,
		CancelFuncs: []context.CancelFunc{cancel},
	}
}

func AfterEach(t *testing.T) {
	cancel()
}
