// Package project finds and classifies Arduino projects.
package project

import (
	"fmt"

	"github.com/arduino/arduino-check/configuration"
	"github.com/arduino/arduino-check/project/library"
	"github.com/arduino/arduino-check/project/packageindex"
	"github.com/arduino/arduino-check/project/platform"
	"github.com/arduino/arduino-check/project/projecttype"
	"github.com/arduino/arduino-check/project/sketch"
	"github.com/arduino/go-paths-helper"
	"github.com/sirupsen/logrus"
)

// Type is the type for project definitions.
type Type struct {
	Path             *paths.Path
	ProjectType      projecttype.Type
	SuperprojectType projecttype.Type
}

// FindProjects searches the target path configured by the user for projects of the type configured by the user as well as the subprojects of those project.
// It returns a slice containing the definitions of each found project.
func FindProjects() ([]Type, error) {
	targetPath := configuration.TargetPath()
	superprojectTypeFilter := configuration.SuperprojectTypeFilter()
	recursive := configuration.Recursive()

	var foundProjects []Type

	// If targetPath is a file, targetPath itself is the project, so it's only necessary to determine/verify the type.
	if targetPath.IsNotDir() {
		logrus.Debug("Projects path is file")
		// The filename provides additional information about the project type. So rather than using isProject(), which doesn't make use this information, use a specialized function that does.
		isProject, projectType := isProjectIndicatorFile(targetPath, superprojectTypeFilter)
		if isProject {
			foundProject := Type{
				Path:             targetPath.Parent(),
				ProjectType:      projectType,
				SuperprojectType: projectType,
			}
			foundProjects = append(foundProjects, foundProject)

			foundProjects = append(foundProjects, findSubprojects(foundProject, projectType)...)

			return foundProjects, nil
		}

		return nil, fmt.Errorf("specified path %s is not an Arduino project", targetPath)
	}

	foundProjects = append(foundProjects, findProjectsUnderPath(targetPath, superprojectTypeFilter, recursive)...)

	if foundProjects == nil {
		return nil, fmt.Errorf("no projects found under %s", targetPath)
	}

	return foundProjects, nil
}

// findProjectsUnderPath finds projects of the given type and subprojects of those projects. It returns a slice containing the definitions of all found projects.
func findProjectsUnderPath(targetPath *paths.Path, projectType projecttype.Type, recursive bool) []Type {
	var foundProjects []Type

	isProject, foundProjectType := isProject(targetPath, projectType)
	if isProject {
		logrus.Tracef("%s is %s", targetPath, projectType)
		foundProject := Type{
			Path:        targetPath,
			ProjectType: foundProjectType,
			// findSubprojects() will overwrite this with the correct value when the project is a subproject.
			SuperprojectType: foundProjectType,
		}
		foundProjects = append(foundProjects, foundProject)

		foundProjects = append(foundProjects, findSubprojects(foundProject, foundProject.ProjectType)...)

		// Don't search recursively past a project.
		return foundProjects
	}

	if recursive {
		// targetPath was not a project, so search the subfolders.
		directoryListing, _ := targetPath.ReadDir()
		directoryListing.FilterDirs()
		for _, potentialProjectDirectory := range directoryListing {
			foundProjects = append(foundProjects, findProjectsUnderPath(potentialProjectDirectory, projectType, recursive)...)
		}
	}

	return foundProjects
}

// findSubprojects finds subprojects of the given project.
// For example, the subprojects of a library are its example sketches.
func findSubprojects(superproject Type, apexSuperprojectType projecttype.Type) []Type {
	subprojectFolderNames := []string{}
	var subProjectType projecttype.Type
	var searchPathsRecursively bool

	// Determine possible subproject paths
	switch superproject.ProjectType {
	case projecttype.Sketch:
		// Sketches don't have subprojects
		return nil
	case projecttype.Library:
		subprojectFolderNames = append(subprojectFolderNames, library.ExamplesFolderSupportedNames()...)
		subProjectType = projecttype.Sketch
		searchPathsRecursively = true // Examples sketches can be under nested subfolders
	case projecttype.Platform:
		subprojectFolderNames = append(subprojectFolderNames, platform.BundledLibrariesFolderNames()...)
		subProjectType = projecttype.Library
		searchPathsRecursively = false // Bundled libraries must be in the root of the libraries folder
	case projecttype.PackageIndex:
		// Platform indexes don't have subprojects
		return nil
	default:
		panic(fmt.Sprintf("Subproject discovery not configured for project type: %s", superproject.ProjectType))
	}

	// Search the subproject paths for projects
	var immediateSubprojects []Type
	for _, subprojectFolderName := range subprojectFolderNames {
		subprojectPath := superproject.Path.Join(subprojectFolderName)
		immediateSubprojects = append(immediateSubprojects, findProjectsUnderPath(subprojectPath, subProjectType, searchPathsRecursively)...)
	}

	var allSubprojects []Type
	// Subprojects may have their own subprojects.
	for _, immediateSubproject := range immediateSubprojects {
		// Subprojects at all levels should have SuperprojectType set to the top level superproject's type, not the immediate parent's type.
		immediateSubproject.SuperprojectType = apexSuperprojectType
		// Each parent project should be followed in the list by its subprojects.
		allSubprojects = append(allSubprojects, immediateSubproject)
		allSubprojects = append(allSubprojects, findSubprojects(immediateSubproject, apexSuperprojectType)...)
	}

	return allSubprojects
}

// isProject determines if a path contains an Arduino project, and if so which type.
func isProject(potentialProjectPath *paths.Path, projectType projecttype.Type) (bool, projecttype.Type) {
	logrus.Tracef("Checking if %s is %s", potentialProjectPath, projectType)
	if (projectType == projecttype.All || projectType == projecttype.Sketch) && isSketch(potentialProjectPath) {
		logrus.Tracef("%s is %s", potentialProjectPath, projecttype.Sketch)
		return true, projecttype.Sketch
	} else if (projectType == projecttype.All || projectType == projecttype.Library) && isLibrary(potentialProjectPath) {
		logrus.Tracef("%s is %s", potentialProjectPath, projecttype.Library)
		return true, projecttype.Library
	} else if (projectType == projecttype.All || projectType == projecttype.Platform) && isPlatform(potentialProjectPath) {
		logrus.Tracef("%s is %s", potentialProjectPath, projecttype.Platform)
		return true, projecttype.Platform
	} else if (projectType == projecttype.All || projectType == projecttype.PackageIndex) && isPackageIndex(potentialProjectPath) {
		logrus.Tracef("%s is %s", potentialProjectPath, projecttype.PackageIndex)
		return true, projecttype.PackageIndex
	}
	return false, projecttype.Not
}

// isProject determines if a file is the indicator file for an Arduino project, and if so which type.
func isProjectIndicatorFile(potentialProjectFilePath *paths.Path, projectType projecttype.Type) (bool, projecttype.Type) {
	logrus.Tracef("Checking if %s is %s indicator file", potentialProjectFilePath, projectType)
	if (projectType == projecttype.All || projectType == projecttype.Sketch) && isSketchIndicatorFile(potentialProjectFilePath) {
		logrus.Tracef("%s is %s indicator file", potentialProjectFilePath, projecttype.Sketch)
		return true, projecttype.Sketch
	} else if (projectType == projecttype.All || projectType == projecttype.Library) && isLibraryIndicatorFile(potentialProjectFilePath) {
		logrus.Tracef("%s is %s indicator file", potentialProjectFilePath, projecttype.Library)
		return true, projecttype.Library
	} else if (projectType == projecttype.All || projectType == projecttype.Platform) && isPlatformIndicatorFile(potentialProjectFilePath) {
		logrus.Tracef("%s is %s indicator file", potentialProjectFilePath, projecttype.Platform)
		return true, projecttype.Platform
	} else if (projectType == projecttype.All || projectType == projecttype.PackageIndex) && isPackageIndexIndicatorFile(potentialProjectFilePath) {
		logrus.Tracef("%s is %s indicator file", potentialProjectFilePath, projecttype.PackageIndex)
		return true, projecttype.PackageIndex
	}
	logrus.Tracef("%s is not indicator file", potentialProjectFilePath)
	return false, projecttype.Not
}

// isSketch determines whether a path is an Arduino sketch.
// Note: this intentionally does not determine the validity of the sketch, only that the developer's intent was for it to be a sketch.
func isSketch(potentialProjectPath *paths.Path) bool {
	directoryListing, _ := potentialProjectPath.ReadDir()
	directoryListing.FilterOutDirs()
	for _, potentialSketchFile := range directoryListing {
		if isSketchIndicatorFile(potentialSketchFile) {
			return true
		}
	}

	// No file was found with a valid main sketch file extension.
	return false
}

func isSketchIndicatorFile(filePath *paths.Path) bool {
	return sketch.HasMainFileValidExtension(filePath)
}

// isLibrary determines if a path is an Arduino library.
// Note: this intentionally does not determine the validity of the library, only that the developer's intent was for it to be a library.
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

	// None of the files required for a valid Arduino library were found.
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

// isPlatform determines if a path is an Arduino boards platform.
// Note: this intentionally does not determine the validity of the platform, only that the developer's intent was for it to be a platform.
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
	return platform.IsConfigurationFile(filePath)
}

func isStrictPlatformIndicatorFile(filePath *paths.Path) bool {
	return platform.IsRequiredConfigurationFile(filePath)
}

// isPackageIndex determines if a path contains an Arduino package index.
// Note: this intentionally does not determine the validity of the package index, only that the developer's intent was for it to be a package index.
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
	return packageindex.HasValidExtension(filePath)
}

func isStrictPackageIndexIndicatorFile(filePath *paths.Path) bool {
	return packageindex.HasValidFilename(filePath, true)
}
