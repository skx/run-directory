# run-directory

This project is a naive port of `run-parts` to golang.

## About run-parts

`run-parts` executes every (executable) file in a directory, in order,
it is usually shipped with a system cron-package where it is used, but it
is also a useful command which users can use in their own shell-scripts.

In the case of `cron` it is often used to execute a series of scripts,
for example you might have a directory of files to be run every hour, and
that could be achieved via :

    $ run-parts /etc/cron.hourly/


# Motivation

The Debian version of run-parts allows processing to terminate if one
of the scripts fails, but unfortunately the CentOS version of `run-parts`
doesn't support this ability.


# Installation

There are two ways to install this project from source, which depend on the version of the [go](https://golang.org/) version you're using.

If you prefer you can fetch a binary from [our release page](https://github.com/skx/run-directory/releases).

## Build without Go Modules (Go before 1.11)

    go get -u github.com/skx/run-directory

## Build with Go Modules (Go 1.11 or higher)

    git clone https://github.com/skx/run-directory ;# make sure to clone outside of GOPATH
    cd run-directory
    go install


# Github Setup

This repository is configured to run tests upon every commit, and when
pull-requests are created/updated.  The testing is carried out via
[.github/run-tests.sh](.github/run-tests.sh) which is used by the
[github-action-tester](https://github.com/skx/github-action-tester) action.

Releases are automated in a similar fashion via [.github/build](.github/build),
and the [github-action-publish-binaries](https://github.com/skx/github-action-publish-binaries) action.



Steve
