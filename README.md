# xpense-trackr

`xpense-trackr` is a full-stack web application for tracking expenses, managing budgets, and analyzing spending habits. It provides users with a comprehensive toolset to manage their finances effectively.

## Features

- **Expense Management:** Add, edit, and delete expenses with ease. Categorize expenses by type, date, and amount.
- **Category Management:** Create, update, and delete categories to organize expenses efficiently.
- **Budget Setting:** Set budget limits for different expense categories and track spending against budget targets.
- **Dashboard Overview:** View summary information, charts, and graphs to monitor spending trends and budget progress.
- **User Authentication:** Secure user authentication and authorization system to protect user data and privacy.
- **Responsive Design:** Mobile-friendly interface ensures seamless user experience across devices.

## Technologies Used

- **Backend:** Go (Golang) with Gorilla Mux for routing, JWT for authentication, and GORM for database interaction.
- **Frontend:** React with Redux for state management, React Router for navigation, and Axios for API requests.
- **Database:** MySQL or PostgreSQL for data storage and retrieval.
- **Deployment:** Docker for containerization and deployment to cloud platforms like AWS, GCP, or Azure.

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 20+
- MySQL or PostgreSQL
- Docker (optional)

### Installation

1. Clone the repository:

   ```
    git clone https://github.com/your-username/xpense-trackr.git
   ```

2. Navigate to the project directory:
   ```
    cd xpense-trackr
   ```

3. Install backend dependencies:
   ```
    go mod tidy
   ```

4. Install frontend dependencies:
   ```
   cd frontend
   npm install
   ```

5. Start the backend server:
   ```
    go run main.go
   ```

6. Start the frontend development server:
   ```
   cd frontend
   npm start
   ```

7. After following these steps, you can open your web browser and navigate to http://localhost:3000 to access the xpense-trackr application.
