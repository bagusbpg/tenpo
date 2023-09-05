# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.7] - 2023-09-05

### Fixed

- Use request context when logging

## [0.1.6] - 2023-08-28

### Changed

- Make Server struct pubic

## [0.1.5] - 2023-08-28

### Changed

- Delete Server interface

## [0.1.4] - 2023-08-25

### Changed

- Bump Go version to 1.21

## [0.1.3] - 2023-05-21

### Added

- Embed RequestID in request's context

## [0.1.2] - 2023-05-13

### Changed

- Include method PATCH when AddRoute
- Log handler name

## [0.1.1] - 2023-05-07

### Fixed

- Return error when failed to ListenAndServe and Shutdown

## [0.1.0] - 2023-05-05

### Added

- Initial implementation of http