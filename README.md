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

3. **Stop Services:**

   ```bash
   make stop
   ```

   This command stops the services

4. **Start Services:**

   ```bash
   make start
   ```

   This command starts the services.

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

### 1. Create todo

#### `POST /todos`

**Request body:**

```json
{
  "title": "sample title 10",
  "description": "sample description"
}
```

**Response:**

```json
{
  "status": "success",
  "data": {
    "id": "f07efb3e-8c26-4335-996a-6624984254dd",
    "title": "sample title 10",
    "description": "sample description",
    "created_at": "2025-02-27T22:47:28.144958119Z"
  }
}
```

### 2. Retrieve todo by ID.

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

### 3. List Todos

#### `GET /todos`

**Response:**

```json
{
  "status": "success",
  "data": [
    {
      "id": "b8cd0486-e748-4250-b5e9-44ef9943ed34",
      "title": "sample title 10",
      "description": "sample description",
      "created_at": "2025-02-27T22:27:37.631142Z"
    },
    {
      "id": "07434e19-5a11-41a4-8844-707d311059b3",
      "title": "sample title 10",
      "description": "sample description",
      "created_at": "2025-02-27T22:27:03.878102Z"
    },
    {
      "id": "442d0e4f-b73c-4f2b-aac7-5a3ee3b7016b",
      "title": "sample title 5",
      "description": "sample description",
      "created_at": "2025-02-27T22:25:53.754117Z"
    }
  ]
}
```

### 4. Update todo

#### `PATCH /todos/:id`

**Request body:**

```json
{
  "title": "sample title 10", // optional
  "description": "sample description" // optional
}
```

**Response:**

```json
{
  "status": "success",
  "data": {
    "id": "07434e19-5a11-41a4-8844-707d311059b3",
    "title": "updated sample title 10",
    "description": "updated sample title 10",
    "created_at": "2025-02-27T22:27:03.878102Z",
    "updated_at": "2025-02-27T22:33:44.10342Z"
  }
}
```

### Delete todo by ID

#### `DELETE /todos/:id`

**Request Path Variables:**

- `id` (required)

**Response:**

```json
{
    "status": "success",
    "data": {
        "success": true
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
