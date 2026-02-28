use actix_web::{web, App, HttpRequest, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::sync::Mutex;

// Define the Todo struct with serde attributes for JSON (same fields as Go)
#[derive(Debug, Serialize, Deserialize, Clone)]
struct Todo {
    id: usize,
    is_done: bool,
    name: String,
    description: String,
}

// App state: holds todos and next ID, wrapped in Mutex for safe sharing
struct AppState {
    todos: Mutex<Vec<Todo>>,
    next_id: Mutex<usize>,
}

// Initialize with two sample todos (matching Go's initial data)
fn init_state() -> AppState {
    AppState {
        todos: Mutex::new(vec![
            Todo {
                id: 1,
                is_done: false,
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
        next_id: Mutex::new(3), // next ID after the two samples
    }
}

// GET /todos – return all todos
async fn get_todos(data: web::Data<AppState>) -> impl Responder {
    let todos = data.todos.lock().unwrap();
    HttpResponse::Ok().json(&*todos)
}

// POST /todos – add a new todo
async fn create_todo(todo: web::Json<Todo>, data: web::Data<AppState>) -> impl Responder {
    let mut new_todo = todo.into_inner();

    // Assign ID and increment next_id
    let mut next_id = data.next_id.lock().unwrap();
    new_todo.id = *next_id;
    *next_id += 1;

    // Append to todos
    let mut todos = data.todos.lock().unwrap();
    todos.push(new_todo.clone());

    HttpResponse::Created().json(new_todo)
}

// GET /todos/{id} – get one todo by ID
async fn get_todo_by_id(path: web::Path<usize>, data: web::Data<AppState>) -> impl Responder {
    let id = path.into_inner();
    let todos = data.todos.lock().unwrap();

    if let Some(todo) = todos.iter().find(|t| t.id == id) {
        HttpResponse::Ok().json(todo)
    } else {
        HttpResponse::NotFound().json(serde_json::json!({ "error": "Todo not found" }))
    }
}

// PUT /todos/{id} – update an existing todo
async fn update_todo(
    path: web::Path<usize>,
    updated: web::Json<Todo>,
    data: web::Data<AppState>,
) -> impl Responder {
    let id = path.into_inner();
    let mut updated_todo = updated.into_inner();
    let mut todos = data.todos.lock().unwrap();

    if let Some(todo) = todos.iter_mut().find(|t| t.id == id) {
        // Preserve ID, update other fields
        updated_todo.id = id;
        *todo = updated_todo.clone();
        HttpResponse::Ok().json(updated_todo)
    } else {
        HttpResponse::NotFound().json(serde_json::json!({ "error": "Todo not found" }))
    }
}

// DELETE /todos/{id} – remove a todo
async fn delete_todo(path: web::Path<usize>, data: web::Data<AppState>) -> impl Responder {
    let id = path.into_inner();
    let mut todos = data.todos.lock().unwrap();

    if let Some(pos) = todos.iter().position(|t| t.id == id) {
        todos.remove(pos);
        HttpResponse::NoContent().finish()
    } else {
        HttpResponse::NotFound().json(serde_json::json!({ "error": "Todo not found" }))
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    // Initialize state with sample data
    let app_state = web::Data::new(init_state());

    println!("Starting server at http://localhost:8080");

    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone()) // share state across workers
            .route("/todos", web::get().to(get_todos))
            .route("/todos", web::post().to(create_todo))
            .route("/todos/{id}", web::get().to(get_todo_by_id))
            .route("/todos/{id}", web::put().to(update_todo))
            .route("/todos/{id}", web::delete().to(delete_todo))
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
