# Drone Delivery Management Backend

Create a Drone delivery management backend in any programming language.

The API should be authenticated via a JWT, which is handed out by an endpoint that takes a name (for security purposes, pretend that this API will be behind an allow list) and a type of user (admin, enduser, or drone). All other endpoints will use the above JWT for validation. The JWT should be a simple self-signed token and held as a bearer.

The service for the API can be in any language.
The transport for the API can be REST, gRPC, Connect, NATS service or Thrift (multiple is alsofine).

## Acceptance Criteria (ACs)

### Drones should be able to

- [x] Reserve a job.
- [x] Grab an order from a location (origin or broken drone).
- [x] Mark an order they have gotten as delivered or failed.
- [x] Mark themselves as broken (and in need of an order handoff).
- [x] Update their location (use latitude/longitude), and get a status update as a heartbeat.
- [x] Get details on the order they are currently assigned.

### Endusers should be able to

- [x] Submit orders for jobs, with an origin and destination.
- [x] Withdraw orders that have not yet been picked up.
- [x] Get details on orders they have submitted, including its current progress, location and ETA.

### Admins should be able to

- [x] Get multiple orders in ulk, even if they did not submit them.
- [x] Change the origin or destination for an order.
- [x] Get a list of drones.
- [x] Mark drones as broken or fixed.

### Additional Rules

Any time a drone is broken it will stop and put up a job for its goods to be picked up by a different drone (even if it gets marked as fixed).

### Submission Instructions

- Please submit the code as a GitHub, GitLab, or other code hosting site link.
- Write some tests for the code.
- Polish is a big part of the evaluation.
