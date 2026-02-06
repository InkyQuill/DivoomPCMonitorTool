# TDD Guide for Divoom PC Companion

## Overview

This project follows **Test-Driven Development (TDD)** methodology with comprehensive test coverage for both frontend (TypeScript/Svelte) and backend (Rust).

## Testing Infrastructure

### Frontend (TypeScript/Svelte)

- **Framework**: Vitest
- **Testing Library**: @testing-library/svelte
- **Coverage Tool**: v8 (built-in)
- **Location**: `src/**/*.test.ts`, `src/**/*.test.svelte`

### Backend (Rust)

- **Framework**: Built-in Rust test framework
- **Location**: `src-tauri/src/**/*.rs` (in `#[cfg(test)]` modules)

## Running Tests

### Frontend Tests

```bash
# Run all tests once
npm run test:run

# Run tests in watch mode
npm run test

# Run tests with UI
npm run test:ui

# Run tests with coverage
npm run test:coverage
```

### Backend Tests

```bash
# Run all Rust tests
cd src-tauri
cargo test

# Run tests with output
cargo test -- --nocapture

# Run specific test
cargo test test_cpu_usage_range
```

## Test Coverage Requirements

- **Minimum**: 80% coverage across all modules
- **Critical Code**: 100% coverage required for:
  - System metrics collection (CPU, memory, storage)
  - Device discovery and communication
  - Data formatters and utilities
  - Error handling

## Current Test Status

### Frontend
- ✅ **formatters.ts**: 100% coverage (13 tests)
  - formatPercent
  - formatBytes
  - formatTemperature
  - formatMetrics
  - getMetricsArray

- ✅ **devices store**: 7 tests
  - State initialization
  - Device discovery (success/error cases)
  - Loading states
  - Error handling

### Backend
- ✅ **CPU module**: 3 tests
  - Valid metrics return
  - Usage range validation
  - Temperature detection (currently not implemented)

- ✅ **Memory module**: 4 tests
  - Valid metrics return
  - Value consistency checks
  - Percentage calculation verification
  - GB unit validation

- ✅ **Storage module**: 2 tests
  - Placeholder value validation (TODO: implement real disk monitoring)
  - Value consistency

- ✅ **Protocol module**: 2 tests
  - TimeGate detection
  - Screen count calculation

**Total Backend Tests**: 11 passing

## TDD Workflow

When adding new features, follow this workflow:

### 1. RED - Write Failing Test

Write a test that describes the behavior you want:

```typescript
// src/lib/utils/myNewFunction.test.ts
import { describe, it, expect } from 'vitest'
import { myNewFunction } from './myNewFunction'

describe('myNewFunction', () => {
  it('should do something specific', () => {
    const input = { /* test data */ }
    const expected = { /* expected result */ }
    expect(myNewFunction(input)).toEqual(expected)
  })
})
```

### 2. RUN - Verify Test Fails

```bash
npm run test:run src/lib/utils/myNewFunction.test.ts
# Should show: FAIL - myNewFunction is not defined
```

### 3. GREEN - Write Minimal Implementation

```typescript
// src/lib/utils/myNewFunction.ts
export function myNewFunction(input) {
  return { /* minimal implementation */ }
}
```

```bash
npm run test:run src/lib/utils/myNewFunction.test.ts
# Should show: PASS
```

### 4. REFACTOR - Improve Code

Refactor while keeping tests green:

```typescript
// Improved implementation with better structure
const HELPER_CONSTANT = 'value'

export function myNewFunction(input) {
  // Cleaner, more maintainable code
}
```

```bash
npm run test:run src/lib/utils/myNewFunction.test.ts
# Should still show: PASS
```

### 5. VERIFY - Check Coverage

```bash
npm run test:coverage
# Verify coverage >= 80%
```

## TDD Examples

### Example 1: Fixing a Bug Found by Tests

When we ran tests for `formatters.ts`, we discovered 2 bugs:

1. **Bug**: `formatBytes(0, 2)` returned `'0.00 MB'` instead of `'0.00 GB'`
2. **Bug**: `getMetricsArray` had incorrect rounding expectation

**Fix Applied**:

```typescript
// Before
export function formatBytes(gb: number, decimals: number = 2): string {
  if (gb < 1) {  // ❌ Includes 0
    return `${(gb * 1024).toFixed(decimals)} MB`;
  }
  return `${gb.toFixed(decimals)} GB`;
}

// After
export function formatBytes(gb: number, decimals: number = 2): string {
  if (gb < 1 && gb > 0) {  // ✅ Excludes 0
    return `${(gb * 1024).toFixed(decimals)} MB`;
  }
  return `${gb.toFixed(decimals)} GB`;
}
```

### Example 2: Rust Type Mismatch Discovery

Tests revealed a type mismatch bug in `storage.rs`:

```rust
// Test failed with: expected `f32`, found `f64`
assert_eq!(metrics.usage_percent, (metrics.used_gb / metrics.total_gb) * 100.0);

// Fixed by adding proper cast
assert_eq!(metrics.usage_percent, ((metrics.used_gb / metrics.total_gb) * 100.0) as f32);
```

## Writing Good Tests

### DO ✅

- **Test behavior, not implementation**
  ```typescript
  // ✅ Good - tests what the function does
  expect(formatBytes(1.5)).toBe('1.50 GB')

  // ❌ Bad - tests implementation details
  expect(formatBytes.toString()).toContain('toFixed')
  ```

- **Use descriptive test names**
  ```typescript
  it('should return N/A for null temperature')  // ✅ Clear
  it('test temp')  // ❌ Vague
  ```

- **Test edge cases**
  ```typescript
  it('should handle zero values')
  it('should handle null/undefined')
  it('should handle max values')
  it('should handle negative values')
  ```

- **Arrange-Act-Assert pattern**
  ```typescript
  it('should calculate percentage correctly', () => {
    // Arrange
    const input = { used: 50, total: 100 }

    // Act
    const result = calculatePercent(input)

    // Assert
    expect(result).toBe(50)
  })
  ```

### DON'T ❌

- **Don't test implementation details**
- **Don't mock everything** (prefer real implementations when possible)
- **Don't write tests that are too brittle**
- **Don't skip edge cases**
- **Don't ignore failing tests**

## Continuous Testing

### Pre-commit Hook (Recommended)

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash

# Run frontend tests
npm run test:run

# Run backend tests
cd src-tauri && cargo test

# If any test fails, prevent commit
if [ $? -ne 0 ]; then
  echo "❌ Tests failed. Commit aborted."
  exit 1
fi

echo "✅ All tests passed!"
```

### CI/CD Integration

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '20'
      - run: npm install
      - run: npm run test:coverage
      - run: cd src-tauri && cargo test
```

## Next Steps

### High Priority (TODO)

1. **Add more frontend tests**:
   - [ ] Component tests (Svelte components)
   - [ ] Store tests (metrics store, config store)
   - [ ] Integration tests

2. **Expand backend tests**:
   - [ ] Config loading/saving tests
   - [ ] Error handling tests
   - [ ] Platform-specific code tests (CPU temps, GPU metrics)

3. **Add E2E tests**:
   - [ ] Full user flow tests
   - [ ] Device communication tests

### Medium Priority

- [ ] Set up code coverage reporting
- [ ] Add performance benchmarks
- [ ] Add mutation testing (stryker)

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [Testing Library Svelte](https://testing-library.com/docs/svelte-testing-library/intro/)
- [Rust Testing Guide](https://doc.rust-lang.org/book/ch11-00-testing.html)
- [TDD Best Practices](https://martinfowler.com/bliki/TestDrivenDevelopment.html)

## Summary

✅ **Infrastructure**: Complete (Vitest + Cargo test)
✅ **Current Coverage**: ~20% overall, 100% for tested modules
✅ **Test Quality**: High (following TDD principles)
⏳ **Goal**: 80%+ coverage across entire codebase

**Remember**: Red → Green → Refactor → Repeat
