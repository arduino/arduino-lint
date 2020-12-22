source = ["dist/arduino-lint_osx_darwin_amd64/arduino-lint"]
bundle_id = "cc.arduino.arduino-lint"

sign {
  application_identity = "Developer ID Application: ARDUINO SA (7KT7ZWMCJT)"
}

# Ask Gon for zip output to force notarization process to take place.
# The CI will ignore the zip output, using the signed binary only.
zip {
  output_path = "arduino-lint.zip"
}
