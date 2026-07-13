###
```bash
#Reference:
https://www.educative.io/blog/50-golang-interview-questions
```

## What are the benefits of using Go compared to other languages?
1. Unlike other languages which started as academic experiments, Go code is pragmatically designed. Every feature and syntax decision is engineered to make life easier for the programmer.

2. Golang is optimized for concurrency and works well at scale.

3. Golang is often considered more readable than other languages due to a single standard code format.

4. Automatic garbage collection is notably more efficient than Java or Python because it executes concurrently alongside the program.

## What are string literals?
A string literal is a string constant formed by concatenating characters. The two forms of string literal are raw and interpreted string literals.

Raw string literals are written within backticks (foo) and are filled with uninterpreted UTF-8 characters. Interpreted string literals are what we commonly think of as strings, written within double quotes and containing any character except newline and unfinished double quotes.

## What data types does Golang use?#
Golang uses the following types:

1. Method
1. Boolean
1. Numeric
1. String
1. Array
1. Slice
1. Struct
1. Pointer
1. Function
1. Interface
1. Map
1. Channel

## What are packages in a Go program?
Packages (pkg) are directories within your Go workspace that contain Go source files or other packages. Every function, variable, and type from your source files are stored in the linked package. Every Go source file belongs to a package, which is declared at the top of the file using:

```go
package <packagename>
```
You can import and export packages to reuse exported functions or types using:
```go
import <packagename>
```
Golang’s standard package is `fmt`, which contains formatting and printing functionalities like `Println()`.

## What form of type conversion does Go support? Convert an integer to a float.
Go supports explicit type conversion to satisfy its strict typing requirements.

```go
i := 55      //int

j := 67.8    //float64

sum := i + int(j) //j is converted to int
```

## What is a goroutine? How do you stop it?#
A goroutine is a function or method that executes concurrently alongside any other goroutines using a special goroutine thread. Goroutine threads are more lightweight than standard threads, with most Golang programs using thousands of goroutines at once.

To create a goroutine, add the go keyword before the function declaration.

```go
go f(x, y, z)
```

You can stop a goroutine by sending it a signal channel. Goroutines can only respond to signals if told to check, so you’ll need to include checks in logical places such as at the top of your for loop.

```go
package main
func main() {
  quit := make(chan bool)
  go func() {
    for {
        select {
        case <-quit:
            return
        default:
            // …
        }
  }
}()
// …
quit <- true
```