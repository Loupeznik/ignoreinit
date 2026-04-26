package src

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type fakeGitignoreClient struct {
	listTemplates   []string
	listErrs        []error
	downloadContent []byte
	downloadErrs    []error
	listCalls       int
	downloadCalls   int
	downloadedPath  string
}

func TestHandleParamsDefaultsLocation(t *testing.T) {
	language, location, err := handleParams([]string{"go"})
	if err != nil {
		t.Fatalf("handleParams returned error: %v", err)
	}

	if language != "go" || location != "." {
		t.Fatalf("handleParams() = %q, %q; want go, .", language, location)
	}
}

func TestHandleParamsRequiresLanguage(t *testing.T) {
	_, _, err := handleParams(nil)
	if err == nil {
		t.Fatal("handleParams() error = nil; want an error")
	}
}

func TestFindTemplateIsCaseInsensitive(t *testing.T) {
	template := findTemplate("go", []string{"Global/Linux.gitignore", "Go.gitignore"})
	if template != "Go.gitignore" {
		t.Fatalf("findTemplate() = %q; want Go.gitignore", template)
	}
}

func TestWriteIgnoreCreatesFileWithContentMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	if err := writeIgnore(path, []byte("bin/\n"), true, false); err != nil {
		t.Fatalf("writeIgnore() returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if !strings.HasPrefix(string(content), generatedHeader) {
		t.Fatalf("content = %q; want generated header", string(content))
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() returned error: %v", err)
	}

	if mode := info.Mode().Perm(); mode != gitignoreFileMode {
		t.Fatalf("file mode = %v; want %v", mode, os.FileMode(gitignoreFileMode))
	}
}

func TestWriteIgnoreMergesWithBlankLineSeparator(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte("bin/"), gitignoreFileMode); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	if err := writeIgnore(path, []byte("dist/\n"), false, true); err != nil {
		t.Fatalf("writeIgnore() returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if got := string(content); got != "bin/\n\ndist/\n" {
		t.Fatalf("merged content = %q; want blank-line separated content", got)
	}
}

func TestWriteIgnoreWrapsWriteErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatalf("Mkdir() returned error: %v", err)
	}

	err := writeIgnore(path, []byte("bin/\n"), false, false)
	if err == nil || !strings.Contains(err.Error(), "could not write") {
		t.Fatalf("writeIgnore() error = %v; want wrapped write error", err)
	}
}

func TestFetchIgnoreRetriesTemplateList(t *testing.T) {
	oldRetryDelay := retryDelay
	retryDelay = 0
	defer func() { retryDelay = oldRetryDelay }()

	client := &fakeGitignoreClient{
		listTemplates:   []string{"Go.gitignore"},
		listErrs:        []error{errors.New("temporary failure")},
		downloadContent: []byte("bin/\n"),
	}

	content, err := fetchIgnore("go", client)
	if err != nil {
		t.Fatalf("fetchIgnore() returned error: %v", err)
	}

	if string(content) != "bin/\n" {
		t.Fatalf("fetchIgnore() = %q; want bin/", string(content))
	}

	if client.listCalls != 2 {
		t.Fatalf("list calls = %d; want 2", client.listCalls)
	}
}

func TestFetchIgnoreReturnsActionableMissingTemplateError(t *testing.T) {
	oldRetryDelay := retryDelay
	retryDelay = 0
	defer func() { retryDelay = oldRetryDelay }()

	client := &fakeGitignoreClient{listTemplates: []string{"Go.gitignore"}}

	_, err := fetchIgnore("rust", client)
	if err == nil || !strings.Contains(err.Error(), "check the language name") {
		t.Fatalf("fetchIgnore() error = %v; want actionable missing template error", err)
	}
}

func TestFetchIgnoreRetriesDownload(t *testing.T) {
	oldRetryDelay := retryDelay
	retryDelay = 0
	defer func() { retryDelay = oldRetryDelay }()

	client := &fakeGitignoreClient{
		listTemplates:   []string{"Go.gitignore"},
		downloadContent: []byte("bin/\n"),
		downloadErrs:    []error{errors.New("temporary failure")},
	}

	_, err := fetchIgnore("go", client)
	if err != nil {
		t.Fatalf("fetchIgnore() returned error: %v", err)
	}

	if client.downloadCalls != 2 {
		t.Fatalf("download calls = %d; want 2", client.downloadCalls)
	}

	if client.downloadedPath != "Go.gitignore" {
		t.Fatalf("download path = %q; want Go.gitignore", client.downloadedPath)
	}
}

func TestWithRetryStopsOnContextCancellation(t *testing.T) {
	oldRetryDelay := retryDelay
	retryDelay = time.Hour
	defer func() { retryDelay = oldRetryDelay }()

	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := withRetry(ctx, func() error {
		calls++
		cancel()
		return errors.New("temporary failure")
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("withRetry() error = %v; want context canceled", err)
	}

	if calls != 1 {
		t.Fatalf("calls = %d; want 1", calls)
	}
}

func (c *fakeGitignoreClient) ListTemplates(ctx context.Context) ([]string, error) {
	c.listCalls++
	if err := nextErr(&c.listErrs); err != nil {
		return nil, err
	}

	return c.listTemplates, nil
}

func (c *fakeGitignoreClient) DownloadTemplate(ctx context.Context, templatePath string) ([]byte, error) {
	c.downloadCalls++
	c.downloadedPath = templatePath
	if err := nextErr(&c.downloadErrs); err != nil {
		return nil, err
	}

	return c.downloadContent, nil
}

func nextErr(errs *[]error) error {
	if len(*errs) == 0 {
		return nil
	}

	err := (*errs)[0]
	*errs = (*errs)[1:]
	return err
}
