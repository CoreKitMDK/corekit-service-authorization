# CoreKit Service Authorization

This repository hosts the `corekit-service-authorization` project, a Go-based microservice designed for robust and efficient management of access rights within the CoreKit ecosystem. Leveraging Knative functions, this service provides a set of granular, event-driven endpoints for handling authorization concerns.

**Key Features:**
- **Rights Management:** Supports operations for granting, retrieving, checking, and revoking user and entity rights to various resources.
- **PostgreSQL Integration:** Utilizes a PostgreSQL database for persistent storage of authorization data, ensuring data integrity and scalability.
- **Knative Functions:** Implemented as a collection of serverless functions, enabling auto-scaling and cost-effective deployment on Kubernetes.
- **Modular Design:** Structured with a clear separation of concerns, including a Data Access Layer (DAL) for database interactions and distinct handlers for each authorization operation.
- **Observability:** Integrated with CoreKit's logging and tracing services for enhanced monitoring and debugging capabilities.

This service is a critical component for enforcing access control policies across CoreKit applications, providing a centralized and performant authorization solution.