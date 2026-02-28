use actix_web::{middleware::Logger, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};
use serde_json::json;
use std::sync::Mutex;

// Full Todo struct (used for responses and storage)
#[derive(Debug, Serialize, Clone)]
struct Todo {
    id: usize,
    is_done: bool,
    name: String,
    description: String,
}

// Struct for creating a new todo (no id field)
#[derive(Debug, Deserialize)]
struct CreateTodoRequest {
    is_done: bool,
    name: String,
    description: String,
}

// Struct for updating a todo (no id field, same fields as create)
type UpdateTodoRequest = CreateTodoRequest; // for simplicity

struct AppState {
    todos: Mutex<Vec<Todo>>,
    next_id: Mutex<usize>,
}

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
        next_id: Mutex::new(3),
    }
}

async fn get_todos(data: web::Data<AppState>) -> impl Responder {
    let todos = data.todos.lock().unwrap();
    HttpResponse::Ok().json(&*todos)
}

async fn create_todo(
    req: web::Json<CreateTodoRequest>,
    data: web::Data<AppState>,
) -> impl Responder {
    let mut next_id = data.next_id.lock().unwrap();
    let new_todo = Todo {
        id: *next_id,
        is_done: req.is_done,
        name: req.name.clone(),
        description: req.description.clone(),
    };
    *next_id += 1;
    let mut todos = data.todos.lock().unwrap();
    todos.push(new_todo.clone());
    HttpResponse::Created().json(new_todo)
}

async fn get_todo_by_id(path: web::Path<usize>, data: web::Data<AppState>) -> impl Responder {
    let id = path.into_inner();
    let todos = data.todos.lock().unwrap();
    if let Some(todo) = todos.iter().find(|t| t.id == id) {
        HttpResponse::Ok().json(todo)
    } else {
        HttpResponse::NotFound().json(json!({ "error": "Todo not found" }))
    }
}

async fn update_todo(
    path: web::Path<usize>,
    req: web::Json<UpdateTodoRequest>,
    data: web::Data<AppState>,
) -> impl Responder {
    let id = path.into_inner();
    let mut todos = data.todos.lock().unwrap();
    if let Some(todo) = todos.iter_mut().find(|t| t.id == id) {
        todo.is_done = req.is_done;
        todo.name = req.name.clone();
        todo.description = req.description.clone();
        HttpResponse::Ok().json(todo)
    } else {
        HttpResponse::NotFound().json(json!({ "error": "Todo not found" }))
    }
}

async fn delete_todo(path: web::Path<usize>, data: web::Data<AppState>) -> impl Responder {
    let id = path.into_inner();
    let mut todos = data.todos.lock().unwrap();
    if let Some(pos) = todos.iter().position(|t| t.id == id) {
        todos.remove(pos);
        HttpResponse::NoContent().finish()
    } else {
        HttpResponse::NotFound().json(json!({ "error": "Todo not found" }))
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init_from_env(env_logger::Env::new().default_filter_or("info"));
    let app_state = web::Data::new(init_state());
    println!("Starting server at http://localhost:8080");

    HttpServer::new(move || {
        App::new()
            .app_data(app_state.clone())
            .wrap(Logger::new("%r %s %Dms"))
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