# Ecommerce PLATFORM - Gonic FULLBUILD with MongoDB

Ecommerce PLATFORM is an e-commerce platform built with Gonic FULLBUILD and MongoDB. This platform provides a solid foundation for creating a high-performance and scalable e-commerce website.

## Features

- Product Management: Adding, editing, and deleting products.
- Category Management: Creating, editing, and deleting product categories.
- Shopping Cart: Managing products added to the shopping cart and the order process.
- Authentication: User registration and login system.
- Order Management: Tracking and managing customer orders.

## Requirements

- Golang (Gonic FULLBUILD)
- MongoDB database
- Gonic Dependencies (github.com/gin-gonic/gin)
- MongoDB Driver for Go (go.mongodb.org/mongo-driver)

## Getting Started

To run the Ecommerce PLATFORM project with MongoDB, follow these steps:

1. **Clone the Repository**
```
git clone https://github.com/Nutstothemoo/Ecommerce.git
```



2. **Install Dependencies**

go mod tidy 


3. **Configure MongoDB**

Configure the MongoDB connection in the configuration file.

4. **Run the Application**

go run *.go

The application will now be accessible at `http://localhost:8080`.

## API Endpoints

## User Routes

- `POST /users/signup` : Create a new user account.
- `POST /users/login` : Log in with an existing user account.

## Admin Routes

- `POST /admin/addproduct` : Add a new product to the platform.

## Product Routes

- `GET /users/productview` : Retrieve all products.
- `GET /users/search` : Search for products using a query.


## Contribution

Contributions are welcome! If you encounter any issues or want to add new features, feel free to open an issue or submit a pull request.









