# BEL
Just another implementation.

# PG DOCS AND SOURCE
You can find documentation for `bel` from Paul Graham's web [site](http://www.paulgraham.com/bel.html).

For my convenience, I copied those documents to the `pgdocs/` directory.

The source for `bel` is copied from [paulgraham/bel.bel](https://sep.yimg.com/ty/cdn/paulgraham/bel.bel?t=1570993483&).
I have a copy saved in this repository as `pgdocs/bel.bel`.

I didn't find a license for the files on his site, so treat the files as his, all rights reserved.

# BUILDING BEL
The "Building Bel" documentation is heavily sources from the tutorial [Building Lisp](https://www.lwh.jp/lisp/index.html) by Leo Howell.

# HISTORY
Logic for `bel` is derived from `schism`, which is derived from `Mini-Scheme` by Atushi Moriwaki and Akira Kida, which is derived from Matsuda and Saigo.

# Stringer

    go:generate stringer -type opcode bel/opcode.go
