# Knative Function: rights-has

This Knative function, `rights-has`, is a crucial part of the `corekit-service-authorization` microservice, designed to efficiently determine if a given entity possesses specific rights to a set of resources.

**Functionality:**
- **Input:** It accepts an HTTP POST request with a JSON payload structured as a `HasRightsRequest`. This request includes the `Entity` (a UUID) for which rights are being checked, and a list of `Resources` (UUID strings) against which the entity's permissions are to be evaluated.
- **Processing:** The function queries the underlying PostgreSQL database through the Authorization Data Access Layer (DAL). It retrieves all active rights for the specified `Entity` that match any of the provided `Resources`. The DAL's `HasRights` method then processes these retrieved rights to determine the effective permissions for each requested resource.
- **Output:** It returns an HTTP 200 OK response with a JSON payload representing a `HasRightsResponse`. This response contains the `Entity` that was queried, a map of `Rights` (where keys are resource identifiers and values are `Right` objects indicating the entity's permissions for that resource), a `Valid` boolean flag, and an `Error` string if the operation encountered issues.

This function is essential for real-time authorization checks, allowing other services to quickly verify access privileges.