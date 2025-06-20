# Releasing

1. **Update the version** in `cmd/tap-incident/cmd/VERSION`:
   ```
   echo "0.6.0" > cmd/tap-incident/cmd/VERSION
   ```

2. **Update the CHANGELOG.md** with your changes:

3. **Commit and push** to master:
   ```bash
   git add cmd/tap-incident/cmd/VERSION CHANGELOG.md
   git commit -m "Release v0.6.0"
   git push origin master
   ```

4. **The automation takes over**:
   - GitHub Actions detects the version change
   - Creates and pushes tag `v0.6.0`
   - GoReleaser builds and publishes everything
