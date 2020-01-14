# godiff
A File/Directory diff-like comparison tool with HTML output.
Enhanced to support CSV files with Primary/Combinational Keys

This program can be use to compare files and directories for differences.
When comparing directories, it iterates through all files in both directories
and compare files having the same name.

See example output [here:](http://raw.githack.com/spcau/godiff/master/example.html)

## How to use godiff - general all files

	godiff file1 file2 > results.html
	godiff directory1 directory > results.html

## How to use godiff - csv files

	* Compare csv files with primary key
	godiff -key <CaseSensitive-Column-name> file1 file2
	* Compare csv files with combinational keys
	 godiff -key <CommaSeperatedCaseSensitive-Column-names> file1 file2
	* Output diff files to different name / folder
	 godiff -csv <diff-csv-file-name> -html <diff-html-file-name> -diff-dir <output-dir> -key <column-name> file1 file2
	* Measure the time taken to generate diff files
	 godiff -timeit -key <Column-name> file1 file2
See `godiff -h` for all the available command line options

## Features

* When comparing two directory, place all the differences into a single html file.
* Supports UTF8 file.
* Show differences within a line
* Options for ignore case, white spaces compare, blank lines etc.
* Compare csv files and generate diff csv file
* Compare csv files with single / combinational primary keys
* CSV files Columns / Rows can be any order
* Measure time taken to create diff files
* Diff files can be saved in different folder


## Description

I need a program to to compare 2 directories, and report differences in all
files. Much like gnudiff, but with a nicer output. And I also like to try out
the go programming language, so I created __godiff__.

The _diff_ algorithm implemented here is based on 
_"An O(ND) Difference Algorithm and its Variations"_
by Eugene Myers Algorithmica Vol. 1 No. 2, 1986, p 251. 

__godiff__ always tries to produce the minimal differences, 
just like gnudiff with the "-d" option.

## Go Language

This program is created in the go programming language.
For more information about _go_, see [golang.org](http://golang.org)

## How to Build

On Linux or Darwin OS

	go build -o godiff godiff.go godiff_unix.go

On Windows

	go build -o godiff.exe  godiff.go godiff_windows.go

## Prebuild Binary for Windows

https://github.com/rsrini7/godiff/releases/download/snapshot/godiff3.7z
