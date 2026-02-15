# ASCII-Art-Web

## Description

`ascii-art-web` is a Go HTTP server that exposes a web GUI for ASCII-art generation.
The user submits text and selects a banner, and the server returns the rendered ASCII art.

Supported banners:
- `standard`
- `shadow`
- `thinkertoy`

Implemented endpoints:
- `GET /` renders the main page.
- `POST /ascii-art` receives form data (`text`, `banner`) and returns the result page.

HTTP status codes:
- `200 OK`: successful request
- `400 Bad Request`: invalid method, invalid form data, invalid input/banner
- `404 Not Found`: missing route/template/banner file
- `500 Internal Server Error`: unexpected server/template/render errors

## Authors

- agaliand

## Architecture

High-level architecture:
- `main.go`: server bootstrap (starts HTTP server)
- `web.go`: HTTP routing + handlers + status-code mapping
- `templates/index.html`: main GUI template (form + output rendering)
- `pipeline/web_render.go`: web rendering entrypoint (`RenderASCII`)
- `pipeline/*`: reusable core logic (validate, load banner, tokenize, render)

Request flow:
1. Browser requests `GET /`.
2. Server loads `templates/index.html` and renders the form.
3. User submits `POST /ascii-art` with text and selected banner.
4. Handler validates/parses input and calls `pipeline.RenderASCII`.
5. Pipeline validates input, loads banner, renders ASCII lines.
6. Handler renders template with result or returns proper HTTP error code.

## Usage: how to run

From project root:

```bash
go run .
```

Server runs on:

```text
http://localhost:8080
```

Manual check:
1. Open `http://localhost:8080`
2. Enter text
3. Select banner
4. Click submit

## Tests

Run all tests:

```bash
go test ./...
```

Run a specific test file/package example:

```bash
go test ./tests -run WebRender
```

Note:
- Use `go test ...` (not `go run .test/...`).

## Implementation details: algorithm

ASCII render algorithm used by web mode:
1. Normalize input.
2. Validate input (`ValidateInput`).
3. Load selected banner file (`LoadBanner`).
4. Tokenize input (`Tokenize`).
5. Render token rows to ASCII-art lines (`RenderLines`).
6. Join lines and return output text.

Web handler algorithm:
1. Validate method and parse form.
2. Read `text` and `banner`.
3. Call `RenderASCII`.
4. Map known errors to `400/404/500`.
5. Render HTML with output on success.
