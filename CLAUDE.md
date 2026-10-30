# CLAUDE.md

inkwell: a point-and-click adventure library for Go on Ebitengine, module
`git.teletypegames.org/engines/inkwell`. The documentation lives in the source:
the package comments are the manual and the doc comments the reference. Every
release tag publishes it to the wiki at
https://wiki.teletypegames.org/godoc/inkwell/latest/ ; `README.md` only links there.

## Documenting code: write for Go doc

The reference is generated with doc2go, and CI fails on any exported
identifier without a doc comment (`revive` with `revive.toml`). Every change
keeps the code fully documented:

- Every exported type, function, method (also on unexported receivers'
  exported types), constant and variable gets a doc comment directly above
  it, starting with its name: `// Scene is …`, `// Register adds …`. A
  grouped `const (…)` / `var (…)` block may share one comment above the group,
  but give each exported name its own line comment when the names mean
  different things.
- Say what the thing is or does, what it takes, what it returns, the errors
  and panics, and the state it changes. Not its history: the why of a change
  goes in the commit.
- A re-export in the root package (`type Scene = scene.Scene`,
  `var Say = act.Say`) gets a one-line comment that says what it is and links
  the original: `// Scene is a room of the game; see [scene.Scene].`
- Every package has a package comment in its `doc.go` (`// Package act …`).
  The manual — how the pieces fit, with examples — goes in the package
  comments, under `# Heading` sections.
- Go doc syntax, not Markdown: `[Name]`, `[pkg.Name]` and `[Type.Method]`
  links, `# Heading` (only in package comments), lists with `-`, code as an
  indented block (one tab or extra spaces after `//`). No backticks, no `**`.
- Check before committing:

    go build ./... && go vet ./...
    revive -config revive.toml -set_exit_status ./...     # must print nothing
    doc2go -out /tmp/inkwell-doc -home git.teletypegames.org/engines/inkwell ./...   # look at it

## Releases

A release is a `vX.Y.Z` tag on master. Its pipeline runs the checks above and
then `publish-godoc`: doc2go builds the HTML, `meta.json` records version,
commit and date, and the archive is uploaded as the Gitea generic package
`engines/inkwell-godoc/<version>`; the wiki-fetcher unpacks it under `/godoc/`.
