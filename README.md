# About
This is a Golang package that provides nmap→sqlite functionality such as the following

## Types
This package provides a type file in *pkg/nmap/types.go* that contains a comprehensive list of nmap types used in the process of unmarshalling nmap xml output. 

I chose to use a bespoke type file instead of relying on 3rd party packages providing such types for the following reasons:
1. Given the history of nmap, some elements and attributes don't exist anymore and should be removed from the type file.
2. Some nmap type packages had struct tags that were misspelled and would result in incorrect unmarshalling.
3. Even if the types provided were fine, some packages have a little extra functionality that I don't need when compiling the package.

By having my own type file, I have the freedom to remove the unnecessary xml elements and correct the struct tags.

## Incremental XML parser
This package includes a function that takes an `io.Reader` representing nmap xml output and can incrementally parse that data stream
on a per element basis.

There were a few reasons for this:
1. From what I have seen, there is no package that incrementally parses nmap xml output, and that is why I specifically created a function for this.
I don't see the need to parse the entire xml output at once. Nmap does a good job of immediately outputting what it has found: I think it's best to 
immediately parse this data so consumers can see immediate results instead of waiting for nmap to finish.
2. I originally was thinking of using [this tool to do this nmap→sqlite conversion](github.com/vdjagilev/nmap-formatter) for me. However, it is not an incremental parser.
Furthermore, even if it was -- or you didn't care about incremental parsing -- this would add a build step and additional complexity to whatever project was using the tool.
For example, assuming you wanted everything to be in Go, you would need to create an nmap `exec.Command`, run & wait for stdout to be closed, then create a nmap-formatter `exec.Command`
and pass the stdout of the first command to the stdin of the second command. Even if this was fine for you, you are still unable to store multiple scans, and
you don't have the capability to update parts of the database with new data (i.e. from incremental parsing).

## SQLite database
I used [sqlc](https://sqlc.dev) to mostly generate the sqlite database code to ensure type safety while automating the process of generating functions to work with the database.

By using sqlc, I can focus more on the queries used to interact with the database and less on the database code itself.
Additionally, when you use sqlc, you need a database schema to work with and this forces you to have a high level overview of the database structure.
It also makes it easier for people using my package to understand how to use the database.

I say mostly because there were a few caveats with sqlc that I needed to work around; the database schema file is what I think accurately represents nmap scan xml output and is normalized
to the best of by ability. However, given this schema, certain tables need access to the id of other tables. 
I could generate and reuse a UUID for new entries, but I feel like it would be a hassle when autoincrement can handle that for me, and UUIDs take up more space.
However, when autoincrementing, you need that id that it autoincremented to use for other inserts.
Sqlc doesn't really provide a way to do this (i.e using `sql.Result.LastInsertId()`) and I had to manually edit the generated code to add this functionality.
That is why .other.sqlc.yaml exists: if I messed up and needed to regenerate the code that specifically needs to return the last insert id,
I can just run `sqlc generate -f .other.sqlc.yaml` and it will regenerate the code for me, although I need to manually edit the generated code again.

<hr>

# Usage
I provided a Makefile to make it easier to manage, test and benchmark this package.

## Generating the database
Run `make build-db` to delete the existing database and generate a new one.

## Genereating the sqlc code
Run `make gen` to generate the main 'automatic' sqlc code.

`make gen-alt` is used to generate the code that needs to return the last insert id.
Sqlc doesn't really provide a way to do this, and you need to manually edit the generated code 
to add this functionality.

## Testing
There are two main tests you can do

`make test-xml` will test the xml parsing function on a comprehensive xml test file

`make test-nmap` will test the xml parsing function on a real nmap scan

## Benchmarking
You can benchmark the performance of the xml parsing function by running `make bench-xml`
