package pdf

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PdfOptions struct {
	Title string // Optional PDF metadata title
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

// GeneratePDF generates a PDF from HTML with options, using a shared Chrome context
func GeneratePDF(html string, options *PdfOptions) ([]byte, error) {
	if browserCtx == nil {
		return nil, fmt.Errorf("chrome not initialized, call InitChrome() first")
	}

	taskCtx, cancel := chromedp.NewContext(browserCtx)
	defer cancel()

	// Timeout for individual task
	taskCtx, cancel = context.WithTimeout(taskCtx, 30*time.Second)
	defer cancel()

	var pdfBuf []byte
	tasks := chromedp.Tasks{
		chromedp.Navigate("data:text/html," + htmlEscape(html)),
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

	if err := chromedp.Run(taskCtx, tasks); err != nil {
		return nil, err
	}

	if options != nil && options.Title != "" {
		pdfBufWithTitle, err := setPDFTitleSafe(pdfBuf, options.Title)
		if err == nil {
			return pdfBufWithTitle, nil
		}
	}

	return pdfBuf, nil
}

// Helper to escape HTML for URL
func htmlEscape(html string) string {
	return url.PathEscape(html)
}

// setPDFTitleSafe sets metadata title using unique temp files
func setPDFTitleSafe(buf []byte, title string) ([]byte, error) {
	tempDir := os.TempDir()
	randomID := randomHex(8)
	inputPath := filepath.Join(tempDir, fmt.Sprintf("input_%s.pdf", randomID))
	outputPath := filepath.Join(tempDir, fmt.Sprintf("output_%s.pdf", randomID))

	if err := os.WriteFile(inputPath, buf, 0644); err != nil {
		return nil, err
	}

	cmd := exec.Command("exiftool", "-Title="+title, "-o", outputPath, inputPath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	result, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	_ = os.Remove(inputPath)
	_ = os.Remove(outputPath)

	return result, nil
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func waitUntilHTMLRendered(ctx context.Context, timeout time.Duration) error {
	const checkDuration = 50 * time.Millisecond
	const minStableIterations = 4

	var lastSize int
	stableCount := 0
	maxChecks := int(timeout / checkDuration)

	for i := 0; i < maxChecks; i++ {
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
