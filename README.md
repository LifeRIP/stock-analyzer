# Stock-Analyzer

Stock-Analyzer is an application designed to analyze stock data retrieved from an external API and provide useful insights for investing.

## Requirements

- [Docker](https://www.docker.com/) installed on your system.
- `.env` file in the backend directory.

## Setup

1. **Configure the environment**:
   - In the backend directory, make sure to include a `.env` file based on the provided `.env.example`.
   - You can copy and rename it using the following command:
     ```sh
     cp backend/.env.example backend/.env
     ```
   - Then, edit the `.env` file with the necessary configurations.

2. **Start the project**:
   - From the project root, run the following command:
     ```sh
     docker-compose up
     ```
   - This will start the necessary containers to run the application.

## Usage

Once the containers are up and running, you can access the application through the following URLs:

- **CockroachDB Web UI**: [http://localhost:8080/](http://localhost:8080/)
- **API**: [http://localhost:8081/](http://localhost:8081/)
- **App Web UI**: [http://localhost:8082/](http://localhost:8082/)

