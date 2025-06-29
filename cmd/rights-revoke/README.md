# Knative Function: rights-revoke

This Knative function, `rights-revoke`, is a vital component of the `corekit-service-authorization` microservice, responsible for deactivating or revoking specific access rights previously granted to an entity.

**Functionality:**
- **Input:** It accepts an HTTP POST request with a JSON payload structured as a `RevokeRightsRequest`. This request includes the `Entity` (a UUID) from which rights are to be revoked, and a list of `Rights` (UUIDs) representing the specific rights to be deactivated.
- **Processing:** Upon receiving a valid request, the function interacts with the underlying PostgreSQL database through the Authorization Data Access Layer (DAL). It performs an `UPDATE` operation on the `rights` table, setting the `active` flag to `false` for the specified rights belonging to the given entity. This effectively deactivates the rights without physically deleting the records, allowing for auditing and potential re-activation if needed.
- **Output:** It returns an HTTP 200 OK response with a JSON payload representing a `RevokeRightsResponse`. This response includes the `Entity` for which rights were processed, a `Valid` boolean flag indicating the success of the operation, and an `Error` string if any issues occurred during processing.

This function is essential for maintaining strict access control and ensuring that entities only retain necessary permissions within the CoreKit ecosystem.