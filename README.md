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


## API Endpoints

- **Gestion des Produits**
- `GET /products` : Récupérer tous les produits
- `GET /products/{id}` : Récupérer un produit par son identifiant
- `POST /products` : Ajouter un nouveau produit
- `PUT /products/{id}` : Mettre à jour un produit existant
- `DELETE /products/{id}` : Supprimer un produit

- **Gestion des Catégories**
- `GET /categories` : Récupérer toutes les catégories
- `GET /categories/{id}` : Récupérer une catégorie par son identifiant
- `POST /categories` : Ajouter une nouvelle catégorie
- `PUT /categories/{id}` : Mettre à jour une catégorie existante
- `DELETE /categories/{id}` : Supprimer une catégorie

- **Panier d'Achat**
- `GET /cart` : Récupérer les produits dans le panier
- `POST /cart` : Ajouter un produit au panier
- `DELETE /cart/{id}` : Supprimer un produit du panier
- `DELETE /cart` : Vider le panier

- **Authentification**
- `POST /signup` : Créer un nouveau compte utilisateur
- `POST /login` : Se connecter avec un compte utilisateur existant

- **Gestion des Commandes**
- `GET /orders` : Récupérer toutes les commandes
- `GET /orders/{id}` : Récupérer une commande par son identifiant
- `POST /orders` : Passer une nouvelle commande

---







