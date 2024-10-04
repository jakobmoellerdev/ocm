package git_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"os"
	"time"

	_ "embed"

	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mandelsoft/filepath/pkg/filepath"
	"github.com/mandelsoft/vfs/pkg/cwdfs"
	"github.com/mandelsoft/vfs/pkg/osfs"

	"ocm.software/ocm/api/datacontext/attrs/tmpcache"
	"ocm.software/ocm/api/datacontext/attrs/vfsattr"
	"ocm.software/ocm/api/ocm"
	"ocm.software/ocm/api/ocm/cpi"
	me "ocm.software/ocm/api/ocm/extensions/accessmethods/git"
)

//go:embed testdata/repo
var testData embed.FS

var _ = Describe("Method", func() {
	var (
		ctx                 ocm.Context
		expectedBlobContent []byte
		accessSpec          *me.AccessSpec
	)

	ctx = ocm.New()

	BeforeEach(func() {
		tempVFS, err := cwdfs.New(osfs.New(), GinkgoT().TempDir())
		Expect(err).ToNot(HaveOccurred())
		tmpcache.Set(ctx, &tmpcache.Attribute{Path: ".", Filesystem: tempVFS})
		vfsattr.Set(ctx, tempVFS)
	})

	BeforeEach(func() {
		repoDir := GinkgoT().TempDir() + filepath.PathSeparatorString + "repo"

		repo := Must(git.PlainInit(repoDir, false))

		repoBase := filepath.Join("testdata", "repo")
		repoTestData := Must(testData.ReadDir(repoBase))

		for _, entry := range repoTestData {
			path := filepath.Join(repoBase, entry.Name())
			repoPath := filepath.Join(repoDir, entry.Name())

			file := Must(testData.Open(path))

			fileInRepo := Must(os.OpenFile(
				repoPath,
				os.O_CREATE|os.O_RDWR|os.O_TRUNC,
				0o600,
			))

			Must(io.Copy(fileInRepo, file))

			Expect(fileInRepo.Close()).To(Succeed())
			Expect(file.Close()).To(Succeed())
		}

		wt := Must(repo.Worktree())
		Expect(wt.AddGlob("*")).To(Succeed())
		Must(wt.Commit("OCM Test Commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "OCM Test",
				Email: "dummy@ocm.software",
				When:  time.Now(),
			},
		}))

		accessSpec = me.New(
			fmt.Sprintf("file://%s", repoDir),
			string(plumbing.Master),
			".",
		)

		expectedBlobContent = Must(testData.ReadFile(filepath.Join("testdata", "repo", "file_in_repo")))
	})

	It("downloads artifacts", func() {
		m := Must(accessSpec.AccessMethod(&cpi.DummyComponentVersionAccess{Context: ctx}))
		content := Must(m.Get())
		unzippedContent := Must(gzip.NewReader(bytes.NewReader(content)))

		r := tar.NewReader(unzippedContent)

		file := Must(r.Next())
		Expect(file.Name).To(Equal("file_in_repo"))
		Expect(file.Size).To(Equal(int64(len(expectedBlobContent))))

		data := Must(io.ReadAll(r))
		Expect(data).To(Equal(expectedBlobContent))
	})
})
