# os2ds-pdf-PoC

## What is this?

This repository contain a Proof of Concept regarding the decision to
rewrite parts of os2datascanner's engine in a new language with the goal
of gaining a considerable speed-up in performance.

One of the slowest engine components is the pdf-scanner.
So, in this PoC, we will benchmark `python` (with `PyPDF2`) against
alternatives.

## Languages included in the PoC

Choosing a language for a project is not at all trivial.
The main goal is a gain in performance, but there are other factors to
consider such as ease-of-use, developer experience and workforce availability.
Speed isn't the only valid concern.

To narrow the field, we consider the fastest languages according to 
"The Computer Language Benchmarks Game":

![How many times more CPU seconds](https://benchmarksgame-team.pages.debian.net/benchmarksgame/download/fastest.svg)


At the time of writting some of the fastest languages are:

- C++ (g++)
- Rust
- C# .NET
- Julia
- Haskell (ghc)
- OCaml
- Go

Since the main goal of the PoC is reading and analyzing pdf-files, which in general involves
string and text manipulation, the language has to be well-suited for such tasks.

That is, languages that have poor support for string and text manipulation have already been eliminated from the list
(C, Fortran, etc.).

We prefer languages, whose reference implementation can produce an executable binary for native targets, since this eliminates
overhead (i.e. there is no need to ship the runtime alongside the code).

Therefore, C# .NET is eliminated. This is also due to the fact that C# .NET is a Microsoft controlled platform.

C++ may be one of the fastest languages, however it is also the most complicated and difficult language to learn and use.
A famous quote from the inventor of C++, Bjarne Stroustrup: "C makes it easy to shoot yourself in the foot; C++ makes it harder, 
but when you do it blows your whole leg off".

For this reason and the fact that finding (good) C++ developers is extremely difficult.

As for the remaining languages: Rust, Julia, Haskell, OCaml and Go, developers are equally difficult to find
since the majority of these languages are fairly new (or have only gained popularity in the last couple of years).
This means that the ease of learning a language becomes a factor.

Regarding this difficulty, we have ranked the contenders from easiest (1) to most difficult (5) to learn
from the perspective of a python developer.

- 1: Julia
- 2: Go
- 3: Haskell
- 4: OCaml
- 5: Rust

The syntax of Julia is quite similar to that of python, which makes it a great candidate.
Julia is garbage-collected and is JIT-compiled to LLVM IR, which makes it both fast and
easy to use.

Go is intended to be a successor to C, and aims to be a very minimalistic language.
It only has 26 reserved keywords, so the barrier to entry is fairly low, aside from the
fact that it has pointers. Go is compiled to native machine code though LLVM, is garbage-collected
and libraries are statically linked.

The thing about Haskell is that the language itself is not that difficult to learn.
But, the functional paradigm that haskell is built on may be a big stone that many stumble
over. However, for programmers that are used to the functional style, haskell code
can be very compact and elegant. The reference implementation of Haskell, 
the Glorious Glasgow Haskell Compiler (ghc), compiles code to native machine code 
though LLVM and is garbage-collected.

OCaml belongs to the ML-family of languages (ML stands for Meta Language) and
supports both imperative and declarative programming styles.

Rust may be the fastest out of the contenders, but is it also the most difficult to learn.
Why is that? Borrowing and Ownership (+ Lifetimes) is an alternative way of managing memory
compared to garbage-collection, and it requires almost as much knowledge about memory management
as C++. There are two modes: safe and unsafe. In safe-mode, the rust-compiler will check the code
for memory-allocation issues at compile-time and throw an error, i.e. it is impossible to compile
a memory-unsafe program in this mode.
