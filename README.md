# PowerEdit

Project Gutenberg is a repository for thousands of books in the public domain which can be downloaded as .txt files .pdf's or viewed in the web.

These are often created by one to several proof-readers who carefully compare a facsimile against a text file that has been copied from that facsimile.

Inherent in this work is the possibility for mistakes. This can take the form of additions to the text, typos, and omissions from the text including whole missing phrases.

PowerEdit takes as input 2 files, the official Project Gutenberg text and another more authoritative (though not error free) text, and walks the user through every discrepancy between the texts,
prompting the user to decide what to do in each case, providing simple but powerful options.

Users can save there work and PowerEdit allows them to pick right back up at the next discrepancy.

### Note:

Because either file may contain regular, obvious discrepancies (one file has page numbers while the other does not, footnotes, randome symbols, etc) it is very useful to first clean each file up of obvious mistakes by using tools such as find-and-replace in conjunction with regex or other pattern matching. Otherwise the number of discrepancies that one must work though could be drastically increased.

## Install

From the root of the project run `make build`

This will create the `poweredit` executable which you can move into you PATH

## Usage

To start a new job pass two files to the command, a file to edit, and the more authoritative version:

`poweredit <text file to edit>  <text file to compare to>`

To resume a job,

View list of in progress jobs:
`poweredit jobs`

Copy the job to resume and:

`poweredit <name of job>`

Alternatively, provide the CSV which is tracking the job you wish to resume

`poweredit <a_jobfile.csv>`

## Editing options

At each discrepancy you will be prompted to resolve the discrepancy with one of the following options:

```
<xy|x> - enter two numbers to advance cursors: x for file under edit, y for source file;  a single digit entry will advance cursor for file under edit by x

a - to add missing token to file under edit
e - edit typo, sets current word of file under edit to current word of source file
ex - edit typo in source, sets current word of source file to current word of file under edit
me - manually enter a custom word set current token for file under edit and source file to this word
d - delete token from file under edit
x - delete current token from source file
v - save changes and quit
q - quit without saving any changes made
```