# Product Manager API

A RESTful API built with Go for managing products, categories, sales, and user authentication. This project provides a complete product management solution with JWT-based authentication, PostgreSQL database integration, and Swagger documentation.

## ✨ Features

* **User Authentication**: JWT-based authentication with registration, login, logout, and token refresh
* **Product Management**: Full CRUD operations for products with category association
* **Category Management**: Complete category management system
* **Sales Tracking**: Sales recording with product associations
* **Database Integration**: PostgreSQL with GORM ORM for automatic migrations
* **API Documentation**: Interactive Swagger UI documentation
* **Middleware Support**: Authentication middleware and CORS support
* **Secure**: Password hashing and JWT token validation

## 🧰 Project Overview

The `productmanagerapi` is structured to provide a modular and scalable approach to product management through a RESTful API. The project is organized into several key directories:

* **`cmd/`**: Contains the main application entry point and Swagger documentation
* **`config/`**: Database configuration and application settings
* **`controllers/`**: HTTP request handlers for all API endpoints
* **`models/`**: GORM data models (User, Product, Category, Sale, SaleProduct)
* **`responseFormatter/`**: Standardized API response formatting
* **`routes/`**: API endpoint routing configuration
* **`services/`**: Business logic layer for data operations
* **`types/`**: Custom type definitions
* **`utils/`**: Utility functions including authentication middleware

## 🏗️ Technology Stack

* **Language**: Go 1.23.5
* **Database**: PostgreSQL with GORM ORM
* **Authentication**: JWT (JSON Web Tokens)
* **Documentation**: Swagger/OpenAPI
* **HTTP Routing**: Native Go HTTP multiplexer

## 📊 Data Models

The application uses the following data models:

### User

* ID, Username, Password, Email, Role
* Supports user authentication and role-based access

### Category

* ID, Name, Description
* Used to organize products into categories

### Product

* ID, Name, Description, Price, Stock, CategoryID
* Core product information with stock tracking

### Sale

* ID, Products (SaleProduct array), Total
* Represents a complete sale transaction

### SaleProduct

* ID, SaleID, ProductID, Quantity, Total
* Junction table for many-to-many relationship between sales and products

---

## 🚀 Getting Started

### Prerequisites

* Go 1.16 or higher (Go 1.23.5 recommended)
* PostgreSQL database
* Git

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Bakemono-san/productmanagerapi.git
   cd productmanagerapi
   ```

2. **Set up PostgreSQL database:**

   Create a PostgreSQL database named `product_management` and update the connection string in `config/dbconfig.go` if needed:

   ```go
   const dsn = "host=localhost user=bakemono password=bakemono dbname=product_management port=5432 sslmode=disable"
   ```

3. **Download dependencies:**

   ```bash
   go mod download
   ```

4. **Run the application:**

   ```bash
   go run cmd/main.go
   ```

   The server will start on `http://localhost:2002`

### Database Migration

The application automatically creates the required database tables on startup using GORM's AutoMigrate feature for the following models:

* Users
* Categories
* Products
* Sales
* SaleProducts



This command starts the API server.

---

## 📚 API Documentation

The API provides endpoints for managing products, categories, sales, and user authentication. Most endpoints require JWT authentication (except auth endpoints and Swagger docs).

### Authentication Endpoints

* **POST** `/auth/register` - Register a new user
* **POST** `/auth/login` - User login
* **POST** `/refresh-token` - Refresh JWT token
* **POST** `/logout` - User logout

### Product Endpoints

* **GET** `/products` - Retrieve all products
* **GET** `/product?id={id}` - Retrieve a specific product by ID
* **POST** `/create-product` - Create a new product
* **PUT** `/update-product` - Update an existing product
* **DELETE** `/delete-product` - Delete a product

### Category Endpoints

* **GET** `/categories` - Retrieve all categories
* **GET** `/category?id={id}` - Retrieve a specific category by ID
* **POST** `/create-category` - Create a new category
* **PUT** `/update-category` - Update an existing category
* **DELETE** `/delete-category` - Delete a category

### Sales Endpoints

* **GET** `/sales` - Retrieve all sales
* **POST** `/create-sale` - Create a new sale
* **DELETE** `/delete-sale` - Delete a sale

### API Documentation

* **GET** `/swagger/*` - Swagger UI documentation

All endpoints (except auth and swagger) require JWT authentication via the Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

---

## 🛠️ Configuration

### Database Configuration

The application uses PostgreSQL as the database. Configuration is managed in `config/dbconfig.go`:

```go
const dsn = "host=localhost user=bakemono password=bakemono dbname=product_management port=5432 sslmode=disable"
```

Update this connection string according to your PostgreSQL setup:

* **host**: PostgreSQL server host (default: localhost)
* **user**: Database username
* **password**: Database password
* **dbname**: Database name (default: product_management)
* **port**: PostgreSQL port (default: 5432)

### JWT Configuration

JWT secret key is configured in `config/dbconfig.go`:

```go
var SECRET_KEY = "BAKEMONO_SECRET"
```

**Important**: Change this to a secure secret key in production.

### Server Configuration

The API server runs on port 2002 by default. This can be modified in `cmd/main.go`.

---

## 🧪 Testing the API

### Using cURL

First, register a user and login to get a JWT token:

```bash
# Register a new user
curl -X POST http://localhost:2002/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"testpass","role":"admin"}'

# Login to get JWT token
curl -X POST http://localhost:2002/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'
```

Then use the token to access protected endpoints:

```bash
# Get all products (replace YOUR_TOKEN with actual token)
curl -X GET http://localhost:2002/products \
  -H "Authorization: Bearer YOUR_TOKEN"

# Create a category
curl -X POST http://localhost:2002/create-category \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics","description":"Electronic products"}'

# Create a product
curl -X POST http://localhost:2002/create-product \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","description":"Gaming laptop","price":999.99,"stock":10,"categoryId":1}'
```

### Using Swagger UI

Access the interactive API documentation at:

```url
http://localhost:2002/swagger/index.html
```

### Using Postman

Import the API endpoints into Postman and set up Bearer token authentication in the Authorization tab.

## 🔧 Development

### Project Structure

```text
productmanagerapi/
├── cmd/
│   ├── main.go              # Application entry point
│   └── docs/                # Swagger documentation files
├── config/
│   └── dbconfig.go          # Database and app configuration
├── controllers/
│   ├── authController.go    # Authentication endpoints
│   ├── categoryController.go # Category CRUD operations
│   ├── productController.go  # Product CRUD operations
│   └── saleController.go     # Sales management
├── models/
│   └── model.go             # GORM data models
├── routes/
│   └── routes.go            # API route definitions
├── services/
│   ├── authService.go       # Authentication business logic
│   ├── categoryService.go   # Category business logic
│   ├── productService.go    # Product business logic
│   └── salesService.go      # Sales business logic
├── types/
│   └── types.go             # Custom type definitions
├── utils/
│   └── utils.go             # Utility functions and middleware
└── responseFormatter/
    └── responseFormatter.go # API response formatting
```

### Key Dependencies

* `gorm.io/gorm` - ORM for database operations
* `gorm.io/driver/postgres` - PostgreSQL driver for GORM
* `github.com/golang-jwt/jwt/v5` - JWT implementation
* `github.com/swaggo/swag` - Swagger documentation generator

---

## 🤝 Contributing

Contributions are welcome! If you'd like to contribute to the `productmanagerapi` project, please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature/your-feature-name`.
3. Make your changes and commit them: `git commit -m 'Add your feature'`.
4. Push to the branch: `git push origin feature/your-feature-name`.
5. Open a pull request.

Please ensure your code adheres to the existing coding standards and includes appropriate tests.

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

## 📫 Contact

For questions or support, please open an issue in the repository or contact the maintainer directly.
