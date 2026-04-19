package git

import (
	"errors"
	"path/filepath"
	"strings"
)

type RepoSpec struct {
	Location string
	Revision string
	local    bool
}

func ParseRepoSpec(raw string) (RepoSpec, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return RepoSpec{}, errors.New("repository specification must not be empty")
	}

	source, revision := splitRevision(trimmed)
	if source == "" {
		return RepoSpec{}, errors.New("repository location must not be empty")
	}

	spec := RepoSpec{Location: source, Revision: revision, local: isLocalPath(source)}
	return spec, nil
}

func splitRevision(input string) (string, string) {
	atIdx := strings.LastIndex(input, "@")
	if atIdx == -1 {
		return strings.TrimSpace(input), ""
	}

	slashIdx := strings.LastIndex(input, "/")
	if slashIdx >= atIdx {
		return strings.TrimSpace(input), ""
	}

	revision := strings.TrimSpace(input[atIdx+1:])
	if revision == "" {
		return strings.TrimSpace(input[:atIdx]), ""
	}

	return strings.TrimSpace(input[:atIdx]), revision
}

func isLocalPath(path string) bool {
	if strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../") {
		return true
	}
	if strings.HasPrefix(path, "~/") {
		return true
	}
	if filepath.IsAbs(path) {
		return true
	}
	if strings.HasPrefix(path, "file://") {
		return true
	}
	return false
}

func (s RepoSpec) CloneURL() string {
	if s.local {
		return s.Location
	}
	if strings.HasPrefix(s.Location, "git@") {
		return s.Location
	}
	if strings.Contains(s.Location, "://") {
		return s.Location
	}
	return "https://" + s.Location
}

func (s RepoSpec) IsLocal() bool {
	return s.local
}
