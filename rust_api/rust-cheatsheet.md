```markdown
# Rust Cheatsheet: The Essential 20%

This cheatsheet covers the most commonly used Rust concepts that will let you write ~80% of everyday code. It assumes you have a basic understanding of programming but are new to Rust. Each section includes concise explanations and small examples.

---

## 1. Variables and Mutability

```rust
let x = 5;          // immutable by default
let mut y = 10;     // mutable with `mut`
y += 5;             // OK

const MAX: u32 = 100;   // constant, type must be annotated

// Shadowing: reuse variable name
let x = x + 1;      // new variable shadows old one
```

- Variables are **immutable** by default – use `mut` to make them mutable.
- **Shadowing** lets you redeclare a variable with the same name, possibly changing its type.
- **Constants** are always immutable, must have a type annotation, and can be declared in any scope.

---

## 2. Basic Data Types

### Scalar Types

```rust
let a: i32 = -42;       // signed integer (i8, i16, i32, i64, i128, isize)
let b: u32 = 42;        // unsigned integer (u8, u16, u32, u64, u128, usize)
let c: f64 = 3.14;      // float (f32, f64)
let d: bool = true;     // boolean
let e: char = '🦀';     // character (4 bytes, Unicode)
```

### Compound Types

```rust
// Tuple
let tup: (i32, f64, char) = (42, 3.14, 'c');
let (x, y, z) = tup;                // destructuring
let first = tup.0;                   // index access

// Array (fixed length, same type)
let arr: [i32; 3] = [1, 2, 3];
let slice = &arr[1..3];              // slice: reference to part of array
```

- Arrays are stack-allocated; slices are **views** into sequences (arrays, Vecs, etc.).
- Use `Vec` for dynamic arrays (see Collections section).

---

## 3. Functions

```rust
fn add(x: i32, y: i32) -> i32 {
    x + y               // last expression is returned (no semicolon)
}

fn print_sum(x: i32, y: i32) {
    println!("{}", x + y);  // no return value (unit type `()`)
}
```

- Function parameters **must** have type annotations.
- Return type is specified with `->`. The last expression is implicitly returned; use `return` for early returns.
- `println!` is a macro (note the `!`) for printing to stdout.

---

## 4. Control Flow

### `if` expressions

```rust
let number = 7;
if number < 5 {
    println!("small");
} else if number == 7 {
    println!("lucky");
} else {
    println!("large");
}

// `if` can be used in a let statement
let result = if number > 5 { "big" } else { "small" };
```

### Loops

```rust
// loop (infinite until break)
let mut counter = 0;
let result = loop {
    counter += 1;
    if counter == 10 {
        break counter * 2;  // break with value
    }
};

// while
while counter < 20 {
    counter += 1;
}

// for over ranges
for i in 0..5 {          // exclusive range 0,1,2,3,4
    println!("{}", i);
}

for i in 0..=5 {         // inclusive range 0..=5
    println!("{}", i);
}

// for over collections
let vec = vec![10, 20, 30];
for item in vec.iter() {      // borrow each element
    println!("{}", item);
}
```

---

## 5. Ownership and Borrowing

**Ownership rules:**

1. Each value has a single **owner**.
2. When the owner goes out of scope, the value is dropped.
3. Ownership can be moved.

```rust
let s1 = String::from("hello");
let s2 = s1;            // s1 is MOVED to s2, s1 is no longer valid
// println!("{}", s1);   // ERROR: borrow of moved value

let s3 = s2.clone();    // deep copy (expensive)
println!("{}", s2);     // still works
```

**Borrowing:** references without taking ownership.

```rust
fn calculate_length(s: &String) -> usize {   // immutable reference
    s.len()
} // s goes out of scope but nothing is dropped

fn change(s: &mut String) {                  // mutable reference
    s.push_str(" world");
}

let mut s = String::from("hello");
let r1 = &s;                                 // multiple immutable borrows OK
let r2 = &s;
// let r3 = &mut s;                           // ERROR: cannot borrow as mutable because immutable borrows exist
println!("{} {}", r1, r2);

let r3 = &mut s;                              // after r1, r2 go out of scope, this is fine
```

- At any time, you can have **either** one mutable reference **or** any number of immutable references.
- References must always be valid (no dangling references).

---

## 6. Structs

```rust
// Define a struct
struct User {
    username: String,
    email: String,
    active: bool,
}

// Create an instance
let mut user1 = User {
    username: String::from("alice"),
    email: String::from("alice@example.com"),
    active: true,
};

// Access fields
user1.email = String::from("new@example.com");

// Field init shorthand
fn build_user(username: String, email: String) -> User {
    User {
        username,    // same as username: username
        email,
        active: true,
    }
}

// Tuple structs
struct Color(i32, i32, i32);
let black = Color(0, 0, 0);

// Unit struct (rarely used)
struct AlwaysEqual;
```

- Structs can have methods and associated functions using `impl` blocks.

```rust
impl User {
    // Associated function (static method)
    fn new(username: String, email: String) -> User {
        User { username, email, active: true }
    }

    // Method (takes self reference)
    fn deactivate(&mut self) {
        self.active = false;
    }

    // Getter (immutable self)
    fn email(&self) -> &String {
        &self.email
    }
}

let mut user = User::new(String::from("bob"), String::from("bob@example.com"));
user.deactivate();
```

---

## 7. Enums and Pattern Matching

### Enum definition

```rust
enum Message {
    Quit,
    Move { x: i32, y: i32 },   // struct-like variant
    Write(String),             // tuple variant
    ChangeColor(i32, i32, i32),
}
```

### `match` (exhaustive)

```rust
let msg = Message::Write(String::from("hello"));

match msg {
    Message::Quit => println!("Quit"),
    Message::Move { x, y } => println!("Move to ({}, {})", x, y),
    Message::Write(text) => println!("Text: {}", text),
    Message::ChangeColor(r, g, b) => println!("Color: {} {} {}", r, g, b),
}
```

### `if let` (when only one pattern matters)

```rust
if let Message::Write(text) = msg {
    println!("Message: {}", text);
} else {
    // optional else
}
```

### Common enums: `Option` and `Result`

```rust
enum Option<T> {
    Some(T),
    None,
}

enum Result<T, E> {
    Ok(T),
    Err(E),
}
```

- `Option` is used when a value may be absent.
- `Result` is used for operations that may fail.

```rust
fn divide(numerator: f64, denominator: f64) -> Option<f64> {
    if denominator == 0.0 {
        None
    } else {
        Some(numerator / denominator)
    }
}

match divide(10.0, 2.0) {
    Some(result) => println!("Result: {}", result),
    None => println!("Cannot divide by zero"),
}
```

---

## 8. Collections

### `Vec<T>` – dynamic array

```rust
let mut v: Vec<i32> = Vec::new();
v.push(5);
v.push(6);
v.push(7);

// Macro to create with values
let v2 = vec![1, 2, 3];

// Access
let third = &v2[2];               // panics if out of bounds
let third = v2.get(2);             // returns Option<&i32>

// Iterate
for i in &v2 {
    println!("{}", i);
}

// Mutable iteration
for i in &mut v {
    *i += 50;
}
```

### `HashMap<K, V>` – key-value store

```rust
use std::collections::HashMap;

let mut scores = HashMap::new();
scores.insert(String::from("Blue"), 10);
scores.insert(String::from("Yellow"), 50);

// Access
let team = String::from("Blue");
let score = scores.get(&team);          // Option<&i32>

// Iterate
for (key, value) in &scores {
    println!("{}: {}", key, value);
}

// Insert if key not present
scores.entry(String::from("Blue")).or_insert(25);
```

---

## 9. Error Handling

### Unrecoverable errors with `panic!`

```rust
panic!("crash and burn");   // program exits with message
```

### Recoverable errors with `Result`

```rust
use std::fs::File;
use std::io::ErrorKind;

let f = File::open("hello.txt");
let f = match f {
    Ok(file) => file,
    Err(error) => match error.kind() {
        ErrorKind::NotFound => match File::create("hello.txt") {
            Ok(fc) => fc,
            Err(e) => panic!("Problem creating file: {:?}", e),
        },
        other_error => panic!("Problem opening file: {:?}", other_error),
    },
};
```

### Propagating errors (shortcut `?`)

```rust
use std::fs::File;
use std::io::{self, Read};

fn read_username_from_file() -> Result<String, io::Error> {
    let mut f = File::open("hello.txt")?;   // if Err, returns early
    let mut s = String::new();
    f.read_to_string(&mut s)?;               // same here
    Ok(s)
}

// Even shorter chaining
fn read_username() -> Result<String, io::Error> {
    let mut s = String::new();
    File::open("hello.txt")?.read_to_string(&mut s)?;
    Ok(s)
}
```

- The `?` operator can only be used in functions that return `Result` or `Option` (or another compatible type).

### Common combinators

```rust
let result: Result<i32, &str> = Ok(5);
let doubled = result.map(|x| x * 2);               // Ok(10)
let checked = result.and_then(|x| if x > 0 { Ok(x) } else { Err("negative") });

let option: Option<i32> = Some(5);
let doubled_opt = option.map(|x| x * 2);            // Some(10)
```

---

## 10. Traits (Interfaces)

Traits define shared behavior.

```rust
// Define a trait
trait Summary {
    fn summarize(&self) -> String;

    // default implementation
    fn default_summary(&self) -> String {
        String::from("(Read more...)")
    }
}

// Implement for a type
struct Article {
    title: String,
    author: String,
}

impl Summary for Article {
    fn summarize(&self) -> String {
        format!("{} by {}", self.title, self.author)
    }
}

// Use trait as parameter
fn notify(item: &impl Summary) {
    println!("Breaking news! {}", item.summarize());
}

// Trait bound syntax
fn notify<T: Summary>(item: &T) {
    println!("{}", item.summarize());
}

// Multiple bounds
fn some_function<T: Summary + Display>(item: &T) {}
```

- Traits can have **associated types** and **generic parameters**.
- Derivable traits: `Debug`, `Clone`, `Copy`, `PartialEq`, `Eq`, etc.

```rust
#[derive(Debug, Clone, PartialEq)]
struct Point {
    x: i32,
    y: i32,
}
```

---

## 11. Common Macros

```rust
println!("Hello, {}!", "world");    // print with newline
print!("no newline");                // print without newline
eprintln!("Error: {}", msg);          // print to stderr

format!("{}-{}", "hello", 42);        // returns String

panic!("something went wrong");       // crash with message

assert!(true);                         // panic if false
assert_eq!(5, 5);                      // panic if not equal

vec![1, 2, 3];                          // create a Vec
```

- Macros are indicated by `!` – they can take a variable number of arguments.

---

## 12. Modules and Use

```rust
// In lib.rs or main.rs
mod math {
    pub fn add(a: i32, b: i32) -> i32 {   // `pub` makes it public
        a + b
    }

    mod inner {                            // private by default
        pub fn multiply(a: i32, b: i32) -> i32 {
            a * b
        }
    }
}

// Bring into scope
use math::add;
use math::inner::multiply;                 // requires `pub` on inner too

// Nested imports
use std::{collections::HashMap, io::{self, Write}};

// External crate (in Cargo.toml)
use serde_json::json;
```

- Files and directories also define modules (e.g., `mod foo;` looks for `foo.rs` or `foo/mod.rs`).

---

## 13. Common Patterns

### `if let` for concise matching

```rust
let config_max = Some(3u8);
if let Some(max) = config_max {
    println!("Maximum is {}", max);
}
```

### `while let` for repeated matching

```rust
let mut stack = Vec::new();
stack.push(1);
stack.push(2);
while let Some(top) = stack.pop() {
    println!("{}", top);
}
```

### Closures (anonymous functions)

```rust
let add_one = |x| x + 1;
let result = add_one(5);      // 6

// With type annotations
let expensive_closure = |num: u32| -> u32 {
    println!("calculating...");
    num * 2
};

// Closures can capture variables from environment
let x = 10;
let equal_to_x = |z| z == x;
```

### Iterators and combinators

```rust
let numbers = vec![1, 2, 3, 4, 5];
let doubled: Vec<i32> = numbers.iter().map(|x| x * 2).collect();   // [2,4,6,8,10]
let even: Vec<&i32> = numbers.iter().filter(|&&x| x % 2 == 0).collect(); // [2,4]
let sum: i32 = numbers.iter().sum();                               // 15
```

- Iterators are **lazy** – they do nothing until consumed (e.g., by `collect`, `sum`, `for` loop).

---

## 14. Testing

```rust
#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        let result = 2 + 2;
        assert_eq!(result, 4);
    }

    #[test]
    #[should_panic(expected = "divide by zero")]
    fn test_panic() {
        panic!("divide by zero");
    }
}
```

- Run tests with `cargo test`.

---

## Summary

These are the building blocks you'll use in most Rust programs. Remember:

- **Ownership & borrowing** are unique to Rust – they ensure memory safety without a garbage collector.
- **Pattern matching** (`match`, `if let`) is powerful and concise.
- **Traits** provide polymorphism.
- **Error handling** with `Result` and `Option` encourages explicit handling of fallibility.

For more depth, refer to [The Rust Book](https://doc.rust-lang.org/book/). Happy coding! 🦀

```
