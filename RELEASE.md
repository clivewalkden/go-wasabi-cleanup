# Release

To release the application follow these steps:

1. Merge development into main `git merge development`
2. Make sure dependencies are up to date `go mod tidy`
3. Tag the release `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push the tag `git push origin v1.0.0`
5. Add the tokens:
   1. `GITHUB_TOKEN` - GitHub token
   2. `HOMEBREW_GITHUB_API_TOKEN` - Homebrew GitHub token
6. Run `goreleaser release --clean`

The release will be created and the binaries will be uploaded to the GitHub release page. The Homebrew tap will also be updated.