// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original timestamp: 2025/04/21 09:44
// Original filename: /repomanagement.go

package helperFunctions

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
)

type RepoInfo struct {
	Scheme        string // e.g., https, ssh ("" if unknown)
	Host          string // e.g., git.example.com
	Port          string // optional (if included in the URL)
	TopLevelOwner string // e.g., group or username
	FullOwnerPath string // e.g., group/subgroup
	Repo          string // e.g., repo (no .git)
	RawURL        string // original input, for debug
}

// ExtractRepoInfo: extracts various info from the URL or SSH-style URL
// Returns an error if the URL is invalid or if the repo name is missing
func ExtractRepoInfo(raw string) (RepoInfo, error) {
	original := raw
	var scheme, host, port string

	if strings.TrimSpace(raw) == "" {
		return RepoInfo{}, errors.New("empty input")
	}

	// Handle SCP-style URLs: git@host:owner/repo(.git)
	if strings.Contains(raw, ":") && !strings.Contains(raw, "://") {
		// Example: git@host:org/repo.git
		userHostSplit := strings.SplitN(raw, "@", 2)
		pathSplit := strings.SplitN(raw, ":", 2)

		if len(pathSplit) == 2 {
			scheme = "ssh"
			raw = "/" + pathSplit[1]

			if len(userHostSplit) == 2 {
				host = strings.Split(userHostSplit[1], ":")[0]
			} else {
				host = strings.Split(pathSplit[0], ":")[0]
			}
		}
	} else {
		// Handle standard URLs: http(s)://host[:port]/owner/repo(.git)
		if u, err := url.Parse(raw); err == nil {
			scheme = u.Scheme
			host = u.Hostname()
			port = u.Port()
			raw = u.Path
		}
	}

	raw = strings.TrimSuffix(raw, "/")
	raw = filepath.ToSlash(raw)
	segments := strings.Split(strings.TrimPrefix(raw, "/"), "/")

	if len(segments) < 2 {
		return RepoInfo{}, errors.New("could not determine owner and repo from URL")
	}

	repo := strings.TrimSuffix(segments[len(segments)-1], ".git")
	ownerPath := strings.Join(segments[:len(segments)-1], "/")
	topOwner := strings.Split(ownerPath, "/")[0]

	if topOwner == "" || repo == "" {
		return RepoInfo{}, errors.New("missing owner or repo name")
	}

	return RepoInfo{
		Scheme:        scheme,
		Host:          host,
		Port:          port,
		TopLevelOwner: topOwner,
		FullOwnerPath: ownerPath,
		Repo:          repo,
		RawURL:        original,
	}, nil
}
