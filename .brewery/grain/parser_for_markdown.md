Parser for Markdown
===================

Description
-----------
We need to be able to parse a Markdown file and pull out the proper properties.

 - [x] H1 equivalent block (both the *Setext* and *atx* styles) is the `title` property
 - [x] H2 equivalent blocks (both the *Setext* and *atx* styles) are property names
 - [x] Blocks within H2 are property values

Comments
--------
Different regular expressions:

 * H1 inline - /^#{1}\s*\w.+$/gm
 * H1 line   - /^\w.+\r?\n={2,}$/gm
 * H2 inline - /^#{2}\s*\w.+$/gm
 * H2 line   - /^\w.+\r?\n-{2,}$/gm

Status
------
Testing

Assigned
--------
@bmallred

Created
-------
2016-08-24

Modified
--------
2016-08-30

Tags
----
 - issue
 - parsing
 