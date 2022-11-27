package templates

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/backdround/deploy-configs/pkg/fsutility"
	"github.com/backdround/deploy-configs/pkg/fstestutility"
)

func TestSuccessfulMakeTemplate(t *testing.T) {
	t.Run("OutputFileDoesntExist", func(t *testing.T) {
		// Creates a template file
		templateFile, cleanup :=
			fstestutility.CreateTemporaryFileWithData("{{.var1}} {{.var2}}")
		defer cleanup()

		expandData := map[string]string{
			"var1": "value1",
			"var2": "value2",
		}

		// Creates an output path
		outputPath := fstestutility.GetAvailableTempPath()
		defer os.Remove(outputPath)

		// Creates test data
		template := Template{
			Name:       "test-template",
			InputPath:  templateFile,
			OutputPath: outputPath,
			Data:       expandData,
		}

		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Success", containsString("test-template")).Once()

		// Executes the test
		NewTemplateMaker(logger).makeTemplate(template)

		// Asserts that the output file exists and expanded
		outputPathType := fsutility.GetFileType(outputPath)
		require.Equal(t, fsutility.Regular.String(), outputPathType.String())

		resultData, err := os.ReadFile(outputPath)
		fstestutility.AssertNoError(err)
		require.Equal(t, "value1 value2", string(resultData))
	})

	t.Run("OutputFileExists", func(t *testing.T) {
		// Creates a template file
		templateData := "{{.var1}} {{.var2}}"
		templateFile, templateCleanup :=
			fstestutility.CreateTemporaryFileWithData(templateData)
		defer templateCleanup()

		expandData := map[string]string{
			"var1": "value1",
			"var2": "value2",
		}

		// Creates an output file
		outputPath, outputCleanup :=
			fstestutility.CreateTemporaryFileWithData("some file data")
		defer outputCleanup()

		// Creates test data
		template := Template{
			Name:       "test-template",
			InputPath:  templateFile,
			OutputPath: outputPath,
			Data:       expandData,
		}

		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Success", containsString("test-template")).Once()

		// Executes the test
		NewTemplateMaker(logger).makeTemplate(template)

		// Asserts that the output file exists and expanded
		outputPathType := fsutility.GetFileType(outputPath)
		require.Equal(t, fsutility.Regular.String(), outputPathType.String())

		resultData, err := os.ReadFile(outputPath)
		fstestutility.AssertNoError(err)
		require.Equal(t, "value1 value2", string(resultData))
	})

	t.Run("OutputDirectoryDoesntExist", func(t *testing.T) {
		// Creates a template file
		templateFile, cleanup :=
			fstestutility.CreateTemporaryFileWithData("{{.var1}} {{.var2}}")
		defer cleanup()

		expandData := map[string]string{
			"var1": "value1",
			"var2": "value2",
		}

		// Creates an output file
		outputDirectory := fstestutility.GetAvailableTempPath()
		outputPath := path.Join(outputDirectory, "file")
		defer os.RemoveAll(outputDirectory)

		// Creates test data
		template := Template{
			Name:       "test-template",
			InputPath:  templateFile,
			OutputPath: outputPath,
			Data:       expandData,
		}

		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Success", containsString("test-template")).Once()

		// Executes the test
		NewTemplateMaker(logger).makeTemplate(template)

		// Asserts that the output file exists and expanded
		outputPathType := fsutility.GetFileType(outputPath)
		require.Equal(t, fsutility.Regular.String(), outputPathType.String())

		resultData, err := os.ReadFile(outputPath)
		fstestutility.AssertNoError(err)
		require.Equal(t, "value1 value2", string(resultData))
	})
}

func TestFailMakeTemplate(t *testing.T) {
	t.Run("TemplateFileDoesntExist", func(t *testing.T) {
		templatePath := fstestutility.GetAvailableTempPath()
		outputPath := fstestutility.GetAvailableTempPath()

		// Creates test data
		template := Template{
			Name:       "test-template",
			InputPath:  templatePath,
			OutputPath: outputPath,
			Data:       "some data",
		}

		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Fail", containsString("test-template")).Once()

		// Executes the test
		NewTemplateMaker(logger).makeTemplate(template)

		// Asserts that an output file doesn't exist
		outputPathType := fsutility.GetFileType(outputPath)
		require.Equal(t, fsutility.Notexisting.String(), outputPathType.String())
	})

	t.Run("DataDoesntCorrespond", func(t *testing.T) {
		// Creates a template file
		templateData := "{{.var1}} {{.var2}}"
		templateFile, cleanup :=
			fstestutility.CreateTemporaryFileWithData(templateData)
		defer cleanup()

		// Creates test data
		template := Template{
			Name:       "test-template",
			InputPath:  templateFile,
			OutputPath: fstestutility.GetAvailableTempPath(),
			Data:       "it doesn't have required fields",
		}

		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Fail", containsString("test-template")).Once()

		// Executes the test
		NewTemplateMaker(logger).makeTemplate(template)

		// Asserts that an output file doesn't exist
		outputPathType := fsutility.GetFileType(template.OutputPath)
		require.Equal(t, fsutility.Notexisting.String(), outputPathType.String())
	})

	t.Run("TemplateInvalid", func(t *testing.T) {
		// Creates a template file
		templateData := `{{`
		templateFile, cleanup :=
			fstestutility.CreateTemporaryFileWithData(templateData)
		defer cleanup()

		// Creates test data
		template := Template{
			Name:       "test-template",
			InputPath:  templateFile,
			OutputPath: fstestutility.GetAvailableTempPath(),
			Data:       "it is't required",
		}

		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Fail", containsString("test-template")).Once()

		// Executes the test
		NewTemplateMaker(logger).makeTemplate(template)

		// Asserts that a output file doesn't exist
		outputPathType := fsutility.GetFileType(template.OutputPath)
		require.Equal(t, fsutility.Notexisting.String(), outputPathType.String())
	})
}

func TestSkipMakeTemplate(t *testing.T) {
	// Creates the template file
	templateFile, templateCleanup :=
		fstestutility.CreateTemporaryFileWithData("{{.var1}} {{.var2}}")
	defer templateCleanup()

	expandData := map[string]string{
		"var1": "value1",
		"var2": "value2",
	}

	// Creates an output file
	outputFile, outputCleanup :=
		fstestutility.CreateTemporaryFileWithData("value1 value2")
	defer outputCleanup()

	// Creates test data
	template := Template{
		Name:       "test-template",
		InputPath:  templateFile,
		OutputPath: outputFile,
		Data:       expandData,
	}

	logger := &LoggerMock{}
	defer logger.AssertExpectations(t)
	logger.On("Log", containsString("test-template")).Once()

	// Executes the test
	NewTemplateMaker(logger).makeTemplate(template)

	// Asserts that the output file exists
	outputPathType := fsutility.GetFileType(outputFile)
	require.Equal(t, fsutility.Regular.String(), outputPathType.String())

	// Asserts that the file hasn't been changed
	resultData, err := os.ReadFile(outputFile)
	fstestutility.AssertNoError(err)
	require.Equal(t, "value1 value2", string(resultData))
}
