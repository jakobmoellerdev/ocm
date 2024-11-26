<!-- Release notes generated using configuration in .github/config/release.yml at refs/heads/releases/v0.19 -->

## What's Changed
### 🚀 Features
* feat(log): log http requests for OCI and docker based on trace level by injecting a logger by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1118
### 🧰 Maintenance
* chore: change guide for 0.18.0 by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1066
* chore: allow publishing to Brew via custom script by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1059
* chore: remove ocm inception during build CTF aggregation by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1065
* chore: release branches as `releases/vX.Y` instead of `releases/vX.Y.Z` by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1071
* chore: cleanup release action by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1076
* chore: bump version to 0.19.0-dev by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1084
* chore: disable mandatory period comments by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1079
* chore: make sure that version bumping happens everytime by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1090
* chore: directly reference integration tests as job by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1096
* chore: also create a branch based on the tag to avoid dangling commits by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1098
* chore: add correct labels for flake nix job by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1100
* chore: allow triggering blackduck scans manually by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1104
* chore: remove the int test repository dispatch by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1106
* chore: rework labeling jobs by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1103
* chore: ensure that PR titles must be semantic by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1108
* chore: move process options to struct by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1109
* chore: automatically set github actions label by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1112
* chore: remove releasenotes.yaml by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1111
* docs: document complex artifact transfer by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1113
* chore: setup release to reuse CTF from components workflow by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1077
* docs: revise RELEASE_PROCESS.md by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1086
* chore: label prs based on conventional commit by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1121
* docs: finally some working examples for a lot of commands by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1123
* chore: let's not store the release notes in the repository by @hilmarf in https://github.com/open-component-model/ocm/pull/1120
* chore: fixup release action versioning and notes process by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1124
* chore: make sure we release to brew too with our release by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1125
* chore(deps): bump anchore/sbom-action from 0.17.7 to 0.17.8 in the ci group by @dependabot in https://github.com/open-component-model/ocm/pull/1128
* chore(deps): bump the go group with 17 updates by @dependabot in https://github.com/open-component-model/ocm/pull/1127
* chore: publish to website as other by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1126
* chore(github_actions): push-to-winget: permissions: `{"contents":"write","pull_requests":"write"}` by @hilmarf in https://github.com/open-component-model/ocm/pull/1133
* chore(signing): correct Fulcio service to correct address by @morri-son in https://github.com/open-component-model/ocm/pull/1135
* chore(github_actions): using now classic secret of OCM_CI_ROBOT by @hilmarf in https://github.com/open-component-model/ocm/pull/1137
* chore: rework release note handling by @jakobmoellerdev in https://github.com/open-component-model/ocm/pull/1139


**Full Changelog**: https://github.com/open-component-model/ocm/compare/v0.18...v0.19.0