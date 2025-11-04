# Windy Webcams Platform Importer 
This project is a Go application designed to fetch asset platform data from the [Windy Webcams API](https://api.windy.com/webcams/docs) and store it in a PostgreSQL database. The     
application is structured following Clean Architecture principles to ensure modularity and maintainability.                                                                            
                                                                                                                                                                                       
## Features                                                                                                                                                                            
                                                                                                                                                                                       
- Fetches platform data from the Windy API.                                                                                                                                            
- Stores the data in a PostgreSQL database using GORM.                                                                                                                                 
- Uses GORM's auto-migration to create the necessary database schema.                                                                                                                  
- Handles configuration via environment variables.                                                                                                                                     
                                                                                                                                                                                       
## Prerequisites                                                                                                                                                                       
                                                                                                                                                                                       
Before you begin, ensure you have the following installed:                                                                                                                             
                                                                                                                                                                                       
- **Go**: Version 1.18 or later.                                                                                                                                                       
- **PostgreSQL**: A running instance of PostgreSQL.                                                                                                                                    
                                                                                                                                                                                       
## Getting Started                                                                                                                                                                     
                                                                                                                                                                                       
Follow these steps to set up and run the project locally.                                                                                                                              
                                                                                                                                                                                       
### 1. Clone the Repository                                                                                                                                                            
                                                                                                                                                                                       
```sh                                                                                                                                                                                  
git clone <your-repository-url>                                                                                                                                                        
cd <repository-directory>                                                                                                                                                              
                                                                                                                                                                                       

2. Set Up the Database                                                                                                                                                                 

Make sure your PostgreSQL server is running. Connect to it and create a new database for this project.                                                                                 

                                                                                                                                                                                       
CREATE DATABASE webcams;                                                                                                                                                               
                                                                                                                                                                                       

3. Configure Environment Variables                                                                                                                                                     

The application uses a .env file for configuration.                                                                                                                                    

 1 Copy the example file:                                                                                                                                                              
                                                                                                                                                                                       
   cp .env.example .env                                                                                                                                                                
                                                                                                                                                                                       
 2 Edit the .env file and provide your Windy API key and PostgreSQL connection details:                                                                                                
                                                                                                                                                                                       
   # Windy API Key                                                                                                                                                                     
   WINDY_API_KEY=your_actual_api_key_here                                                                                                                                              
                                                                                                                                                                                       
   # PostgreSQL connection details                                                                                                                                                     
   DB_HOST=localhost                                                                                                                                                                   
   DB_PORT=5432                                                                                                                                                                        
   DB_USER=your_postgres_user                                                                                                                                                          
   DB_PASSWORD=your_postgres_password                                                                                                                                                  
   DB_NAME=webcams                                                                                                                                                                     
   DB_SSLMODE=disable                                                                                                                                                                  
                                                                                                                                                                                       

4. Install Dependencies                                                                                                                                                                

Install the required Go modules:                                                                                                                                                       

                                                                                                                                                                                       
go mod tidy                                                                                                                                                                            
                                                                                                                                                                                       


Running the Application                                                                                                                                                                

To run the application, execute the main.go file. This will trigger a one-time operation to connect to the database, run migrations, fetch platform data from the Windy API, and store 
it.                                                                                                                                                                                    

                                                                                                                                                                                       
go run cmd/app/main.go                                                                                                                                                                 
                                                                                                                                                                                       

You should see output in your console indicating the progress:                                                                                                                         

                                                                                                                                                                                       
Database migration completed.                                                                                                                                                          
Fetching and storing platforms...                                                                                                                                                      
Successfully saved <N> platforms.                                                                                                                                                      
Process finished successfully.                                                                                                                                                         
                                                                                                                                                                                       


Project Structure                                                                                                                                                                      

The codebase is organized to separate concerns, making it easier to maintain and test.                                                                                                 

 • cmd/app: Contains the main application entrypoint.                                                                                                                                  
 • configs: Handles loading configuration from environment variables.                                                                                                                  
 • internal/: Contains the core application logic.                                                                                                                                     
    • domain: Defines the primary data structures (e.g., Platform).                                                                                                                    
    • platform: Implements the business logic (service) and data access (repository) for platforms.                                                                                    
 • pkg/: Contains shared, reusable packages.                                                                                                                                           
    • database: Manages the PostgreSQL database connection and migrations.                                                                                                             
    • windy: Provides a client for interacting with the Windy Webcams API. 