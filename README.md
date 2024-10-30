# Consultant Management System

The Consultant Management System is a web application designed to manage consultants, their tasks, and generate reports. It includes features such as user authentication, consultant management, task assignment, and reporting.

## Features

- **User Authentication**: Secure login and registration system.
- **Consultant Management**: Add, edit, and delete consultants.
- **Task Assignment**: Assign tasks to consultants.
- **Reporting**: Generate and view reports on consultant activities.
- **Consistent UI**: Sidebar and banner for consistent navigation and styling across all pages.

## Technologies Used

- **Backend**: Go (Golang)
- **Frontend**: HTML, Tailwind CSS, HTMX
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Environment Management**: `godotenv`

## Getting Started

### Prerequisites

- Go (Golang) installed
- PostgreSQL installed and running
- Git installed

### Installation

1. **Clone the repository**:
    ```sh
    git clone https://github.com/your-username/consultant-management.git
    cd consultant-management
    ```

2. **Set up environment variables**:
    Create a `.env` file in the project root and add your database credentials:
    ```env
    DB_USER=your_db_user
    DB_NAME=your_db_name
    DB_PASSWORD=your_db_password
    DB_SSLMODE=disable
    ```

3. **Install dependencies**:
    ```sh
    go get -u github.com/gorilla/mux
    go get -u github.com/joho/godotenv
    go get -u github.com/dgrijalva/jwt-go
    go get -u golang.org/x/crypto/bcrypt
    ```

4. **Run the application**:
    ```sh
    go run main.go
    ```

### Usage

1. **Access the application**:
    Open your web browser and navigate to `http://localhost:8080`.

2. **Register a new user**:
    Navigate to `http://localhost:8080/register` to create a new account.

3. **Login**:
    Navigate to `http://localhost:8080/login` to log in with your credentials.

4. **Manage Consultants**:
    Navigate to `http://localhost:8080/consultants` to add, edit, and delete consultants.

5. **Assign Tasks**:
    Use the task assignment form on the consultants page to assign tasks to consultants.

6. **Generate Reports**:
    Navigate to `http://localhost:8080/report` to generate and view reports on consultant activities.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
