# CI Status

[![CI](https://github.com/InkyQuill/DivoomPCMonitorTool/workflows/CI/badge.svg)](https://github.com/InkyQuill/DivoomPCMonitorTool/actions/workflows/ci.yml)

---

## Test Results

| Platform | Status | Details |
|----------|--------|---------|
| **Linux** | [![CI](https://img.shields.io/github/actions-workflow/status/InkyQuill/DivoomPCMonitorTool/CI.yml?label=Linux)](https://github.com/InkyQuill/DivoomPCMonitorTool/actions/workflows/ci.yml) | Backend + Frontend |
| **Windows** | [![CI](https://img.shields.io/github/actions-workflow/status/InkyQuill/DivoomPCMonitorTool/CI.yml?label=Windows)](https://github.com/InkyQuill/DivoomPCMonitorTool/actions/workflows/ci.yml) | Backend |
| **macOS** | [![CI](https://img.shields.io/github/actions-workflow/status/InkyQuill/DivoomPCMonitorTool/CI.yml?label=macOS)](https://github.com/InkyQuill/DivoomPCMonitorTool/actions/workflows/ci.yml) | Backend |

---

## Automated Testing

This repository uses GitHub Actions for continuous integration:

- ✅ **Automated tests** run on every push and pull request
- ✅ **Multi-platform testing** (Linux, Windows, macOS)
- ✅ **Fast feedback** - Know immediately if something breaks
- ✅ **Zero cost** - Free for public repositories

### Test Coverage

Current test coverage:
- **Backend (Rust)**: 111 tests, ~85% coverage
- **Frontend (Svelte)**: 20 tests, 100% coverage of tested modules
- **Total**: 131 tests passing

See [ALL_TESTS_SUMMARY.md](./ALL_TESTS_SUMMARY.md) for details.
