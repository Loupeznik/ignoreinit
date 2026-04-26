package src

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type fakeGitignoreClient struct {
	listTemplates    []string
	listErrs         []error
	downloadContent  []byte
	downloadContents [][]byte
	downloadErrs     []error
	listCalls        int
	downloadCalls    int
	downloadedPath   string
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

func TestHandleGenerationParamsDefaultsLocationForMultipleTemplates(t *testing.T) {
	templates, location, err := handleGenerationParams([]string{"go", "node", "terraform"})
	if err != nil {
		t.Fatalf("handleGenerationParams() returned error: %v", err)
	}

	if got := strings.Join(templates, ", "); got != "go, node, terraform" {
		t.Fatalf("templates = %q; want go, node, terraform", got)
	}

	if location != "." {
		t.Fatalf("location = %q; want .", location)
	}
}

func TestHandleGenerationParamsUsesExplicitLocation(t *testing.T) {
	templates, location, err := handleGenerationParams([]string{"go", "node", "./project"})
	if err != nil {
		t.Fatalf("handleGenerationParams() returned error: %v", err)
	}

	if got := strings.Join(templates, ", "); got != "go, node" {
		t.Fatalf("templates = %q; want go, node", got)
	}

	if location != "./project" {
		t.Fatalf("location = %q; want ./project", location)
	}
}

func TestHandleGenerationParamsIgnoresPrintFlagArgs(t *testing.T) {
	templates, location, err := handleGenerationParams([]string{"go", "node", "--print=true"})
	if err != nil {
		t.Fatalf("handleGenerationParams() returned error: %v", err)
	}

	if got := strings.Join(templates, ", "); got != "go, node" {
		t.Fatalf("templates = %q; want go, node", got)
	}

	if location != "." {
		t.Fatalf("location = %q; want .", location)
	}
}

func TestNormalizeGenerationPrintArgsKeepsPrintFromConsumingTemplates(t *testing.T) {
	args := NormalizeGenerationPrintArgs([]string{"ignoreinit", "init", "--print", "go", "node"})

	if got := strings.Join(args, " "); got != "ignoreinit init --print=true go node" {
		t.Fatalf("NormalizeGenerationPrintArgs() = %q; want print flag normalized", got)
	}
}

func TestNormalizeGenerationPrintArgsKeepsExplicitBoolValues(t *testing.T) {
	args := NormalizeGenerationPrintArgs([]string{"ignoreinit", "init", "--print", "false", "go"})

	if got := strings.Join(args, " "); got != "ignoreinit init --print false go" {
		t.Fatalf("NormalizeGenerationPrintArgs() = %q; want explicit bool value unchanged", got)
	}
}

func TestHandleGenerationParamsDoesNotStatImplicitLocation(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	if err := os.Mkdir("project", 0755); err != nil {
		t.Fatalf("Mkdir() returned error: %v", err)
	}

	templates, location, err := handleGenerationParams([]string{"go", "project"})
	if err != nil {
		t.Fatalf("handleGenerationParams() returned error: %v", err)
	}

	if got := strings.Join(templates, ", "); got != "go, project" {
		t.Fatalf("templates = %q; want go, project", got)
	}

	if location != "." {
		t.Fatalf("location = %q; want .", location)
	}
}

func TestFindTemplateIsCaseInsensitive(t *testing.T) {
	template := findTemplate("go", []string{"Global/Linux.gitignore", "Go.gitignore"})
	if template != "Go.gitignore" {
		t.Fatalf("findTemplate() = %q; want Go.gitignore", template)
	}
}

func TestListTemplateNamesReturnsSortedNames(t *testing.T) {
	oldRetryDelay := retryDelay
	retryDelay = 0
	defer func() { retryDelay = oldRetryDelay }()

	client := &fakeGitignoreClient{listTemplates: []string{"Node.gitignore", "Global/Linux.gitignore", "Go.gitignore"}}

	names, err := listTemplateNames(client)
	if err != nil {
		t.Fatalf("listTemplateNames() returned error: %v", err)
	}

	want := "Go, Linux, Node"
	if got := strings.Join(names, ", "); got != want {
		t.Fatalf("listTemplateNames() = %q; want %q", got, want)
	}
}

func TestSearchNamesReturnsCloseMatches(t *testing.T) {
	names := []string{"Go", "Node", "Terraform", "TeX", "VisualStudioCode"}

	matches := searchNames("terfrm", names)

	if len(matches) == 0 || matches[0] != "Terraform" {
		t.Fatalf("searchNames() = %v; want Terraform as first match", matches)
	}
}

func TestSearchNamesReturnsContainsMatchesBeforeSubsequenceMatches(t *testing.T) {
	names := []string{"VisualStudioCode", "CodeKit", "Cloud9"}

	matches := searchNames("code", names)

	want := "CodeKit, VisualStudioCode"
	if got := strings.Join(matches[:2], ", "); got != want {
		t.Fatalf("searchNames() first matches = %q; want %q", got, want)
	}
}

func TestGenerateCompletionSupportsBash(t *testing.T) {
	completion, err := generateCompletion("bash")
	if err != nil {
		t.Fatalf("generateCompletion() returned error: %v", err)
	}

	if !strings.Contains(completion, "complete -F _ignoreinit_completion ignoreinit") {
		t.Fatalf("bash completion = %q; want bash completion function", completion)
	}
}

func TestGenerateCompletionRejectsUnsupportedShell(t *testing.T) {
	_, err := generateCompletion("elvish")
	if err == nil || !strings.Contains(err.Error(), "unsupported shell") {
		t.Fatalf("generateCompletion() error = %v; want unsupported shell error", err)
	}
}

func TestWriteIgnoreCreatesFileWithContentMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	if err := writeIgnore(path, []byte("bin/\n"), true, false, false); err != nil {
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

	if mode := info.Mode().Perm(); mode&0111 != 0 {
		t.Fatalf("file mode = %v; should not be executable", mode)
	}
}

func TestWriteIgnoreMergesWithBlankLineSeparator(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte("bin/"), gitignoreFileMode); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	if err := writeIgnore(path, []byte("dist/\n"), false, true, false); err != nil {
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

func TestWriteIgnoreMergesCRLFContentWithCleanSeparator(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte("bin/\r\n"), gitignoreFileMode); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	if err := writeIgnore(path, []byte("dist/\n"), false, true, false); err != nil {
		t.Fatalf("writeIgnore() returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if got := string(content); got != "bin/\n\ndist/\n" {
		t.Fatalf("merged content = %q; want normalized blank-line separator", got)
	}
}

func TestWriteIgnoreMergesPreservesDuplicatePatterns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte("bin/\n"), gitignoreFileMode); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	if err := writeIgnore(path, []byte("bin/\ndist/\n"), false, true, false); err != nil {
		t.Fatalf("writeIgnore() returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if got := string(content); got != "bin/\n\nbin/\ndist/\n" {
		t.Fatalf("merged content = %q; want duplicate pattern order preserved", got)
	}
}

func TestWriteIgnoreMergesPreservesOrderSensitiveNegations(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.WriteFile(path, []byte("!important.log\n"), gitignoreFileMode); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	if err := writeIgnore(path, []byte("*.log\n!important.log\n"), false, true, false); err != nil {
		t.Fatalf("writeIgnore() returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if got := string(content); got != "!important.log\n\n*.log\n!important.log\n" {
		t.Fatalf("merged content = %q; want order-sensitive negation preserved", got)
	}
}

func TestWriteIgnoreMergesWithoutDuplicateSectionMarkers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	existing := "# >>> ignoreinit: Go\nbin/\n# <<< ignoreinit: Go\n"
	if err := os.WriteFile(path, []byte(existing), gitignoreFileMode); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	if err := writeIgnore(path, []byte(existing), false, true, false); err != nil {
		t.Fatalf("writeIgnore() returned error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() returned error: %v", err)
	}

	if got := strings.Count(string(content), "# >>> ignoreinit: Go"); got != 1 {
		t.Fatalf("section start count = %d; want 1", got)
	}
	if got := strings.Count(string(content), "# <<< ignoreinit: Go"); got != 1 {
		t.Fatalf("section end count = %d; want 1", got)
	}
}

func TestWriteIgnoreWrapsWriteErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatalf("Mkdir() returned error: %v", err)
	}

	err := writeIgnore(path, []byte("bin/\n"), false, false, false)
	if err == nil || !strings.Contains(err.Error(), "could not write") {
		t.Fatalf("writeIgnore() error = %v; want wrapped write error", err)
	}
}

func TestWriteIgnorePrintsToStdout(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".gitignore")

	output := captureStdout(t, func() {
		if err := writeIgnore(path, []byte("bin/\n"), true, false, true); err != nil {
			t.Fatalf("writeIgnore() returned error: %v", err)
		}
	})

	if !strings.HasPrefix(output, generatedHeader) || !strings.Contains(output, "bin/\n") {
		t.Fatalf("stdout = %q; want generated content", output)
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("printed write created file or returned unexpected error: %v", err)
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

	if got := string(content); got != "# >>> ignoreinit: Go\nbin/\n# <<< ignoreinit: Go\n" {
		t.Fatalf("fetchIgnore() = %q; want sectioned Go template", got)
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

func TestFetchIgnoresCombinesMultipleTemplates(t *testing.T) {
	oldRetryDelay := retryDelay
	retryDelay = 0
	defer func() { retryDelay = oldRetryDelay }()

	client := &fakeGitignoreClient{
		listTemplates:    []string{"Go.gitignore", "Node.gitignore"},
		downloadContents: [][]byte{[]byte("bin/\n"), []byte("node_modules/\r\n")},
	}

	content, err := fetchIgnores([]string{"go", "node"}, client)
	if err != nil {
		t.Fatalf("fetchIgnores() returned error: %v", err)
	}

	want := "# >>> ignoreinit: Go\nbin/\n# <<< ignoreinit: Go\n\n# >>> ignoreinit: Node\nnode_modules/\n# <<< ignoreinit: Node\n"
	if got := string(content); got != want {
		t.Fatalf("fetchIgnores() = %q; want combined templates", got)
	}
}

func TestFetchIgnoresSkipsDuplicateTemplates(t *testing.T) {
	oldRetryDelay := retryDelay
	retryDelay = 0
	defer func() { retryDelay = oldRetryDelay }()

	client := &fakeGitignoreClient{
		listTemplates:   []string{"Go.gitignore"},
		downloadContent: []byte("bin/\n"),
	}

	content, err := fetchIgnores([]string{"go", "Go"}, client)
	if err != nil {
		t.Fatalf("fetchIgnores() returned error: %v", err)
	}

	if client.downloadCalls != 1 {
		t.Fatalf("download calls = %d; want 1", client.downloadCalls)
	}

	if got := strings.Count(string(content), "# >>> ignoreinit: Go"); got != 1 {
		t.Fatalf("section count = %d; want 1", got)
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

	if len(c.downloadContents) > 0 {
		content := c.downloadContents[0]
		c.downloadContents = c.downloadContents[1:]
		return content, nil
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

func captureStdout(t *testing.T, action func()) string {
	t.Helper()

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe() returned error: %v", err)
	}

	os.Stdout = writer
	writerClosed := false
	defer func() {
		os.Stdout = oldStdout
		if !writerClosed {
			if err := writer.Close(); err != nil {
				t.Logf("Close() returned error during cleanup: %v", err)
			}
		}
		if err := reader.Close(); err != nil {
			t.Logf("Close() returned error during cleanup: %v", err)
		}
	}()

	action()

	if err := writer.Close(); err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}
	writerClosed = true

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll() returned error: %v", err)
	}

	return string(output)
}
