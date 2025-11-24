---
title: "[0007] Standard API Error Response Format"
description: >
    Define a consistent error response structure for all API endpoints to improve client error handling and debugging
type: docs
weight: 7
category: "api-design"
status: "proposed"
date: 2025-11-24
deciders: []
consulted: []
informed: []
---

<!--
ADR Categories:
- strategic: High-level architectural decisions (frameworks, auth strategies, cross-cutting patterns)
- user-journey: Solutions for specific user journey problems (feature implementation approaches)
- api-design: API endpoint design decisions (pagination, filtering, bulk operations)
-->

## Context and Problem Statement

Currently, there is no standardized error response format across API endpoints in the infrastructure. This inconsistency makes it difficult for clients to handle errors uniformly, complicates debugging, and reduces the overall developer experience. How should we structure error responses to ensure consistency, provide adequate debugging information, and follow industry best practices?

<!-- This is an optional element. Feel free to remove. -->
## Decision Drivers

<!--
For Strategic ADRs, consider: scalability, maintainability, team expertise, ecosystem maturity, vendor lock-in
For User Journey ADRs, consider: user experience, security, implementation complexity, timeline
For API Design ADRs, consider: client usage patterns, performance, backward compatibility, API conventions
-->

* Client needs consistent error parsing across all endpoints
* Debugging requires sufficient context and error details
* Security considerations (avoid leaking sensitive implementation details)
* OpenAPI specification compatibility (services use OpenAPI-first design)
* Machine-readable error codes for programmatic handling
* Human-readable error messages for logging and debugging
* Support for field-level validation errors
* Alignment with REST/HTTP standards and industry practices

## Considered Options

* RFC 7807 Problem Details (application/problem+json)
* Custom structured error format with error codes
* Simple error message string
* Google API error response format (errors array)

## Decision Outcome

Chosen option: "RFC 7807 Problem Details (application/problem+json)", because it is an industry standard (RFC), provides extensibility, is widely supported by tooling and libraries, includes standard fields for error identification and debugging, and integrates well with OpenAPI specifications.

<!-- This is an optional element. Feel free to remove. -->
### Consequences

* Good, because RFC 7807 is a well-established standard with broad ecosystem support
* Good, because it provides a consistent structure with required and optional fields
* Good, because it supports extensibility through custom properties
* Good, because OpenAPI 3.x has native support for Problem Details via schemas
* Good, because it includes both machine-readable (type) and human-readable (title, detail) information
* Neutral, because it requires clients to handle a specific content type (application/problem+json)
* Bad, because it adds slight complexity compared to simple error strings
* Bad, because developers need to learn the RFC 7807 structure if unfamiliar

<!-- This is an optional element. Feel free to remove. -->
### Confirmation

Implementation compliance will be confirmed through:
1. OpenAPI schema definitions requiring RFC 7807 structure for error responses
2. Code review process ensuring endpoints use the standard error format
3. Integration tests validating error response structure and content type
4. Go service framework helpers/middleware to generate RFC 7807 responses

<!-- This is an optional element. Feel free to remove. -->
## Pros and Cons of the Options

### RFC 7807 Problem Details (application/problem+json)

RFC 7807 defines a standard JSON structure for HTTP API error responses with fields: `type` (URI reference), `title`, `status` (HTTP status code), `detail`, and `instance` (URI reference to specific occurrence).

* Good, because it is an IETF standard (RFC 7807) with wide industry adoption
* Good, because it provides both machine-readable (`type`, `status`) and human-readable (`title`, `detail`) information
* Good, because it supports extension fields for custom data (e.g., validation errors)
* Good, because OpenAPI 3.x natively supports RFC 7807 schemas
* Good, because many HTTP libraries and frameworks have built-in support
* Good, because the `instance` field helps trace specific error occurrences
* Neutral, because requires specific `Content-Type: application/problem+json` header
* Bad, because it requires more implementation effort than simple error strings
* Bad, because teams unfamiliar with RFC 7807 face a learning curve

### Custom structured error format with error codes

Define a project-specific JSON error structure with custom fields like `error_code`, `message`, `details`, etc.

* Good, because it can be tailored exactly to project needs
* Good, because it's flexible and can evolve with requirements
* Good, because it avoids external standard dependencies
* Neutral, because requires defining and documenting the structure
* Bad, because it lacks ecosystem tooling and library support
* Bad, because clients must learn a custom format instead of a standard
* Bad, because it reinvents a solution that already exists as an RFC
* Bad, because it's harder to integrate with standard OpenAPI tooling

### Simple error message string

Return errors as plain text strings or simple JSON objects with a single `message` field.

* Good, because it is extremely simple to implement
* Good, because it requires minimal client parsing logic
* Good, because it has no learning curve
* Neutral, because it works for very simple APIs
* Bad, because it lacks structure for programmatic error handling
* Bad, because it cannot distinguish between error types without parsing messages
* Bad, because it provides no standard fields for HTTP status, error codes, or metadata
* Bad, because it doesn't scale well for complex validation errors
* Bad, because it complicates debugging without structured fields

### Google API error response format (errors array)

Use Google's error response format with a top-level `error` object containing `code`, `message`, and an `errors` array with detailed error information.

* Good, because it is used by a major tech company (Google)
* Good, because it supports multiple errors in a single response
* Good, because it includes structured error details
* Good, because it is well-documented in Google's API design guide
* Neutral, because it is familiar to developers who use Google APIs
* Bad, because it is not a formal standard (not an RFC)
* Bad, because it has less ecosystem support than RFC 7807
* Bad, because it duplicates HTTP status code in the JSON payload
* Bad, because it requires custom schema definitions rather than using standard patterns

<!-- This is an optional element. Feel free to remove. -->
## More Information

### RFC 7807 Example Response

```json
{
  "type": "https://api.example.com/errors/validation-error",
  "title": "Validation Error",
  "status": 400,
  "detail": "The request body failed validation",
  "instance": "/api/v1/users/create",
  "invalid_fields": [
    {
      "field": "email",
      "reason": "must be a valid email address"
    }
  ]
}
```

### References

* RFC 7807 - Problem Details for HTTP APIs: https://www.rfc-editor.org/rfc/rfc7807.html
* OpenAPI 3.x support for Problem Details
* Go implementation libraries: `github.com/moogar0880/problems` or custom implementation
* Related ADR: Resource Identifier Standard (0006) for `type` and `instance` URI patterns
