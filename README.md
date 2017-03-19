# deep

A shallow package manager for Go

## Warning

**This is currently under heavy development.**

Until you see tests for it and this warning changed / removed you should treat
it like a nuclear bomb for your packages. Don't touch it unless you know what
you are doing and are prepared to take all the blame if something goes horribly
wrong (which it will).

## Current status

This can almost vendor itself but it won't read the lock/manifest files when
running again and it will cause issues when importing the dependencies as again.
More to come in the next days.

## Project goals

Have a package manager for Go projects that is easy to use and provides the
common needed features across a wide enough range of use-cases.

## Steps to implement this

- [x] decide to follow or extend the manifest/lock files of [golang/dep](https://github.com/golang/dep)
  - decision taken, this won't follow the golang/dep manifest/lock files but
it will have a compatibility mode with it
- [x] decide on core functionality
  - core functionality will be: adding dependencies, updating and removing
them. Versions will be optional for each of these operations.
- [ ] implement basic tool which satisfies the core functionality
- [ ] release it into the wild
- [ ] collect and evolve based on the users feedback


## Ideas to filter out

- will work only with semver semantics
- usage:
```shell
deep (optional package@version)
deep update package (optional @version)
deep rm <package>
```
- vendor flattening will be a thing done by default, no questions asked
- stripping of things: ` --strip=tests,examples,main,all `, with ` all ` as
default
- alternative to stripping, `  --keep=none,tests,main,all ` with
` none ` as default
- the manifest will dictate the list of OS/arch combos, otherwise the current
one is assumed
- test dependencies of vendored libs are not downloaded by default, something
should be done about it
- humans will have to sort out conflicts for versions:
  - the major version you want conflicts with the version a lib wants? boom
  - the major version two libs want conflicts? boom


## Usage

To install the latest version use the known way to get package into your system

```bash
go get -u github.com/dlsniper/deep/cmd/deep
```

More usage to come as the project matures and gets functionality added.


## Rational for this tool

Based on the long and agonizing journey seen in the various tools that try to
be package managers for Go, including the latest incarnation of this,
[golang/dep](https://github.com/golang/dep), there's a clear need for something
that just works.

This tool does not solve every problem that can appear, does not attempt to
solve issues around licensing, CGO or any other complex scenarios. If you need
such a tool, probably [golang/dep](https://github.com/golang/dep) is the tool
you are looking for.

There's no reason why this will not evolve in the future to support such
use-cases but simplicity in both user interaction and maintenance will always
be ahead of any thing else.

## Contributing

I don't expect any contributors to this tool, or users for that matter, but if
you choose to use it find any issues with it or things that can improved, I
would like to ask you to please spare a few moments and either file a bug
report (as detailed as possible to help reproduce the issue) or a feature
request.

Should you want to send a PR to this tool, that would be amazing and I thank
you for that, I will do my best to help you get it merged asap. If you want to
change a large portion of the code base or you want to add a major piece of
functionality, please open up an issue first and lets talk about it, then try
to send modular PRs so they can be reviewed and merged asap (I can create
branches in the main repository should this help getting a large functionality
get merged).

Thank you.

## License

As with any other software, this too has a license. This project uses the
[Apache License v2.0](LICENSE.md) which can be found in the file LICENSE.md in
the root of this repository, where this file is also located, or at the
following address: http://www.apache.org/licenses/LICENSE-2.0

Copyright 2017 Florin Pățan

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
