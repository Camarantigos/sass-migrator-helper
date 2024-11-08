
# Sass Migrator Tool

The `sass-migrator-helper` is a Go-based CLI tool that helps you migrate SCSS files by replacing `@import` paths with relative paths. The tool is designed to convert alias-based imports (e.g., `@styles`) to paths relative to a specified source directory, such as `src`, making it easy to ensure correct imports across all files. After updating the imports, the tool runs `sass-migrator` on the specified entry file to modernize SCSS syntax.

## Prerequisites

- Go installed on your machine. You can download it from [golang.org](https://golang.org/dl/).
- `sass-migrator` installed globally via npm:

  ```bash
  npm install -g sass-migrator
  ```

## Installation

### 1. Clone the Repository

Clone this repository to your local machine:

```bash
git clone <repository-url>
cd <repository-folder>
```

### 2. Build the Tool

To create an executable binary of `sass-migrator-helper`, run the following command in the project’s root directory:

```bash
go build -o sass-migrator-helper
```

This will create a binary named `sass-migrator-helper` in the current directory.

### 3. Move the Binary to a Directory in Your PATH

To use `sass-migrator-helper` from anywhere on your system, move it to a directory included in your system’s `PATH`, such as `/usr/local/bin` on Unix-like systems.

```bash
sudo mv sass-migrator-helper /usr/local/bin/
```

Now you can use `sass-migrator-helper` from any location in your terminal.

## Usage

The tool requires three arguments:

- `-sourceDir`: The root directory containing all SCSS files (e.g., `src`).
- `-entryFile`: The main SCSS entry file for running `sass-migrator`.
- `-alias`: The alias used in `@import` statements that should be replaced with relative paths (e.g., `@styles`).

### Example Command

```bash
sass-migrator-helper -sourceDir ./src -entryFile src/assets/styles/main.scss -alias @styles
```

### Command Breakdown

- **`-sourceDir`**: Specifies the root directory (`src`) from which all relative paths will be calculated.
- **`-entryFile`**: Specifies the main SCSS file where `sass-migrator` will begin the migration. This file, and any dependencies it imports, will be updated.
- **`-alias`**: Specifies the alias (e.g., `@styles`) that will be replaced with a relative path from each file's location.

### How It Works

1. **Locate all `.scss` files in `sourceDir`**: The tool recursively finds all SCSS files within the specified source directory (`src`).
2. **Replace `@import` paths**: The tool replaces any `@import` statements using the specified alias (e.g., `@styles`) with the correct relative path based on each file’s location.
3. **Run `sass-migrator`**: After updating the imports, the tool runs `sass-migrator` on the specified entry file with `--migrate-deps` to modernize the SCSS syntax.

### Output

The tool will print the files it updates and any errors encountered during processing. Once complete, the SCSS files should be updated with correct relative imports and modernized by `sass-migrator`.

### Troubleshooting

- **Path Errors**: If `sass-migrator` cannot find certain files, double-check that the paths in `@import` statements match the actual file structure relative to `sourceDir`.
- **Permissions**: Ensure you have write permissions for the files in `sourceDir` and execution permissions for `sass-migrator-helper` in `/usr/local/bin`.

## Example

Here’s a complete example to show how the tool processes files:

1. Suppose you have a directory structure like this:

   ```
   src/
   ├── assets/
   │   └── styles/
   │       ├── abstracts/
   │       │   └── _variables.scss
   │       ├── common/
   │       │   └── _button.scss
   │       └── main.scss
   └── views/
       └── component/
           └── component.scss
   ```

2. If `_button.scss` contains `@import '@styles/abstracts/variables';`, running the tool will replace this with a relative path like `@import '../../assets/styles/abstracts/variables';`.

3. Finally, the tool will run `sass-migrator` on `main.scss`, migrating the syntax as needed.

## License

This project is licensed under the MIT License.
