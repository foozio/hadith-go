# API Docs

This folder serves API documentation at `/docs` via the REST server.

By default, the page attempts to load a local Redoc bundle from `/docs/redoc.standalone.js` and falls back to the public CDN if not available.

For fully self-contained (offline) docs, download the Redoc standalone bundle and place it here:

- File: `web/docs/redoc.standalone.js`
- Source: https://github.com/Redocly/redoc/releases (look for `redoc.standalone.js`)

Then open: http://localhost:8080/docs

The OpenAPI spec is served from `/openapi.yaml` by the server.

