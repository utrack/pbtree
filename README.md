`go get` for Protobuf - the missing link in your pb toolchain.

`pbtree` downloads dependencies for your protofiles and organizes them into a
single directory tree for later building.

This tool is intended to be used before `protoc` and any proto linters.

Protobuf files use C-style imports - that is, any relative import will work if
`protoc` is able to find it under any of include paths (passed via `-I` flag).
That makes the building process very brittle - you need to supply a lot of -I&rsquo;s
and clone all the repos to proper paths.

`pbtree` attempts to solve that problem by rewriting all imports to a single
URI-like import style resembling Go&rsquo;s import rules. Every discovered file is
pulled, vendored and scanned automatically.

[Full documentation](https://github.com/utrack/pbtree/wiki)

# Table of Contents

1.  [Installation](#org3f80848)
2.  [Quick start](#orge5a734e)
3.  [Pulling protofiles from other projects](#orgf10ee9f)
4.  [Canonical import format](#orgab5551d)
5.  [Full documentation](#orgf1b4a16)


<a id="org3f80848"></a>

# Installation

Grab the latest release [here](https://github.com/utrack/pbtree/releases), or fetch it via `go` tool:

    GO111MODULE=on go get github.com/utrack/pbtree@latest


<a id="orge5a734e"></a>

# Quick start

Navigate to your repository and create a new pbtree project, passing full repo name via `--module`:

    » mkdir super-project && cd super-project
    » pbtree init --module github.com/me/super-project
    2020/07/16 17:45:07 new config is ready at '.pbtree.yml', edit away or see 'pbtree help add'

Add directory with protofiles:

    » mkdir protos && curl -o protos/foo.proto "https://gist.githubusercontent.com/utrack/0cac21b0ca1fafb96ef82afe15418037/raw/5ae54db359036736deaf020de8f205154fa57eaa/foo.proto"
    » pbtree add ./protos

Build your protofile tree:

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
change them to Canonical import format or [Map your imports](https://github.com/utrack/pbtree/wiki/Rewriting-import-paths) via pbtree configs.


<a id="orgf10ee9f"></a>

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


<a id="orgab5551d"></a>

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


<a id="orgf1b4a16"></a>

# Full documentation

Check out the [wiki](https://github.com/utrack/pbtree/wiki) for more.

