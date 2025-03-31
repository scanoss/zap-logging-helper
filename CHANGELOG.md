# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Upcoming changes...

## [0.3.2] - 2025-03-31
### Added
- Added 'x-response-id' to response trailer
- Upgraded Go v1.24.x
- Upgraded golangci-lint to v1.64.8

## [0.3.1] - 2024-08-19
### Update
- Upgraded to GO 1.22
- Updated dependency versions

## [0.3.0] - 2023-11-20
### Added
- Added Open Telemetry (OTEL) trace support
- Upgraded to Go 1.20

## [0.2.0] - 2023-05-04
### Added
- Added support to configure log outputs (stdout/stderr/file)

## [0.1.0] - 2023-03-08
### Added
- Added support for Streaming interceptor
### Fixed
- Fixed issues as a result of `golangci` linting

## [0.0.2] - 2022-10-21
### Added
- License headers and documentation

## [0.0.1] - 2022-10-20
### Added
- Initialisation functions for zap
- Atomic logging level setting
- gRPC interceptor to inject logging context and request/response id

[0.0.1]: https://github.com/scanoss/zap-logging-helper/compare/v0.0.0...v0.0.1
[0.0.2]: https://github.com/scanoss/zap-logging-helper/compare/v0.0.1...v0.0.2
[0.1.0]: https://github.com/scanoss/zap-logging-helper/compare/v0.0.2...v0.1.0
[0.2.0]: https://github.com/scanoss/zap-logging-helper/compare/v0.1.0...v0.2.0
[0.3.0]: https://github.com/scanoss/zap-logging-helper/compare/v0.2.0...v0.3.0
[0.3.1]: https://github.com/scanoss/zap-logging-helper/compare/v0.3.0...v0.3.1
[0.3.2]: https://github.com/scanoss/zap-logging-helper/compare/v0.3.1...v0.3.2