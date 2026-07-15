package repomanagement

import "testing"

func TestExtractRepoInfo(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    RepoInfo
		wantErr bool
	}{
		{
			name: "https simple",
			raw:  "https://github.com/foo/bar.git",
			want: RepoInfo{
				Scheme:        "https",
				Host:          "github.com",
				Port:          "",
				TopLevelOwner: "foo",
				FullOwnerPath: "foo",
				Repo:          "bar",
				RawURL:        "https://github.com/foo/bar.git",
			},
		},
		{
			name: "https without .git suffix",
			raw:  "https://github.com/foo/bar",
			want: RepoInfo{
				Scheme:        "https",
				Host:          "github.com",
				Port:          "",
				TopLevelOwner: "foo",
				FullOwnerPath: "foo",
				Repo:          "bar",
				RawURL:        "https://github.com/foo/bar",
			},
		},
		{
			name: "https with port and subgroup",
			raw:  "https://gitlab.example.com:8443/group/subgroup/project.git",
			want: RepoInfo{
				Scheme:        "https",
				Host:          "gitlab.example.com",
				Port:          "8443",
				TopLevelOwner: "group",
				FullOwnerPath: "group/subgroup",
				Repo:          "project",
				RawURL:        "https://gitlab.example.com:8443/group/subgroup/project.git",
			},
		},
		{
			name: "scp-style ssh with user",
			raw:  "git@github.com:foo/bar.git",
			want: RepoInfo{
				Scheme:        "ssh",
				Host:          "github.com",
				Port:          "",
				TopLevelOwner: "foo",
				FullOwnerPath: "foo",
				Repo:          "bar",
				RawURL:        "git@github.com:foo/bar.git",
			},
		},
		{
			name: "scp-style ssh without user",
			raw:  "host:owner/repo.git",
			want: RepoInfo{
				Scheme:        "ssh",
				Host:          "host",
				Port:          "",
				TopLevelOwner: "owner",
				FullOwnerPath: "owner",
				Repo:          "repo",
				RawURL:        "host:owner/repo.git",
			},
		},
		{
			name: "explicit ssh scheme URL",
			raw:  "ssh://git@host.example.com:2222/owner/repo.git",
			want: RepoInfo{
				Scheme:        "ssh",
				Host:          "host.example.com",
				Port:          "2222",
				TopLevelOwner: "owner",
				FullOwnerPath: "owner",
				Repo:          "repo",
				RawURL:        "ssh://git@host.example.com:2222/owner/repo.git",
			},
		},
		{
			name:    "empty input",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			raw:     "   ",
			wantErr: true,
		},
		{
			name:    "missing repo segment",
			raw:     "https://github.com/onlyowner",
			wantErr: true,
		},
		{
			name:    "no owner or repo at all",
			raw:     "just-a-bare-string",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ExtractRepoInfo(c.raw)
			if c.wantErr {
				if err == nil {
					t.Fatalf("ExtractRepoInfo(%q) expected error, got %+v", c.raw, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ExtractRepoInfo(%q) unexpected error: %v", c.raw, err)
			}
			if got != c.want {
				t.Errorf("ExtractRepoInfo(%q) =\n%+v\nwant:\n%+v", c.raw, got, c.want)
			}
		})
	}
}
