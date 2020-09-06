# pbtree

`go get` for Protobuf - the missing link in your pb toolchain.

`pbtree` downloads dependencies for your protofiles and organizes them into a
single directory tree for later building.

This tool should be used before `protoc` and any proto linters.


# Table of Contents

1.  [pbtree](#org27a4941)
2.  [Installation](#org823a6e4)
3.  [Quick start](#orgb1c1a67)
4.  [Pulling protofiles from other projects](#org89f5965)
5.  [Canonical import format](#orgecf611a)
6.  [Full documentation](#orgc96ef68)


<a id="org280a530"></a>

## Problem and approach

Protobuf files use C-style imports - that is, any relative import will work if
`protoc` is able to find it under any of include paths (passed via `-I` flag).
That makes the building process too brittle - you need to supply a lot of -I&rsquo;s
and clone all the repos to proper paths.

`pbtree` attempts to solve that problem by rewriting all imports to a single
URI-like import style resembling Go&rsquo;s import rules; the format looks like
`git.corp/my/repo!/path/file.proto`. Afterwards, it downloads dependencies and
organizes them into a single file tree.


<a id="org7011c23"></a>

## State of a project

This is a first public release, but it is already used in production.
`pbtree` covers most of the usecases of the companies I&rsquo;ve been working with.

If it doesn&rsquo;t work for you - create an issue and we&rsquo;ll get there :)

TODO:

-   [ ] recursive versioning of dependencies, [#1](https://github.com/utrack/pbtree/issues/1)
-   [ ] global cache for dependencies pulled via HTTP
-   [ ] &#x2026;


<a id="org823a6e4"></a>

# Installation

Grab the latest release [here](https://github.com/utrack/pbtree/releases), or fetch it via `go` tool:

    GO111MODULE=on go get github.com/utrack/pbtree@latest


<a id="orgb1c1a67"></a>

# Quick start

Navigate to your repository and create a new pbtree project, passing full repo name via `--module`:

    » mkdir super-project && cd super-project
    » pbtree init github.com/me/super-project
    2020/07/16 17:45:07 new config is ready at '.pbtree.yml', edit away or see 'pbtree help add'

Add directory with protofiles:

    » mkdir protos && curl -o protos/foo.proto "https://gist.githubusercontent.com/utrack/0cac21b0ca1fafb96ef82afe15418037/raw/5ae54db359036736deaf020de8f205154fa57eaa/foo.proto"
    » pbtree add ./protos

Build your proto file tree:

    » pbtree build
    fetcher: using http fetcher for 'github.com/google/protobuf'

Check the generated tree:

    » tree -a
    .
    ├── .pbtree.yml
    ├── protos
    │   └── foo.proto
    └── vendor.pbtree
        └── github.com
            ├── google
            │   └── protobuf!
            │       └── src
            │           └── google
            │               └── protobuf
            │                   └── timestamp.proto
            └── me
                └── super-project!
                    └── protos
                        └── foo.proto

Now, you can use `protoc` to generate your proto(s) with a single command:

    » protoc -I./vendor.pbtree --go_out=. ./vendor.pbtree/github.com/me/super-project\!/protos/foo.proto

Sometimes `pbtree` won&rsquo;t be able to figure out your imports; you can either
change them to Canonical import format by hand or [Map your imports](https://github.com/utrack/pbtree/wiki/Rewriting-import-paths) without
changing files themselves.


<a id="org89f5965"></a>

# Pulling protofiles from other projects

You can pull any 3rd-party protofiles to generate clients for services
described in other repos. To pull them, use `pbtree get`:

    » pbtree get 'github.com/googleapis/googleapis!/google/type/datetime.proto'
    fetcher: using http fetcher for 'github.com/googleapis/googleapis'
    INFO    file successfully added, don't forget to call 'pbtree build'!
    
    » pbtree build
    fetcher: using http fetcher for 'github.com/googleapis/googleapis'
    fetcher: using http fetcher for 'github.com/google/protobuf'
    » tree -al
    .
    ├── .pbtree.yaml
    ├── protos
    │   └── foo.proto
    └── vendor.pbtree
        └── github.com
            ├── google
            │   └── protobuf!
            │       └── src
            │           └── google
            │               └── protobuf
            │                   ├── duration.proto
            │                   └── timestamp.proto
            ├── googleapis
            │   └── googleapis!
            │       └── google
            │           └── type
            │               └── datetime.proto
            └── me
                └── super-project!
                    └── protos
                        └── foo.proto


<a id="orgecf611a"></a>

# Canonical import format

Any import is converted to a single format that looks like
`repository!/dir/file.proto`, e.g.:

    git.enterprise.com/my/project!/api/file.proto
    github.com/google/protobuf!/src/google/protobuf/timestamp.proto

etc.

You can use this import format in your protofiles:

    syntax = "proto3";
    
    import "foo.com/bar/baz!/dir/file.proto"

For any other formats, `pbtree` will try to guesstimate what you want. See
[ImportPathDiscovery](https://github.com/utrack/pbtree/wiki/Import-path-discovery) for detailed info.


<a id="orgc96ef68"></a>

# Full documentation

Check out the [wiki](https://github.com/utrack/pbtree/wiki) for more.

