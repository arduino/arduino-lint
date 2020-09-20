package projects

import (
	"fmt"
	"os"

	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/projects/library"
	"github.com/arduino/arduino-check/projects/packageindex"
	"github.com/arduino/arduino-check/projects/platform"
	"github.com/arduino/arduino-check/projects/projecttype"
	"github.com/arduino/arduino-check/projects/sketch"
	"github.com/arduino/arduino-cli/cli/errorcodes"
	"github.com/arduino/go-paths-helper"
)

type Type struct {
	Path             *paths.Path
	ProjectType      projecttype.Type
	SuperprojectType projecttype.Type
}

func FindProjects() []Type {
	targetPath := configuration.TargetPath()
	superprojectTypeConfiguration := configuration.SuperprojectType()
	recursive := configuration.Recursive()

	var foundProjects []Type

	// If targetPath is a file, targetPath itself is the project, so it's only necessary to determine/verify the type
	if targetPath.IsNotDir() {
		// The filename provides additional information about the project type. So rather than using isProject(), which doesn't make use this information, use a specialized function that does.
		isProject, projectType := isProjectIndicatorFile(targetPath, superprojectTypeConfiguration)
		if isProject {
			foundProject := Type{
				Path:             targetPath.Parent(),
				ProjectType:      projectType,
				SuperprojectType: projectType,
			}
			foundProjects = append(foundProjects, foundProject)

			foundProjects = append(foundProjects, findSubprojects(foundProject, projectType)...)

			return foundProjects
		}

		fmt.Errorf("error: specified path %s is not an Arduino project", targetPath.String())
		os.Exit(errorcodes.ErrGeneric)
	}

	foundProjects = append(foundProjects, findProjects(targetPath, superprojectTypeConfiguration, recursive)...)

	if foundProjects == nil {
		fmt.Errorf("error: no projects found under %s", targetPath.String())
		os.Exit(errorcodes.ErrGeneric)
	}

	return foundProjects
}

func findProjects(targetPath *paths.Path, projectType projecttype.Type, recursive bool) []Type {
	var foundProjects []Type

	isProject, projectType := isProject(targetPath, projectType)
	if isProject {
		foundProject := Type{
			Path:        targetPath,
			ProjectType: projectType,
			// findSubprojects() will overwrite this with the correct value when the project is a subproject
			SuperprojectType: projectType,
		}
		foundProjects = append(foundProjects, foundProject)

		foundProjects = append(foundProjects, findSubprojects(foundProject, foundProject.ProjectType)...)

		// Don't search recursively past a project
		return foundProjects
	}

	if recursive {
		// targetPath was not a project, so search the subfolders
		directoryListing, _ := targetPath.ReadDir()
		directoryListing.FilterDirs()
		for _, potentialProjectDirectory := range directoryListing {
			foundProjects = append(foundProjects, findProjects(potentialProjectDirectory, projectType, recursive)...)
		}
	}

	return foundProjects
}

func findSubprojects(superproject Type, apexSuperprojectType projecttype.Type) []Type {
	var immediateSubprojects []Type

	switch superproject.ProjectType {
	case projecttype.Sketch:
		// Sketches don't have subprojects
		return nil
	case projecttype.Library:
		subprojectPath := superproject.Path.Join("examples")
		immediateSubprojects = append(immediateSubprojects, findProjects(subprojectPath, projecttype.Sketch, true)...)
		// Apparently there is some level of official support for "example" in addition to the specification-compliant "examples"
		// see: https://github.com/arduino/arduino-cli/blob/0.13.0/arduino/libraries/loader.go#L153
		subprojectPath = superproject.Path.Join("example")
		immediateSubprojects = append(immediateSubprojects, findProjects(subprojectPath, projecttype.Sketch, true)...)
	case projecttype.Platform:
		subprojectPath := superproject.Path.Join("libraries")
		immediateSubprojects = append(immediateSubprojects, findProjects(subprojectPath, projecttype.Library, false)...)
	case projecttype.PackageIndex:
		// Platform indexes don't have subprojects
		return nil
	}

	var allSubprojects []Type
	// Subprojects may have their own subprojects
	for _, immediateSubproject := range immediateSubprojects {
		// Subprojects at all levels should have SuperprojectType set to the top level superproject's type, not the immediate parent's type
		immediateSubproject.SuperprojectType = apexSuperprojectType
		// Each parent project should be followed in the list by its subprojects
		allSubprojects = append(allSubprojects, immediateSubproject)
		allSubprojects = append(allSubprojects, findSubprojects(immediateSubproject, apexSuperprojectType)...)
	}

	return allSubprojects
}

// isProject determines if a path contains an Arduino project, and if so which type
func isProject(potentialProjectPath *paths.Path, projectType projecttype.Type) (bool, projecttype.Type) {
	if (projectType == projecttype.All || projectType == projecttype.Sketch) && isSketch(potentialProjectPath) {
		return true, projecttype.Sketch
	} else if (projectType == projecttype.All || projectType == projecttype.Library) && isLibrary(potentialProjectPath) {
		return true, projecttype.Library
	} else if (projectType == projecttype.All || projectType == projecttype.Platform) && isPlatform(potentialProjectPath) {
		return true, projecttype.Platform
	} else if (projectType == projecttype.All || projectType == projecttype.PackageIndex) && isPackageIndex(potentialProjectPath) {
		return true, projecttype.PackageIndex
	}
	return false, projecttype.Not
}

// isProject determines if a file is the indicator file for an Arduino project, and if so which type
func isProjectIndicatorFile(potentialProjectFilePath *paths.Path, projectType projecttype.Type) (bool, projecttype.Type) {
	if (projectType == projecttype.All || projectType == projecttype.Sketch) && isSketchIndicatorFile(potentialProjectFilePath) {
		return true, projecttype.Sketch
	} else if (projectType == projecttype.All || projectType == projecttype.Library) && isLibraryIndicatorFile(potentialProjectFilePath) {
		return true, projecttype.Library
	} else if (projectType == projecttype.All || projectType == projecttype.Platform) && isPlatformIndicatorFile(potentialProjectFilePath) {
		return true, projecttype.Platform
	} else if (projectType == projecttype.All || projectType == projecttype.PackageIndex) && isPackageIndexIndicatorFile(potentialProjectFilePath) {
		return true, projecttype.PackageIndex
	}
	return false, projecttype.Not
}

// isSketch determines if a path is an Arduino sketch
// Note: this intentionally does not determine the validity of the sketch, only that the developer's intent was for it to be a sketch
func isSketch(potentialProjectPath *paths.Path) bool {
	directoryListing, _ := potentialProjectPath.ReadDir()
	directoryListing.FilterOutDirs()
	for _, potentialSketchFile := range directoryListing {
		if isSketchIndicatorFile(potentialSketchFile) {
			return true
		}
	}

	// No file was found with a valid main sketch file extension
	return false
}

func isSketchIndicatorFile(filePath *paths.Path) bool {
	if sketch.HasMainFileValidExtension(filePath) {
		return true
	}
	return false
}

// isLibrary determines if a path is an Arduino library
// Note: this intentionally does not determine the validity of the library, only that the developer's intent was for it to be a library
func isLibrary(potentialProjectPath *paths.Path) bool {
	// Arduino libraries will always have one of the following files in its root folder:
	// - a library.properties metadata file
	// - a header file
	directoryListing, _ := potentialProjectPath.ReadDir()
	directoryListing.FilterOutDirs()
	for _, potentialLibraryFile := range directoryListing {
		if isLibraryIndicatorFile(potentialLibraryFile) {
			return true
		}
	}

	// None of the files required for a valid Arduino library were found
	return false
}

func isLibraryIndicatorFile(filePath *paths.Path) bool {
	if library.IsMetadataFile(filePath) {
		return true
	}

	if library.HasHeaderFileValidExtension(filePath) {
		return true
	}

	return false
}

// isPlatform determines if a path is an Arduino boards platform
// Note: this intentionally does not determine the validity of the platform, only that the developer's intent was for it to be a platform
func isPlatform(potentialProjectPath *paths.Path) bool {
	directoryListing, _ := potentialProjectPath.ReadDir()
	directoryListing.FilterOutDirs()
	for _, potentialPlatformFile := range directoryListing {
		if isStrictPlatformIndicatorFile(potentialPlatformFile) {
			return true
		}
	}

	return false
}

func isPlatformIndicatorFile(filePath *paths.Path) bool {
	if platform.IsConfigurationFile(filePath) {
		return true
	}

	return false
}

func isStrictPlatformIndicatorFile(filePath *paths.Path) bool {
	if platform.IsRequiredConfigurationFile(filePath) {
		return true
	}

	return false
}

// isPackageIndex determines if a path contains an Arduino package index
// Note: this intentionally does not determine the validity of the package index, only that the developer's intent was for it to be a package index
func isPackageIndex(potentialProjectPath *paths.Path) bool {
	directoryListing, _ := potentialProjectPath.ReadDir()
	directoryListing.FilterOutDirs()
	for _, potentialPackageIndexFile := range directoryListing {
		if isStrictPackageIndexIndicatorFile(potentialPackageIndexFile) {
			return true
		}
	}

	return false
}

func isPackageIndexIndicatorFile(filePath *paths.Path) bool {
	if filePath.Ext() == ".json" {
		return true
	}

	return false
}

func isStrictPackageIndexIndicatorFile(filePath *paths.Path) bool {
	if packageindex.HasValidFilename(filePath, true) {
		return true
	}

	return false
}
