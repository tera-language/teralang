# TeraLang

TeraLang is a lightweight, easy-to-use language for quickly creating mock APIs and mock backends for frontend development. It allows you to define routes, responses, and mock data with minimal setup, helping you simulate real API interactions without the need for a full backend server.

## Features

- **Simple Syntax**: Define mock routes using an intuitive and easy-to-understand syntax.
- **Route Definition**: Create mock API routes with different HTTP methods (GET, POST, etc.).
- **Flexible Mock Data**: Return static or dynamic mock data for each route.
- **Importable Files**: Reuse TeraLang files by importing them into other mock backends.
- **Customizable Responses**: Easily define responses for specific routes, including JSON-like data.

## Installation

You can install TeraLang via npm:

```bash
npm install teralang
```

# Usage
## Basic Example

Create a **.tera** file to define routes for your mock API. For example, create a file called **mock-api.tera**:

```
route "/ping" GET: {
  status: 200 // optional
  json: {
    message: "pong"
  }
}
```

To start the server, run:

```bash
npx teralang serve mock-api.tera --port 4000
```

# Importing Other Tera Files
To keep your project organized and reduce clutter, you can import other **.tera** files into your main file. This allows you to split different parts of your mock API into separate files and manage them more easily.

For example, if you have a file called app.tera and want to include it in your main file, you can import it like this:

```
import "./app.tera"
```
## Note:

Imported **.tera** files follow the exact same format as your main .tera file—there's no special syntax or setup required.

You can also chain imports inside those files. For example, app.tera can itself import more **.tera** files:
```
import "./users.tera"
import "./posts.tera"
```

## What is coming soon?

- [ ] Database integration support
- [ ] Live reload on file changes
- [ ] VS Code extension for syntax highlighting