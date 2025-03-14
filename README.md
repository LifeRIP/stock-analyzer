# Stock-Analyzer

Stock-Analyzer es una aplicación diseñada para analizar datos de inventarios y proporcionar información útil para la gestión de stocks.

## Requisitos

- [Docker](https://www.docker.com/) instalado en tu sistema.
- Archivo `.env` en el directorio del backend.

## Configuración

1. **Configurar el entorno**:
   - En el directorio del backend, asegúrate de incluir un archivo `.env` basado en el `.env.example` proporcionado.
   - Puedes copiarlo y renombrarlo con el siguiente comando:
     ```sh
     cp backend/.env.example backend/.env
     ```
   - Luego, edita el archivo `.env` con las configuraciones necesarias.

2. **Iniciar el proyecto**:
   - Desde la raíz del proyecto, ejecuta el siguiente comando:
     ```sh
     docker-compose up
     ```
   - Esto levantará los contenedores necesarios para ejecutar la aplicación.

## Uso

Una vez que los contenedores estén en ejecución, puedes acceder a la aplicación a través de las siguientes URLs:

- **CockroachDB Web UI**: [http://localhost:8080/](http://localhost:8080/)
- **API**: [http://localhost:8081/](http://localhost:8081/)
- **App Web UI**: [http://localhost:8082/](http://localhost:8082/)

