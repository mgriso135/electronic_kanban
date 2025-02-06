# Electronic Kanban Backend and Frontend

## Project Description

This project is a web-based Electronic Kanban system designed to manage and visualize Kanban workflows. It consists of a Go-based backend API and a React-based frontend user interface. The system allows users to manage accounts, products, statuses, status chains, kanban chains, and kanban cards, providing dashboards for suppliers and customers to track and interact with their respective kanbans.

## Features

*   **Accounts Management:** CRUD interface for managing customer and supplier accounts.
*   **Products Management:** CRUD interface for managing product information.
*   **Statuses Management:** CRUD interface for defining Kanban statuses (e.g., "To Do", "In Progress", "Shipped").
*   **Status Chains:** Define ordered sequences of statuses to represent Kanban workflows.
    *   Link multiple statuses to a status chain, ordered by sequence.
    *   Define customer/supplier ownership for each status in a chain.
*   **Kanban Chains:** Define Kanban supply chains connecting customers, suppliers, and products.
    *   Specify customer, supplier, product, lead time, container type, quantity, and linked status chain.
    *   Automatic creation of initial Kanban cards upon Kanban Chain creation.
*   **Kanban Cards:** Digital representations of Kanban cards.
    *   Inherit properties from their Kanban Chain (lead time, container type, quantity, status chain).
    *   Track current status and status history.
    *   Support status updates with history logging.
*   **Supplier Dashboard:**
    *   Dashboard for suppliers to view and manage their Kanban cards, organized by product.
    *   Actionable "Change Status" buttons for statuses owned by the supplier.
*   **Customer Dashboard:**
    *   Dashboard for customers to view and manage their Kanban cards, organized by product.
    *   Actionable "Change Status" buttons for statuses owned by the customer.
*   **Kanban List Interface:** CRUD interface to manage all kanban cards with product filtering.
    *   List, create, edit (partially), and delete Kanban cards.
    *   Product-based filtering for Kanban cards.
*   **Real-time Kanban Status Updates:** (Feature in progress - aiming for real-time updates on dashboards).
*   **Status History Tracking:** Records every Kanban status change in a history log.

## Technologies Used

*   **Backend:**
    *   Go (Golang)
    *   github.com/gorilla/mux (HTTP router and multiplexer)
    *   github.com/joho/godotenv (Loading environment variables from `.env` files)
    *   github.com/lib/pq (PostgreSQL driver for Go)
    *   github.com/rs/cors (CORS middleware)
*   **Frontend:**
    *   React
    *   Node.js / npm
    *   axios (HTTP client for API requests)
    *   react-router-dom (Routing for React application)

## Setup Instructions

**Prerequisites:**

*   **Go:**  [Download and install Go](https://go.dev/dl/) (version 1.22 or later recommended).
*   **Node.js and npm:** [Download and install Node.js](https://nodejs.org/en/download/) (npm is included with Node.js).
*   **PostgreSQL:** [Download and install PostgreSQL](https://www.postgresql.org/download/).

**Backend Setup:**

1.  **Navigate to the backend directory:**
    ```bash
    cd electronic_kanban_software/backend
    ```

2.  **Create `.env` file:**
    *   Create a `.env` file in the `backend` directory.
    *   Copy the content from the provided `.env.txt.txt` file into `.env` and adjust the database connection string and port if necessary.  Example `.env` content:
        ```
        DB_CONNECTION_STRING=postgres://postgres:hellas@localhost:5432/electronic_kanban_2?sslmode=disable
        PORT=8080
        ```
    *   **Database Configuration:** Ensure your PostgreSQL database is running and accessible with the credentials specified in the `DB_CONNECTION_STRING`. You may need to create the database `electronic_kanban_2` if it doesn't exist. You will also need to set up the database schema (SQL scripts for table creation are not provided in this README, please refer to database schema documentation if available).

3.  **Download Go dependencies:**
    ```bash
    go mod download
    ```

4.  **Run the backend server:**
    ```bash
    go run main.go
    ```
    The backend server should start and listen on the port specified in your `.env` file (default: `8080`). Check the console output for any errors during startup.

**Frontend Setup:**

1.  **Navigate to the frontend directory:**
    ```bash
    cd electronic_kanban_software/frontend
    ```

2.  **Install npm dependencies:**
    ```bash
    npm install
    ```

3.  **Run the frontend application:**
    ```bash
    npm start
    ```
    The React frontend application should open in your browser, typically at `http://localhost:3000`. Check the console output for any errors during startup.

**How to Use**

1.  **Access the Application:** Open your web browser and navigate to `http://localhost:3000` (or the port where your React frontend is running).

2.  **Navigation:** Use the navigation links at the top of the page to access different sections:
    *   **Accounts:** Manage customer and supplier accounts.
    *   **Products:** Manage product information.
    *   **Statuses:** Manage Kanban statuses.
    *   **Status Chains:** Define and manage status chains.
    *   **Kanban Chains:** Define and manage Kanban chains, linking customers, suppliers, and products.
    *   **Kanbans:** List and manage individual Kanban cards (CRUD interface with product filtering).
    *   **Supplier Dashboard:** Access the Supplier Dashboard (initially for Supplier ID 1, you can change the ID in the URL).
    *   **Customer Dashboard:** Access the Customer Dashboard (initially for Customer ID 1, you can change the ID in the URL).

3.  **CRUD Interfaces:** For Accounts, Products, Statuses, Status Chains, and Kanban Chains, use the links on the list pages to:
    *   **Create New:**  Navigate to the "new" page to create a new entity.
    *   **Edit:** Click "Edit" link in the list to edit an existing entity.
    *   **Delete:** Click "Delete" button in the list to delete an entity.

4.  **Kanban Dashboard Interaction:**
    *   **Select Supplier/Customer:** On the Supplier and Customer Dashboards, use the dropdown to select a specific supplier or customer to view their Kanban data.
    *   **View Kanban Cards:** Kanban cards are organized by product. Each card displays Kanban information and the current status.
    *   **Change Status (Actionable Statuses):** If a "Change Status" button is visible on a Kanban card (and not disabled), it means it's your role (Supplier or Customer, depending on the dashboard) to take action for that status. Click the button to advance the Kanban to the next status in the chain.

5.  **Kanban List Interface:**
    *   **Product Filter:** Use the "Filter by Product" dropdown to filter the Kanban list to show only Kanbans related to specific products.
    *   **Delete Kanbans:** Use the "Delete" button in the Kanban list to mark Kanbans as inactive (soft delete).
    *   **Create New Kanban:** Use the "Create New Kanban" link to add new Kanban cards.
    *   **Edit Kanban:** Use the "Edit" link to modify editable properties of existing Kanban cards (Lead Time, Container Type, Quantity).

## API Endpoints (Backend)

*   **Accounts:**
    *   `GET /api/accounts`: Get all accounts.
    *   `POST /api/accounts`: Create a new account.
    *   `GET /api/accounts/{id}`: Get account by ID.
    *   `PUT/PATCH /api/accounts/{id}`: Update account by ID.
    *   `DELETE /api/accounts/{id}`: Delete account by ID.
*   **Products:**
    *   `GET /api/products`: Get all products.
    *   `POST /api/products`: Create a new product.
    *   `GET /api/products/{id}`: Get product by ID.
    *   `PUT/PATCH /api/products/{id}`: Update product by ID.
    *   `DELETE /api/products/{id}`: Delete product by ID.
*   **Statuses:**
    *   `GET /api/statuses`: Get all statuses.
    *   `POST /api/statuses`: Create a new status.
    *   `GET /api/statuses/{id}`: Get status by ID.
    *   `PUT/PATCH /api/statuses/{id}`: Update status by ID.
    *   `DELETE /api/statuses/{id}`: Delete status by ID.
*   **Status Chains:**
    *   `GET /api/status-chains`: Get all status chains.
    *   `POST /api/status-chains`: Create a new status chain (including linked statuses).
    *   `GET /api/status-chains/{id}`: Get status chain by ID.
    *   `PUT/PATCH /api/status-chains/{id}`: Update status chain by ID.
    *   `DELETE /api/status-chains/{id}`: Delete status chain by ID.
    *   `GET /api/status-chains/{statusChainId}/statuses`: Get statuses for a specific status chain.
    *   `PUT /api/status-chains/{statusChainId}/statuses`: Update statuses for a status chain (order, customer_supplier).
*   **Kanban Chains:**
    *   `GET /api/kanban-chains`: Get all kanban chains.
    *   `POST /api/kanban-chains`: Create a new kanban chain (including initial kanbans).
    *   `GET /api/kanban-chains/{id}`: Get kanban chain by ID.
    *   `PUT/PATCH /api/kanban-chains/{id}`: Update kanban chain by ID.
    *   `DELETE /api/kanban-chains/{id}`: Delete kanban chain by ID.
*   **Kanbans:**
    *   `GET /api/kanbans`: Get all kanbans (supports optional `product_id` query parameter for filtering).
    *   `POST /api/kanbans`: Create a new kanban.
    *   `GET /api/kanbans/{id}`: Get kanban by ID.
    *   `PUT/PATCH /api/kanbans/{id}`: Update kanban by ID (primarily for status updates).
    *   `DELETE /api/kanbans/{id}`: Delete kanban by ID (soft delete - sets `is_active=false`).
*   **Dashboards:**
    *   `GET /api/dashboards/supplier/{supplierId}`: Get supplier dashboard data for a specific supplier.
    *   `GET /api/dashboards/customer/{customerId}`: Get customer dashboard data for a specific customer.

## Contributing

[Optional: Add contribution guidelines here if you plan to make this project open source or accept contributions.]

## License

[Optional: Specify the project license here, e.g., MIT License, Apache 2.0, etc.]

## Contact
Author: Matteo Griso mgrisoster@gmail.com