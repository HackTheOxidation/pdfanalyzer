# pdfanalyzer

## What is this?

This repository contain a Proof of Concept regarding the decision to
rewrite parts of os2datascanner's engine in a new language with the goal
of gaining a considerable speed-up in performance.

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
images, graphics and so on, for `os2datascanner` to analyzer. Another goal is to be able to
filter out objects of interest in order to reduce the amount of data that needs to be scanned.

### Roadmap

- 1: Extract objects from a pdf file
- 2: Filter extracted objects
- 3: Send filtered objects to RabbitMQ
