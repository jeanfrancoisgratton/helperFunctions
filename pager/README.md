# The pager subpackage
A UNIX-like paging utility, pretty much like "more"
___

This package provides functions that I have been using all over my various tools but the obvious struck: instead of ducplicating code on my filesystem, I should have a single source of truth.

## Overview

You can use this package in two ways:
- Paginate a text file
- Paginate a provided text block

Whichever the source (text file, text block), you can also provide a terminal height value; this means that paginated pages can be smaller than the actual terminal height.

## Installing the package

As simple as `go get github.com/jeanfrancoisgratton/helperFunctions/pager` and adding it to the import statements in your software.

## How to use

As mentioned, the source can be either a text file that gets read, or a provided text block

### Paginating a text file
All you need to do to paginate a text file is to call `PaginateFromFile(file, pageHeight)` where:
- file is the path to the file to be read and paginated
- pageHeight is the desired page height. A value of zero means you'd use the current terminal's height

Note that a value below 10, or a terminal height below 10 (who uses a terminal less than 10 lines tall anyway ??) will prevent paging; the file will be displayed as is.

This function, `PaginateFromFile` is actually a wrapper around the whole Pager structure (the NewPager() function, described below). The text file is converted into a text block and we then paginate the whole thing with NewPager()

### Paginating a text block
You call NewPager(_textblock_, _pageHeight_), where :
- _textblock_ is the actual block to be paginated
- _pageHeight_ is the desired page height. A value of zero means you'd use the current terminal's height


## Extra functionalities

Once you've reach the end of a page the pager will await your next action (a key press), which can be one of the following:

- spacebar (whitespace) : scroll down to the next page
- u : scroll up one page (previous page)
- \n (ENTER) : scroll down a single line
- y : scroll up one single line
- q : abort paging and return to the calling function
- CTRL+C : quit the software right away
