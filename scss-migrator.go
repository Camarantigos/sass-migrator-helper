package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// printHelp displays the usage information for the script.
func printHelp() {
	fmt.Println(`Usage: sass-migrator-helper -sourceDir <source-directory> -entryFile <entry-file> -alias <alias>
	
This script modifies @import paths with a specified alias in SCSS files to use relative paths,
then runs the sass-migrator on each updated SCSS file and the main entry file.

Flags:
  -sourceDir     The root directory containing all SCSS files (e.g., src). Required.
  -entryFile     The main SCSS entry file for sass-migrator. Required.
  -alias         The alias (e.g., @styles) to replace with relative paths. Required.`)
}

// findSCSSFiles recursively searches for .scss files in the specified directory.
func findSCSSFiles(root string) ([]string, error) {
	var scssFiles []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".scss") {
			scssFiles = append(scssFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return scssFiles, nil
}

// replaceAliasWithRelativePath replaces import paths using an alias with relative paths.
func replaceAliasWithRelativePath(filePath, alias, srcRoot string) (bool, error) {
	// Open the file for reading
	contentFile, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer contentFile.Close()

	// Read content from the file
	content, err := io.ReadAll(contentFile)
	if err != nil {
		return false, err
	}

	// Convert the content to a string for manipulation
	fileContent := string(content)

	// Define the regex to find @import statements that use the alias
	re := regexp.MustCompile(`@import ['"]` + regexp.QuoteMeta(alias) + `([^'"]+)['"];`)
	modifiedContent := re.ReplaceAllStringFunc(fileContent, func(match string) string {
		importPath := re.FindStringSubmatch(match)[1]

		// Create a full path by appending importPath to srcRoot
		fullImportPath := filepath.Join(srcRoot, "assets/styles", importPath)

		// Generate a relative path from the src directory to the target import path
		relPath, err := filepath.Rel(filepath.Dir(filePath), fullImportPath)
		if err != nil {
			fmt.Printf("Error creating relative path: %v\n", err)
			return match // Keep the original match in case of an error
		}

		// Return the modified @import statement with the correctly calculated relative path
		return fmt.Sprintf(`@import '%s';`, relPath)
	})

	// If the content was modified, write it back to the file
	if modifiedContent != fileContent {
		// Open the file for writing
		err = os.WriteFile(filePath, []byte(modifiedContent), 0644)
		if err != nil {
			return false, err
		}
		fmt.Printf("Updated imports in: %s\n", filePath)
		return true, nil
	}
	return false, nil
}

// runSassMigrator runs the sass-migrator command on a specific file.
func runSassMigrator(file string) error {
	cmd := exec.Command("sass-migrator", "--migrate-deps", "module", file)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// Define flags for the source directory, entry file, and alias
	sourceDir := flag.String("sourceDir", "", "The root directory containing SCSS files (e.g., src).")
	entryFile := flag.String("entryFile", "", "The main SCSS entry file for migration.")
	alias := flag.String("alias", "", "The alias in import paths to replace with relative paths (e.g., @styles).")
	showHelp := flag.Bool("help", false, "Show help message.")

	flag.Parse()

	// Show help if needed
	if *showHelp || *sourceDir == "" || *entryFile == "" || *alias == "" {
		printHelp()
		return
	}

	// Find .scss files in the source directory
	scssFiles, err := findSCSSFiles(*sourceDir)
	if err != nil {
		fmt.Printf("Error finding .scss files: %v\n", err)
		return
	}

	// Track which files were updated and need to run through sass-migrator
	var filesToMigrate []string

	// Process each .scss file to replace alias imports with relative paths
	for _, file := range scssFiles {
		updated, err := replaceAliasWithRelativePath(file, *alias, *sourceDir)
		if err != nil {
			fmt.Printf("Error updating imports in file %s: %v\n", file, err)
			continue
		}
		if updated {
			filesToMigrate = append(filesToMigrate, file)
		}
	}

	// Run sass-migrator on each updated file
	for _, file := range filesToMigrate {
		fmt.Printf("Running sass-migrator on updated file: %s\n", file)
		err = runSassMigrator(file)
		if err != nil {
			fmt.Printf("Error running sass-migrator on file %s: %v\n", file, err)
		}
	}

	// Finally, run sass-migrator on the main entry file
	fmt.Printf("Running sass-migrator on entry file: %s\n", *entryFile)
	err = runSassMigrator(*entryFile)
	if err != nil {
		fmt.Printf("Error running sass-migrator on file %s: %v\n", *entryFile, err)
	}
}
