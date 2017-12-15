# per_score_cal

## Getting Started

### Dependencies

### Build and run this project

1. To give privilege to ur .sh file
    ```
  chmod +x setupPostgres.sh
    ```
2. Run .sh file to create role and database
    ```
    ./setupPostgres.sh
    ```
3. Run command to migrate database
    ```
    go run main.go createDB
    ```
4. Run command to start server
    ```
    go run main.go serve
    ```
