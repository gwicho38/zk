package file

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/mickael-menu/zk/core/zk"
	"github.com/mickael-menu/zk/util"
	"github.com/mickael-menu/zk/util/assert"
	"github.com/mickael-menu/zk/util/fixtures"
)

var root = fixtures.Path("walk")

func TestWalkRootDir(t *testing.T) {
	dir := zk.Dir{Name: "", Path: root}
	testEqual(t, Walk(dir, "md", &util.NullLogger), []Path{
		{Dir: "", Filename: "a.md", Abs: filepath.Join(root, "a.md")},
		{Dir: "", Filename: "b.md", Abs: filepath.Join(root, "b.md")},
		{Dir: "dir1", Filename: "a.md", Abs: filepath.Join(root, "dir1/a.md")},
		{Dir: "dir1", Filename: "b.md", Abs: filepath.Join(root, "dir1/b.md")},
		{Dir: "dir1/dir1", Filename: "a.md", Abs: filepath.Join(root, "dir1/dir1/a.md")},
		{Dir: "dir2", Filename: "a.md", Abs: filepath.Join(root, "dir2/a.md")},
	})
}

func TestWalkSubDir(t *testing.T) {
	dir := zk.Dir{Name: "dir1", Path: filepath.Join(root, "dir1")}
	testEqual(t, Walk(dir, "md", &util.NullLogger), []Path{
		{Dir: "dir1", Filename: "a.md", Abs: filepath.Join(root, "dir1/a.md")},
		{Dir: "dir1", Filename: "b.md", Abs: filepath.Join(root, "dir1/b.md")},
		{Dir: "dir1/dir1", Filename: "a.md", Abs: filepath.Join(root, "dir1/dir1/a.md")},
	})
}

func TestWalkSubSubDir(t *testing.T) {
	dir := zk.Dir{Name: "dir1/dir1", Path: filepath.Join(root, "dir1/dir1")}
	testEqual(t, Walk(dir, "md", &util.NullLogger), []Path{
		{Dir: "dir1/dir1", Filename: "a.md", Abs: filepath.Join(root, "dir1/dir1/a.md")},
	})
}

func date(s string) time.Time {
	date, _ := time.Parse(time.RFC3339, s)
	return date
}

func testEqual(t *testing.T, actual <-chan Metadata, expected []Path) {
	popExpected := func() (Path, bool) {
		if len(expected) == 0 {
			return Path{}, false
		}
		item := expected[0]
		expected = expected[1:]
		return item, true
	}

	for act := range actual {
		exp, ok := popExpected()
		if !ok {
			t.Errorf("More paths available than expected")
			return
		}
		assert.Equal(t, act.Path, exp)
		assert.NotNil(t, act.Modified)
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected paths: %v", expected)
	}
}
