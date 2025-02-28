# Todo Microservice

This is a microservice-based Todo app that uses gRPC for communication between services.

## Getting Started

### Prerequisites

- Docker
- Go
- Protocol Buffers compiler (protoc) with Go plugins (protoc-gen-go and protoc-gen-go-grpc)

## How to start project

Use the provided Makefile commands:

1. **Clone the Repository:**

   ```bash
   git clone git@github.com:victor-nach/todo.git
   cd todo
   ```

2. **Build and Run Services:**

   ```bash
   make build-run
   ```

   This command builds the Docker images and starts the services.

---

### Models
- Todo model
    ```json
        {
            "id": "string",
            "title": "string",
            "description": "string", // optional
            "created_at": "timestamp",
            "updated_at": "timestamp"
        }
    ```

---

## Endpoints

### Retrieve todo by ID.

#### `GET /todos/:id`

**Request Path Variables:**

- `id` (required)

**Response:**

```json
{
    "status": "success",
    "data": {
        "id": "b8cd0486-e748-4250-b5e9-44ef9943ed34",
        "title": "sample title 10",
        "description": "sample description",
        "created_at": "2025-02-27T22:27:37.631142Z"
    }
}
```

### **Postman collection**

A Postman collection is provided in the root of the project as .postman_collection.json. You can use this collection to quickly explore and test the API endpoints.

To use the Postman collection:

1. Open Postman.
2. Click on "Import" and select the .postman_collection.json file from the project root.
3. Explore the endpoints included in the collection and run sample requests.

---
