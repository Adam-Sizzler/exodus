# Embedded fonts

`jetbrains-mono-{regular,bold}.woff2` are subsets of [JetBrains Mono][jbm]
v2.304, licensed under the SIL Open Font License 1.1 — see `OFL.txt`. JetBrains
Mono declares no Reserved Font Name, so a subset may keep the name.

They are embedded in the SVG output rather than referenced by name. The report
is built almost entirely from box-drawing and block characters, and a viewer
without a font covering them renders broken frames or tofu. Embedding also
fixes the advance width, which is what lets the SVG place a run of text by
column number instead of measuring it.

The subset covers what the report can emit: ASCII, Latin-1 and Latin Extended-A
for country names, box drawing, block elements, geometric shapes and arrows.
Regenerate with:

    scripts/subset-fonts.sh

[jbm]: https://github.com/JetBrains/JetBrainsMono
