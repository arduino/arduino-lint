// This file is part of arduino-check.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-check.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

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
	var foundProjects []Type

	for _, targetPath := range configuration.TargetPaths() {
		foundProjectsForTargetPath, err := findProjects(targetPath)
		if err != nil {
			return nil, err
		}
		foundProjects = append(foundProjects, foundProjectsForTargetPath...)
	}
	return foundProjects, nil
}

func findProjects(targetPath *paths.Path) ([]Type, error) {
	var foundProjects []Type

	// If targetPath is a file, targetPath itself is the project, so it's only necessary to determine/verify the type.
	if targetPath.IsNotDir() {
		logrus.Debug("Projects path is file")
		// The filename provides additional information about the project type. So rather than using isProject(), which doesn't make use this information, use a specialized function that does.
		isProject, projectType := isProjectIndicatorFile(targetPath, configuration.SuperprojectTypeFilter())
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

	foundProjects = append(foundProjects, findProjectsUnderPath(targetPath, configuration.SuperprojectTypeFilter(), configuration.Recursive())...)

	if foundProjects == nil {
		return nil, fmt.Errorf("no projects found under %s", targetPath)
	}

	return foundProjects, nil
}

// findProjectsUnderPath finds projects of the given type and subprojects of those projects. It returns a slice containing the definitions of all found projects.
func findProjectsUnderPath(targetPath *paths.Path, projectTypeFilter projecttype.Type, recursive bool) []Type {
	var foundProjects []Type

	isProject, foundProjectType := isProject(targetPath, projectTypeFilter)
	if isProject {
		logrus.Tracef("%s is %s", targetPath, foundProjectType)
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
			foundProjects = append(foundProjects, findProjectsUnderPath(potentialProjectDirectory, projectTypeFilter, recursive)...)
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
func isProject(potentialProjectPath *paths.Path, projectTypeFilter projecttype.Type) (bool, projecttype.Type) {
	logrus.Tracef("Checking if %s is %s", potentialProjectPath, projectTypeFilter)

	projectType := projecttype.Not
	if projectTypeFilter.Matches(projecttype.Sketch) && isSketch(potentialProjectPath) {
		projectType = projecttype.Sketch
	} else if projectTypeFilter.Matches(projecttype.Library) && isLibrary(potentialProjectPath) {
		projectType = projecttype.Library
	} else if projectTypeFilter.Matches(projecttype.Platform) && isPlatform(potentialProjectPath) {
		projectType = projecttype.Platform
	} else if projectTypeFilter.Matches(projecttype.PackageIndex) && isPackageIndex(potentialProjectPath) {
		projectType = projecttype.PackageIndex
	}

	if projectType == projecttype.Not {
		return false, projectType
	}
	logrus.Tracef("%s is %s", potentialProjectPath, projectType)
	return true, projectType
}

// isProject determines if a file is the indicator file for an Arduino project, and if so which type.
func isProjectIndicatorFile(potentialProjectFilePath *paths.Path, projectTypeFilter projecttype.Type) (bool, projecttype.Type) {
	logrus.Tracef("Checking if %s is %s indicator file", potentialProjectFilePath, projectTypeFilter)

	projectType := projecttype.Not
	if projectTypeFilter.Matches(projecttype.Sketch) && isSketchIndicatorFile(potentialProjectFilePath) {
		projectType = projecttype.Sketch
	} else if projectTypeFilter.Matches(projecttype.Library) && isLibraryIndicatorFile(potentialProjectFilePath) {
		projectType = projecttype.Library
	} else if projectTypeFilter.Matches(projecttype.Platform) && isPlatformIndicatorFile(potentialProjectFilePath) {
		projectType = projecttype.Platform
	} else if projectTypeFilter.Matches(projecttype.PackageIndex) && isPackageIndexIndicatorFile(potentialProjectFilePath) {
		projectType = projecttype.PackageIndex
	}

	if projectType == projecttype.Not {
		logrus.Tracef("%s is not indicator file", potentialProjectFilePath)
		return false, projectType
	}
	logrus.Tracef("%s is %s indicator file", potentialProjectFilePath, projectType)
	return true, projectType
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
