// This is a simple Todo list web API built with the Actix-web framework in Rust.
// It demonstrates basic CRUD operations (Create, Read, Update, Delete) for Todo items.
// Each line is commented to explain what it does, especially for someone new to Rust.

// --- Imports: bringing external functionality into scope ---
use actix_web::{               // Actix-web is a web framework for Rust
    middleware::Logger,        // Logger middleware to log HTTP requests
    web,                       // Web module containing routing and data extraction utilities
    App,                       // App struct to build the web application
    HttpResponse,              // For constructing HTTP responses
    HttpServer,                // The HTTP server itself
    Responder,                 // Trait that allows many types to be converted to HTTP responses
};
use serde::{Deserialize, Serialize}; // Serde provides serialization/deserialization (to/from JSON)
use serde_json::json;          // Helper macro to create JSON objects easily
use std::sync::Mutex;          // Mutex allows safe shared access to data across threads

// --- Data structures ---

// The full representation of a Todo item. It will be used when returning data to the client.
// #[derive(...)] automatically implements common traits: Debug for printing, Serialize for JSON output, Clone for duplicating.
#[derive(Debug, Serialize, Clone)]
struct Todo {
    id: usize,                 // Unique identifier for the todo
    is_done: bool,              // Whether the task is completed
    name: String,               // Short name of the task
    description: String,        // Longer description (can be empty)
}

// This struct is used when the client sends a request to create a new todo.
// It does NOT include an id – the server will assign one.
// #[derive(Deserialize)] allows automatic parsing from JSON.
#[derive(Debug, Deserialize)]
struct CreateTodoRequest {
    is_done: bool,
    name: String,
    description: String,
}

// For simplicity, updating a todo uses the same fields as creating one.
// A type alias means UpdateTodoRequest is exactly the same as CreateTodoRequest.
type UpdateTodoRequest = CreateTodoRequest;

// The global state of our application, shared across all threads/requests.
struct AppState {
    // Mutex wraps the Vec so that only one thread can modify it at a time.
    // This is necessary because Actix-web runs multiple threads.
    todos: Mutex<Vec<Todo>>,   // The list of todos
    next_id: Mutex<usize>,     // The next ID to assign to a new todo
}

// --- Initialization ---
// Creates and returns a new AppState with some sample todos.
fn init_state() -> AppState {
    AppState {
        // Mutex::new creates a new mutex containing the given value.
        todos: Mutex::new(vec![
            Todo {
                id: 1,
                is_done: false,
                // .to_string() converts a string literal (&str) into a heap-allocated String.
                name: "Do the laundry".to_string(),
                description: "".to_string(),
            },
            Todo {
                id: 2,
                is_done: false,
                name: "Clean the dishes".to_string(),
                description: "".to_string(),
            },
        ]),
        // The next ID starts at 3 because we already used 1 and 2.
        next_id: Mutex::new(3),
    }
}

// --- Request handlers (each handles one API endpoint) ---

// Handler for GET /todos – returns all todos as JSON.
// `async` means this function can perform asynchronous operations (like waiting for locks).
// `data: web::Data<AppState>` extracts the shared application state from the request.
// `-> impl Responder` means the function returns any type that implements the Responder trait (here, HttpResponse).
async fn get_todos(data: web::Data<AppState>) -> impl Responder {
    // Lock the mutex to gain safe access to the todos vector.
    // `.lock().unwrap()`: lock() returns a Result; unwrap() panics if the lock is poisoned (another thread panicked while holding it).
    // For a simple example, we accept this risk.
    let todos = data.todos.lock().unwrap();
    // Return HTTP 200 OK with the todos serialized to JSON.
    // `&*todos` dereferences the MutexGuard to get a reference to the inner Vec, then passes it to .json().
    HttpResponse::Ok().json(&*todos)
}

// Handler for POST /todos – creates a new todo.
// `req: web::Json<CreateTodoRequest>` automatically parses the request body as JSON into our CreateTodoRequest struct.
async fn create_todo(
    req: web::Json<CreateTodoRequest>,
    data: web::Data<AppState>,
) -> impl Responder {
    // Lock the next_id mutex to read and then increment it.
    let mut next_id = data.next_id.lock().unwrap();
    // Build a new Todo using the current next_id and the data from the request.
    let new_todo = Todo {
        id: *next_id,                // Dereference the MutexGuard to get the usize
        is_done: req.is_done,
        // req.name is a &String due to automatic deref; .clone() creates a new owned String.
        name: req.name.clone(),
        description: req.description.clone(),
    };
    // Increment the next_id for future todos.
    *next_id += 1;
    // Lock the todos vector and push the new todo.
    let mut todos = data.todos.lock().unwrap();
    todos.push(new_todo.clone());     // Clone because we also need to return it in the response
    // Return HTTP 201 Created with the newly created todo as JSON.
    HttpResponse::Created().json(new_todo)
}

// Handler for GET /todos/{id} – returns a single todo by its ID.
// `path: web::Path<usize>` extracts the `id` from the URL path and parses it as usize.
async fn get_todo_by_id(path: web::Path<usize>, data: web::Data<AppState>) -> impl Responder {
    // Get the actual id value from the Path wrapper.
    let id = path.into_inner();
    // Lock the todos vector.
    let todos = data.todos.lock().unwrap();
    // Try to find the todo with the given id.
    // `.iter()` creates an iterator over references to the todos.
    // `.find(|t| t.id == id)` returns an Option<&Todo> – Some if found, None if not.
    if let Some(todo) = todos.iter().find(|t| t.id == id) {
        // If found, return HTTP 200 with the todo.
        HttpResponse::Ok().json(todo)
    } else {
        // If not found, return HTTP 404 with a JSON error message.
        // json!({ "error": "Todo not found" }) creates a serde_json Value.
        HttpResponse::NotFound().json(json!({ "error": "Todo not found" }))
    }
}

// Handler for PUT /todos/{id} – updates an existing todo.
async fn update_todo(
    path: web::Path<usize>,
    req: web::Json<UpdateTodoRequest>,  // Same structure as create request
    data: web::Data<AppState>,
) -> impl Responder {
    let id = path.into_inner();
    // Lock the todos vector for modification (need mutable access).
    let mut todos = data.todos.lock().unwrap();
    // Find the todo with mutable reference.
    // `.iter_mut()` gives mutable references to each element.
    // `.find(|t| t.id == id)` returns Option<&mut Todo>.
    if let Some(todo) = todos.iter_mut().find(|t| t.id == id) {
        // Update the fields with the new data.
        todo.is_done = req.is_done;
        todo.name = req.name.clone();
        todo.description = req.description.clone();
        // Return HTTP 200 with the updated todo.
        HttpResponse::Ok().json(todo)
    } else {
        // Todo not found.
        HttpResponse::NotFound().json(json!({ "error": "Todo not found" }))
    }
}

// Handler for DELETE /todos/{id} – deletes a todo.
async fn delete_todo(path: web::Path<usize>, data: web::Data<AppState>) -> impl Responder {
    let id = path.into_inner();
    // Lock the todos vector for modification.
    let mut todos = data.todos.lock().unwrap();
    // Find the position (index) of the todo with the given id.
    // `.iter()` gives an iterator of references, `.position()` returns Some(index) if found.
    if let Some(pos) = todos.iter().position(|t| t.id == id) {
        // Remove the element at that position.
        todos.remove(pos);
        // Return HTTP 204 No Content (successful deletion, no body).
        HttpResponse::NoContent().finish()
    } else {
        // Todo not found.
        HttpResponse::NotFound().json(json!({ "error": "Todo not found" }))
    }
}

// --- Main function: entry point of the application ---
// The `#[actix_web::main]` macro sets up the Actix runtime (similar to #[tokio::main]).
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize the logger from environment variables (default to "info" level).
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));
    // Create the application state wrapped in `web::Data`, which makes it thread-safe and cloneable.
    let app_state = web::Data::new(init_state());
    println!("Starting server at http://localhost:8080");

    // Build and run the HTTP server.
    HttpServer::new(move || {   // `move` captures app_state by value into the closure
        App::new()
            // Attach the shared application state to the app.
            // `.app_data()` expects a `Data` instance; we pass a cloned reference (cheap clone).
            .app_data(app_state.clone())
            // Add the logger middleware to log each request.
            .wrap(Logger::new("%r %s %Dms"))  // Format: method path status_code time_ms
            // Define routes and attach their handlers.
            .route("/todos", web::get().to(get_todos))      // GET /
            .route("/todos", web::post().to(create_todo))   // POST /
            .route("/todos/{id}", web::get().to(get_todo_by_id))   // GET /{id}
            .route("/todos/{id}", web::put().to(update_todo))      // PUT /{id}
            .route("/todos/{id}", web::delete().to(delete_todo))   // DELETE /{id}
    })
    // Bind the server to localhost port 8080.
    .bind(("127.0.0.1", 8080))?
    // Start the server and wait for it to finish (which it never does normally).
    .run()
    .await
}