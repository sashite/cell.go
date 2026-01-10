# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-01-10

### Changed

- **BREAKING:** `FromIndices` and `MustFromIndices` now use variadic arguments instead of a slice for a more natural Go API.

  ```go
  // Before (v1.0.0)
  cell.FromIndices([]int{4, 3})
  cell.MustFromIndices([]int{1, 1, 1})

  // After (v1.1.0)
  cell.FromIndices(4, 3)
  cell.MustFromIndices(1, 1, 1)

  // Migration: add spread operator to existing slice usage
  cell.FromIndices(indices...)
  ```

## [1.0.0] - 2025-01-10

### Added

- Initial implementation of the [CELL Specification v1.0.0](https://sashite.dev/specs/cell/1.0.0/).
- `Valid` function to validate CELL coordinate strings.
- `Parse` and `MustParse` functions to parse coordinates into dimensional components.
- `Dimensions` function to count the number of dimensions in a coordinate.
- `ToIndices` and `MustToIndices` functions to convert coordinates to 0-indexed integers.
- `FromIndices` and `MustFromIndices` functions to convert indices back to coordinates.
- `Regex` function to access the CELL validation regular expression.
- Support for extended alphabet notation (aa, ab, ..., zz, aaa, ...).
- Comprehensive test suite with game-specific examples (Chess, Shōgi, 3D Tic-Tac-Toe).

[1.1.0]: https://github.com/sashite/cell.go/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/sashite/cell.go/releases/tag/v1.0.0
