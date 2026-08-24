# Emergency fork ownership

This repository temporarily carries the `golang-migrate/migrate` patch required
by boatsetter/platform#2515.

- Owner: `@tgenov`
- Reviewing team: `@boatsetter/cb-platform`
- Upstream base: `18966c755f83de2b250dd8e766bad8b769cd0e25`
- Reviewed implementation merge: `fdbd6bf5e56ecad7db2538a8d4d60915109d942a`
- Reviewed implementation tree: `7caf0a3085ed12474b4641ca6a37b14fa80156a4`
- Review PR: https://github.com/boatsetter/migrate/pull/1
- Tracking issue: https://github.com/boatsetter/platform/issues/2515
- Expiry: 2026-09-24

The fork may only contain the reviewed GitHub source-driver dependency update,
its regression tests, and repository governance or CI metadata. It must not
accumulate unrelated product changes.

Remove the platform dependency on this fork at the first of:

1. a canonical `golang-migrate/migrate` release containing the upstream fix; or
2. removal of runtime `github://` migration fetching under RFC-083/#2403.

If neither exit is available by the expiry date, the owner must obtain and
record an explicit extension in boatsetter/platform#2515.
