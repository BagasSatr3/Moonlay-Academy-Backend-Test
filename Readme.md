````markdown
# Moonlay Todo List API

This API provides endpoints to manage todo lists and sublists. It is built using the Echo framework and GORM for database interactions.

## Getting Started

Follow these instructions to set up and run the Moonlay Todo List API on your local machine.

### Prerequisites

- Go (version 1.18 or later)
- PostgreSQL (version v13 or later)
- Postman (optional, for testing the API)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/moonlay-todo-list-api.git
   ```
````

2. Change into the project directory:

   ```bash
   cd moonlay-todo-list-api
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Set up the PostgreSQL database:

   - Create a database named `moonlay-todo-list`.
   - Update the database configuration in `config/database.go` if needed.

5. Run database migrations:

   ```bash
   go run main.go migrate
   ```

### Usage

1. Start the API server:

   ```bash
   go run main.go
   ```

2. Access the API at [http://localhost:8000](http://localhost:8000).

### API Endpoints

The API provides the following endpoints:

- **Lists:**

  - `GET /lists`: Get lists with optional filtering.
  - `GET /lists/:id`: Get a list by ID.
  - `POST /lists`: Create a new list.
  - `PUT /lists/:id`: Update a list.
  - `DELETE /lists/:id`: Delete a list.

- **Sublists:**
  - `GET /lists/:listID/sublists`: Get sublists for a specific list.
  - `GET /sublists/:id`: Get a sublist by ID.
  - `POST /lists/:listID/sublists`: Create a new sublist for a list.
  - `PUT /sublists/:id`: Update a sublist.
  - `DELETE /sublists/:id`: Delete a sublist.

### File Uploads

- File uploads are supported for creating lists and sublists.
- Supported file types: `.txt` and `.pdf`.
- Maximum file size: 10 MB.

### Running Tests

1. Run all unit tests with coverage:

   ```bash
   go test -cover ./...
   ```

2. Adjust the coverage threshold as needed.

### API Specification

The API specification is available in the [Postman Collection](link-to-postman-collection)
