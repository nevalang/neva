package git

import (
	"fmt"
	"strings"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Checkout updates the repository worktree to the requested revision.
// The revision can be provided as a full reference (refs/... path), a branch
// name, a tag, or any revision string understood by go-git's ResolveRevision.
func Checkout(repo *gitlib.Repository, revision string) error {
	if revision == "" {
		return fmt.Errorf("revision must not be empty")
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("open worktree: %w", err)
	}

	if strings.HasPrefix(revision, "refs/") {
		if err := worktree.Checkout(&gitlib.CheckoutOptions{Branch: plumbing.ReferenceName(revision)}); err == nil {
			return nil
		}
	}

	if err := worktree.Checkout(&gitlib.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(revision)}); err == nil {
		return nil
	}

	if tagRef, err := repo.Reference(plumbing.NewTagReferenceName(revision), true); err == nil {
		if err := worktree.Checkout(&gitlib.CheckoutOptions{Hash: tagRef.Hash()}); err == nil {
			return nil
		}
	}

	if hash, err := repo.ResolveRevision(plumbing.Revision(revision)); err == nil {
		if err := worktree.Checkout(&gitlib.CheckoutOptions{Hash: *hash}); err == nil {
			return nil
		}
	}

	return fmt.Errorf("revision %q not found", revision)
}
