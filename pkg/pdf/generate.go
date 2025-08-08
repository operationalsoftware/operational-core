package pdf

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type PDFDefinition struct {
	HTML  string
	Title string
}

var (
	browserCtx    context.Context
	browserCancel context.CancelFunc
	once          sync.Once
)

// InitChromium initializes the shared Chrome allocator
func InitChromium() {
	once.Do(func() {
		browserCtx, browserCancel = chromedp.NewExecAllocator(
			context.Background(),
			append(chromedp.DefaultExecAllocatorOptions[:],
				chromedp.Flag("disable-gpu", true),
				chromedp.Flag("no-sandbox", true),
				chromedp.Flag("disable-setuid-sandbox", true),
				chromedp.Flag("disable-dev-shm-usage", true),
			)...,
		)
	})
}

// ShutdownChromium terminates the shared Chrome context
func ShutdownChromium() {
	if browserCancel != nil {
		browserCancel()
	}
}

func GeneratePDF(reqCtx context.Context, pdfDef PDFDefinition) ([]byte, error) {
	if browserCtx == nil {
		return nil, fmt.Errorf("chrome not initialized, call InitChrome() first")
	}

	// Derive a new context with timeout
	ctx, cancel := context.WithTimeout(reqCtx, 30*time.Second)
	defer cancel()

	// Create a new tab context derived from the timed request context
	tabCtx, cancelTab := chromedp.NewContext(ctx)
	defer cancelTab()

	var pdfBuf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate("data:text/html," + htmlEscape(pdfDef.HTML)),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return waitUntilHTMLRendered(ctx, 10*time.Second)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				Do(ctx)
			return err
		}),
	}

	if err := chromedp.Run(tabCtx, tasks); err != nil {
		return nil, fmt.Errorf("error running chromedp tasks: %w", err)
	}

	pdfBuf, err := setPDFTitle(pdfBuf, pdfDef.Title)
	if err != nil {
		return nil, fmt.Errorf("error setting PDF title: %w", err)
	}

	return pdfBuf, nil
}

// Helper to escape HTML for URL
func htmlEscape(html string) string {
	return url.PathEscape(html)
}

func waitUntilHTMLRendered(ctx context.Context, timeout time.Duration) error {
	const checkDuration = 50 * time.Millisecond
	const minStableIterations = 4

	var lastSize int
	stableCount := 0
	maxChecks := int(timeout / checkDuration)

	for range maxChecks {
		var htmlSize int
		err := chromedp.Run(ctx, chromedp.Evaluate(`document.documentElement.outerHTML.length`, &htmlSize))
		if err != nil {
			return err
		}

		if lastSize != 0 && htmlSize == lastSize {
			stableCount++
		} else {
			stableCount = 0
		}

		if stableCount >= minStableIterations {
			return nil
		}

		lastSize = htmlSize
		time.Sleep(checkDuration)
	}

	return fmt.Errorf("HTML did not stabilize within %v", timeout)
}

func setPDFTitle(pdfData []byte, title string) ([]byte, error) {
	inputFile, err := os.CreateTemp("", "input-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(inputFile.Name())

	if _, err := inputFile.Write(pdfData); err != nil {
		return nil, err
	}
	if err := inputFile.Close(); err != nil {
		return nil, err
	}

	outputFile, err := os.CreateTemp("", "output-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close()

	// Set metadata using AddProperties
	props := map[string]string{
		"Title": title,
	}

	if err := api.AddPropertiesFile(inputFile.Name(), outputFile.Name(), props, nil); err != nil {
		return nil, err
	}

	return os.ReadFile(outputFile.Name())
}
