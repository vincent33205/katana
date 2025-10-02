# Parser Refactor Suggestions

This document outlines opportunities to split `pkg/engine/parser/parser.go` into smaller, testable units and to isolate functionality for future maintenance.

## Current Pain Points

- **Single giant registry.** `NewResponseParser` wires up every parser, mixing header, body, and custom logic in one slice, which makes it hard to see ownership boundaries or extend specific areas independently.
- **Highly repetitive attribute scrapers.** Most body parsers share the same pattern: select elements, read one or more attributes, and enqueue navigation requests. The copy‑paste structure obscures the minimal differences and increases the risk of missing new attributes.
- **Mixed responsibilities.** The file blends simple attribute extraction with heavier logic (for example, forms and HTMX requests) and even regex driven scraping. These should evolve separately with tailored tests and abstractions.

## Splitting the File

1. **Create concern-focused files (or subpackages).**
   - `header.go`: header parsers and `NewHeaderParsers()`.
   - `body/attributes.go`: simple attribute-driven parsers that can be data-driven.
   - `body/forms.go`: complex form submission builder.
   - `body/htmx.go` and `body/resources.go`: HTMX attributes and media tags.
   - `regex.go`: the custom regex parser.
   Tie them together with a thin `registry.go` that merges the slices.

2. **Adopt data-driven attribute parsers.**
   - Describe parsers with a struct (`selector`, `attributes`, `tag`, `attributeLabel`).
   - Loop through the registry to generate navigation requests, reducing ~30 functions to declarative entries.
   - Allow per-entry hooks for special cases (e.g., `img[src]` skipping `data:` URLs).

3. **Isolate complex flows.**
   - Move form handling (multipart, field guessing) into `formbuilder` helpers that return a fully-populated `navigation.Request`.
   - Extract HTMX handling into its own module with shared helpers for methods and attribute naming.

4. **Introduce parser interfaces.**
   - Define an interface like `type ResponseParser interface { Parse(*navigation.Response) ([]*navigation.Request, error) }`.
   - Wrap existing functions with adapters so larger features (forms/htmx) can expose richer errors and be unit tested outside the registry loop.

## Testing Strategy

- **Table-driven tests for declarative parsers.** Provide HTML snippets and expected URLs, run them through the generic attribute parser, and assert the produced navigation requests.
- **Focused unit tests for complex helpers.**
  - Forms: build tests covering GET/POST, enctype changes, multipart boundaries, and form-fill suggestions.
  - HTMX: validate HTTP method translation and absolute URL resolution.
- **Golden HTML fixtures.** Keep small fixture files under `testdata/` for multi-element cases; parse them with `goquery` to simulate realistic responses.
- **Integration smoke tests.** Once parsers are modular, reuse them in higher-level tests that feed a `navigation.Response` with both headers and body to ensure combined behavior.

## Isolation Improvements

- **Abstract navigation creation.** Inject a small interface for `NewNavigationRequestURLFromResponse` so tests can verify the inputs without depending on the concrete implementation.
- **Wrap third-party dependencies.** Provide adapters around `goquery.Selection` for easier mocking/faking in tests.
- **Surface errors.** Return errors from complex parsers instead of logging directly; let the caller decide whether to log or skip.

These steps should make `parser.go` dramatically smaller, clarify responsibilities, and unlock targeted tests for each parsing concern.
