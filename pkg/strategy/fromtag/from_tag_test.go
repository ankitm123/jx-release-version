package fromtag

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	archiver "github.com/jm33-m0/arc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name             string
		strategy         Strategy
		expected         *semver.Version
		expectedErrorMsg string
	}{
		{
			name: "no prefix",
			strategy: Strategy{
				Dir: "testdata/git-repo",
			},
			expected: semver.MustParse("v2.0.0"),
		},
		{
			name: "fetch tags first",
			strategy: Strategy{
				Dir:       "testdata/git-repo",
				FetchTags: true,
			},
			expected: semver.MustParse("v2.1.0"),
		},
		{
			name: "v1 prefix",
			strategy: Strategy{
				Dir:        "testdata/git-repo",
				TagPattern: "v1",
			},
			expected: semver.MustParse("v1.1.0"),
		},
		{
			name: "v1.0 prefix",
			strategy: Strategy{
				Dir:        "testdata/git-repo",
				TagPattern: "v1.0",
			},
			expected: semver.MustParse("v1.0.1"),
		},
		{
			name: "v3 prefix",
			strategy: Strategy{
				Dir:        "testdata/git-repo",
				TagPattern: "v3",
			},
			expectedErrorMsg: "no semver tags with pattern \"v3\" found",
		},
	}

	// the git repo is stored as a tar.gz archive to make it easy to commit
	gitRepoPath := filepath.Join("testdata", "git-repo")
	err := os.RemoveAll(gitRepoPath)
	require.NoErrorf(t, err, "failed to delete %s", gitRepoPath)
	err = archiver.Unarchive(filepath.Join("testdata", "git-repo.tar.gz"), gitRepoPath)
	require.NoErrorf(t, err, "failed to decompress git repository at %s", gitRepoPath)
	// and it has an origin remote repo (with more tags)
	gitRepoOriginPath := filepath.Join("testdata", "git-repo-origin")
	err = os.RemoveAll(gitRepoOriginPath)
	require.NoErrorf(t, err, "failed to delete %s", gitRepoOriginPath)
	err = archiver.Unarchive(filepath.Join("testdata", "git-repo-origin.tar.gz"), gitRepoOriginPath)
	require.NoErrorf(t, err, "failed to decompress git repository at %s", gitRepoOriginPath)
	// link the 2 repos together
	repo, err := git.PlainOpen(gitRepoPath)
	require.NoErrorf(t, err, "failed to open git repo at %s", gitRepoPath)
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{gitRepoOriginPath},
		Fetch: []config.RefSpec{config.RefSpec("+refs/heads/*:refs/remotes/origin/*")},
	})
	require.NoError(t, err, "failed to create git remote")

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := test.strategy.ReadVersion()
			if len(test.expectedErrorMsg) > 0 {
				require.EqualError(t, err, test.expectedErrorMsg)
				assert.Nil(t, actual)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			}
		})
	}
}
