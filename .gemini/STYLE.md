# Foreign7 Coding Standards & Style Guide

This document outlines the coding standards and conventions for the Foreign7 project. Adherence to these rules ensures consistency and maintainability across the codebase.

## 🚀 Core Principles

### 1. Language & Naming
- **English Only:** All identifiers (variables, functions, types, etc.) and comments MUST be in English.
- **Naming Conventions:**
    - Use `PascalCase` for exported symbols.
    - Use `camelCase` for internal symbols.
    - Use short, descriptive names. Avoid generic names like `data` or `info` when more context can be provided.
    - Acronyms should be consistent (e.g., `JSONData` or `jsonData`, not `JsonData`).

### 2. Modern Go Practices
- **Use `any`:** Always use the `any` keyword instead of `interface{}` for empty interfaces.
- **Explicit Error Handling:** Never ignore errors. Handle them explicitly using the `if err != nil` pattern.
- **Internal Packages:** Keep implementation details in `internal/` to prevent unwanted external usage.

### 3. Documentation
- **Comments:** Write meaningful comments in English.
- **Exported Symbols:** Every exported type, function, and variable should have a documentation comment explaining its purpose.

### 4. Formatting
- **gofmt:** All code must be formatted with `gofmt` or `goimports`.
- **Grouped Imports:** Group imports into three sections: standard library, third-party packages, and internal project packages, separated by a newline.
