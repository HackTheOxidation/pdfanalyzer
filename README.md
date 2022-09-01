# pdfanalyzer

## What is this?

This repository contain a Proof of Concept regarding the decision to
rewrite parts of os2datascanner's engine in a new language with the goal
of gaining a considerable speed-up in terms of performance.

One of the slowest engine components is the pdf-scanner.
So, in this PoC, we will benchmark the current `python` implementation against
this PoC implementation written in `go`.

## Choice of language for the PoC

Choosing a language for a project is not at all trivial.
The main goal is a gain in performance, but there are other factors to
consider such as ease-of-use, developer experience and workforce availability.
Speed isn't the only valid concern.

Since the main goal of the PoC is reading and analyzing pdf-files, which in general involves
string and text manipulation, the language has to be well-suited for such tasks.

Go is intended to be a successor to C, and aims to be a very minimalistic language.
It only has 26 reserved keywords, so the barrier to entry is fairly low, aside from the
fact that it has pointers. However, the use of pointers are rather restricted compared to `c`,
which is a good thing. Go is compiled to native machine code, is garbage-collected
and libraries are statically linked by default, which is very desirable from a performance
perspective.

## Goals of the project

The Goal of the PoC is to make a tool that can extract various objects from a pdf file, including text
images, graphics and so on, for `os2datascanner` to analyze. Another goal is to be able to
filter out objects of interest in order to reduce the amount of data that needs to be scanned.

### Roadmap

The following 5 points defined below the core features of the project.

- 1: Extract objects from a pdf file
- 2: Filter extracted objects
- 3: Process extracted objects if necessary
- 4: Be used both as a CLI-tool and as a library
- 5: Send filtered objects to RabbitMQ

## Getting Started

Currently, you need to have `go` installed on your system to install, run and/or develop on the project.
Head over to [go.dev](https://go.dev/learn/) if you are new to `go`.

### Manual installation from repository with `go install`

If you want to install and use the `pdfanalyzer`-cli, just run the following:

```sh
git clone https://git.magenta.dk/os2datascanner/os2ds-pdf-poc.git
cd os2ds-pdf-poc
go install .
```

This will compile `pdfanalyzer` and place the binary executable in `$GOROOT/bin` (this is probably `$HOME/go/bin`).
If you haven't already, add `$GOROOT/bin` to `$PATH` or move the executable somewhere else. 

### For developers/contributers

If you want to develop on the project, just clone the repo with ssh and you are good to go (no pun intended):

```sh
git clone git@git.magenta.dk:os2datascanner/os2ds-pdf-poc.git
```

