# Version History

## Current Version: 1.5.1

### Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html):
- **MAJOR** version for incompatible API changes
- **MINOR** version for backwards-compatible functionality additions
- **PATCH** version for backwards-compatible bug fixes

---

## Version 1.5.1

**Release Date**: 2026-03-21

**Type**: Minor Release (New Features)

### Configuration Changes (Non-Breaking)
- Configuration file location changed from `~/.ssmshrc` to `~/.config/ssmsh/config` (XDG Base Directory compliant)
  - **Automatic migration** with backup on first run - no user action required
  - Legacy location still fully supported for backwards compatibility
  - Zero downtime upgrade path

### Major Features
- **Tab Completion** for AWS Parameter Store paths and parameters
  - Async/non-blocking architecture prevents terminal freezing
  - Two-TAB pattern: First TAB starts fetch, second shows results
  - Cache warmup on startup for instant common path completion
  - Two-tier caching (in-memory + persistent)
  - Configurable via runtime commands and config file

- **XDG Base Directory Compliance**
  - Config: `~/.config/ssmsh/config`
  - Cache: `~/.cache/ssmsh/cache.gob.gz`
  - History: `~/.config/ssmsh/history`

### New Commands
- `completion [true|false|stats|clear-cache|save-cache|reload-cache]`
- `config [show|edit]`

### New Flags
- `--generate-config` - Generate default configuration file
- `get -d` / `get --decrypt` - Inline parameter decryption

### Improvements
- Region detection from AWS profile config
- Debug mode via `SSMSH_DEBUG=1` environment variable
- Graceful shutdown on SIGINT/SIGTERM
- Performance optimizations with intelligent caching

### Statistics
- **New Lines of Code**: ~3,400 (2,500 production + 900 tests)
- **Test Coverage**: 32.6% overall (config: 52.4%, parameterstore: 71.4%)
- **Performance**:
  - Memory cache: < 1 μs per lookup
  - Persistent cache: < 1 ms to load
  - Background fetch timeout: 2 seconds

---

## Version 1.5.0

**Release Date**: 2024

### Features
- Parameter value visibility improvements
- Enhanced put validation

### Fixes
- Restore parameter value visibility in get command
- Improve put command validation

---

## Version 1.x.x

Earlier versions focused on core Parameter Store functionality:
- Interactive shell for AWS Systems Manager Parameter Store
- Commands: cd, ls, get, put, rm, cp, mv
- Multi-region support
- Parameter history
- Advanced parameters with policies
- Batch command execution
- Profile switching

---

## Migration Notes

### Upgrading to 2.0.0

**Automatic Migration** (Recommended):
1. Simply run `ssmsh`
2. Migration happens automatically on first startup
3. Original config backed up to `~/.ssmshrc.backup`

**Manual Migration**:
```bash
mkdir -p ~/.config/ssmsh
cp ~/.ssmshrc ~/.config/ssmsh/config
```

**Continue Using Legacy Location**:
```bash
export SSMSH_CONFIG=~/.ssmshrc
# or
ssmsh --config ~/.ssmshrc
```

### Configuration Priority
1. `--config` flag (highest)
2. `SSMSH_CONFIG` environment variable
3. `~/.config/ssmsh/config` (XDG)
4. `~/.ssmshrc` (legacy)

---

## Roadmap

### Planned for Future Versions
- Fuzzy matching for tab completion
- Substring search in parameters
- Completion for parameter values
- Smart caching with user pattern learning
- Offline mode with full parameter cache
- Multi-region cache warmup

### Under Consideration
- Export/import functionality
- Integration with CloudWatch Events
- Global parameter search
- Copy between AWS accounts
- Read parameters as environment variables
- Regex/globbing support

---

## Version Numbering Guide

- **1.5.x**: Tab completion, XDG-compliant configuration (with backward compatibility)
- **1.4.x**: Core Parameter Store features with legacy configuration

---

**Current Version**: 1.5.1
**Previous Stable**: 1.5.0
**Latest Features**: Tab completion, XDG config support, inline decrypt
