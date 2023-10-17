# doodocs-assignment



README: Archive Management System
This repository contains a simple Archive Management System written in Go. The system provides basic functionality for managing and interacting with archives and files. It includes a web server that allows users to perform various operations related to archives and files.

## Getting Started
To run this application, follow these steps:

Ensure you have Go installed on your system. If not, you can download and install it from golang.org.

Clone the repository to your local machine:

bash
Copy code
git clone <repository_url>
Change to the project directory:

bash
Copy code
cd <project_directory>
Build and run the application:

bash
Copy code
go run main.go
The application will start the web server, and you can access it in your web browser at http://localhost:4000.

## Project Structure
The project is organized into the following components:

main.go: The main entry point of the application, where the web server is configured and started.

internal/handlers: Contains the HTTP request handlers for various routes.

internal/service: Houses the business logic of the Archive Management System.

## Endpoints
The application provides the following endpoints:

Home Page:

URL: http://localhost:4000
Description: The home page of the application.
Archive Information:

URL: http://localhost:4000/api/archive/information
Description: Retrieves information about archives.
Create Archive:

URL: http://localhost:4000/api/archive/files
Description: Handles the creation of archives.
Send File via Email:

URL: http://localhost:4000/api/mail/file
Description: Handles sending a file via email.
C
