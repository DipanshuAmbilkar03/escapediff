// Package tests runs escapediff the way the skill does: on a real git
// repository, building two revisions of the example shop service.
package tests

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var bin string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "escapediff-bin")
	if err != nil {
		panic(err)
	}
	bin = filepath.Join(dir, "escapediff")
	// Build the single file, as the skill runs it.
	out, err := exec.Command("go", "build", "-o", bin, "../plugin/scripts/escapediff.go").CombinedOutput()
	if err != nil {
		panic(string(out))
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

type result struct {
	NewHeap    []struct{ Func, Text string }
	NewLeaks   []struct{ Func, Text string }
	LostDevirt []struct{ Func string } `json:"lostDevirtualization"`
	LostInline []struct {
		Name      string
		CallSites int `json:"callSitesInBase"`
	} `json:"noLongerInlined"`
	GoneHeap []struct{ Text string }
}

// copyTree copies the Go files and go.mod of an example version into dst.
func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, p)
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		os.MkdirAll(filepath.Dir(target), 0o755)
		return os.WriteFile(target, b, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// repo makes a git repository with the before version committed.
func repo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	copyTree(t, "../examples/shop/before", dir)
	git(t, dir, "init", "-q")
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-q", "-m", "before")
	return dir
}

func escapediff(t *testing.T, dir string, args ...string) (int, string, string) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	err := cmd.Run()
	code := 0
	if e, ok := err.(*exec.ExitError); ok {
		code = e.ExitCode()
	} else if err != nil {
		t.Fatal(err)
	}
	return code, out.String(), errb.String()
}

func snapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && !strings.Contains(p, string(filepath.Separator)+".git"+string(filepath.Separator)) {
			b, _ := os.ReadFile(p)
			files[p] = string(b)
		}
		return nil
	})
	return files
}

func TestWorkingTreeAgainstHEAD(t *testing.T) {
	dir := repo(t)
	copyTree(t, "../examples/shop/after", dir)
	before := snapshot(t, dir)

	code, out, errs := escapediff(t, dir, "git", "--json", "--fail", "HEAD")
	if code != 1 {
		t.Fatalf("exit %d, want 1\n%s\n%s", code, out, errs)
	}
	var r result
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatal(err)
	}
	names := map[string]int{}
	for _, f := range r.LostInline {
		names[f.Name] = f.CallSites
	}
	if n, ok := names["(*Order).Total"]; !ok || n != 2 {
		t.Errorf("(*Order).Total should lose inlining at 2 call sites: %+v", r.LostInline)
	}
	if len(r.LostDevirt) != 1 || r.LostDevirt[0].Func != "Receipt" {
		t.Errorf("lost devirtualization = %+v", r.LostDevirt)
	}
	leaks := map[string]bool{}
	for _, l := range r.NewLeaks {
		leaks[l.Func+": "+l.Text] = true
	}
	for _, want := range []string{"(*Order).Total: leaking param: o", "Batch: leaking param content: orders"} {
		if !leaks[want] {
			t.Errorf("missing leak %q in %+v", want, r.NewLeaks)
		}
	}
	if len(r.GoneHeap) != 1 || r.GoneHeap[0].Text != "id escapes to heap" {
		t.Errorf("gone heap = %+v", r.GoneHeap)
	}
	// The tool must not change the working tree or write binaries into it.
	after := snapshot(t, dir)
	if len(after) != len(before) {
		t.Errorf("working tree changed: %d files before, %d after", len(before), len(after))
	}
	for p, c := range before {
		if after[p] != c {
			t.Errorf("%s changed", p)
		}
	}
}

func TestTwoRefsAndSave(t *testing.T) {
	dir := repo(t)
	copyTree(t, "../examples/shop/after", dir)
	git(t, dir, "commit", "-q", "-am", "after")
	save := t.TempDir()
	code, out, errs := escapediff(t, dir, "git", "--head", "HEAD", "--save", save, "HEAD~1", "./orders")
	if code != 0 {
		t.Fatalf("exit %d\n%s", code, errs)
	}
	if !strings.Contains(out, "No longer inlined:") || !strings.Contains(out, "(*Order).Total") {
		t.Errorf("output:\n%s", out)
	}
	for _, name := range []string{"base.log", "head.log"} {
		b, err := os.ReadFile(filepath.Join(save, name))
		if err != nil || !bytes.Contains(b, []byte("can inline")) {
			t.Errorf("%s not saved: %v", name, err)
		}
	}
	// The saved logs give the same answer through compare.
	code, out2, _ := escapediff(t, dir, "compare", filepath.Join(save, "base.log"), filepath.Join(save, "head.log"))
	if code != 0 || out2 != out {
		t.Errorf("compare on saved logs differs:\n%s\n---\n%s", out2, out)
	}
}

func TestCommentsOnlyChangeIsQuiet(t *testing.T) {
	dir := repo(t)
	copyTree(t, "../examples/shop/comments", dir)
	code, out, errs := escapediff(t, dir, "git", "--fail", "HEAD")
	if code != 0 || !strings.HasPrefix(out, "No changes in heap escapes") {
		t.Fatalf("exit %d\n%s\n%s", code, out, errs)
	}
}

func TestSubdirectory(t *testing.T) {
	dir := repo(t)
	copyTree(t, "../examples/shop/after", dir)
	code, out, errs := escapediff(t, filepath.Join(dir, "orders"), "git", "HEAD", ".")
	if code != 0 || !strings.Contains(out, "(*Order).Total") {
		t.Fatalf("exit %d\n%s\n%s", code, out, errs)
	}
}

func TestBrokenBuildIsReported(t *testing.T) {
	dir := repo(t)
	os.WriteFile(filepath.Join(dir, "orders", "broken.go"), []byte("package orders\n\nfunc Broken() int { return \"x\" }\n"), 0o644)
	code, _, errs := escapediff(t, dir, "git", "HEAD")
	if code != 2 || !strings.Contains(errs, "building the working tree failed") || !strings.Contains(errs, "broken.go") {
		t.Fatalf("exit %d\n%s", code, errs)
	}
}

func TestBadRef(t *testing.T) {
	dir := repo(t)
	code, _, errs := escapediff(t, dir, "git", "no-such-branch")
	if code != 2 || !strings.Contains(errs, "no-such-branch") {
		t.Fatalf("exit %d\n%s", code, errs)
	}
}
