package commands

import (
	"os"
	"path"
	"testing"

	"github.com/backdround/deploy-configs/pkg/fstestutility"
	"github.com/backdround/deploy-configs/pkg/fsutility"
	"github.com/stretchr/testify/require"
)

func TestSuccessfulExecuteCommand(t *testing.T) {
	t.Run("OutputFileDoesntExist", func(t *testing.T) {
		// Creates input file
		inputFileData := "some data"
		inputFile, cleanup := fstestutility.
			CreateTemporaryFileWithData(inputFileData)
		defer cleanup()

		// Creates output path
		outputPath := fstestutility.GetAvailableTempPath()
		defer os.Remove(outputPath)

		// Creates test data
		command := Command{
			Name: "test-command",
			InputPath: inputFile,
			OutputPath: outputPath,
			Command: "cat {{.input}} > {{.output}}",
		}

		// Creates the logger mock
		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Success", containsString("test-command")).Once()

		// Executes the test
		NewCommandExecuter(logger).executeCommand(command)

		// Asserts output file
		outputFileType := fsutility.GetFileType(outputPath)
		require.Equal(t, fsutility.Regular.String(), outputFileType.String())

		outputFileData, err := os.ReadFile(outputPath)
		fstestutility.AssertNoError(err)
		require.Equal(t, inputFileData, string(outputFileData))
	})

	t.Run("OutputDirectoryDoesntExist", func(t *testing.T) {
		// Creates input file
		inputFileData := "some data"
		inputFile, cleanup := fstestutility.
			CreateTemporaryFileWithData(inputFileData)
		defer cleanup()

		// Creates output path
		outputPath := path.Join(fstestutility.GetAvailableTempPath(), "out.txt")
		defer os.Remove(outputPath)

		// Creates test data
		command := Command{
			Name: "test-command",
			InputPath: inputFile,
			OutputPath: outputPath,
			Command: "cat {{.input}} > {{.output}}",
		}

		// Creates the logger mock
		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Success", containsString("test-command")).Once()

		// Executes the test
		NewCommandExecuter(logger).executeCommand(command)

		// Asserts output file
		outputFileType := fsutility.GetFileType(outputPath)
		require.Equal(t, fsutility.Regular.String(), outputFileType.String())

		outputFileData, err := os.ReadFile(outputPath)
		fstestutility.AssertNoError(err)
		require.Equal(t, inputFileData, string(outputFileData))
	})

	t.Run("OutputFileExists", func(t *testing.T) {
		// Creates input file
		inputFileData := "some data"
		inputFile, cleanup := fstestutility.
			CreateTemporaryFileWithData(inputFileData)
		defer cleanup()

		// Creates output file
		outputFile := "output_file.*.txt"
		cleanup = fstestutility.CreateTemporaryFiles(&outputFile)
		defer cleanup()

		// Creates test data
		command := Command{
			Name: "test-command",
			InputPath: inputFile,
			OutputPath: outputFile,
			Command: "cat {{.input}} > {{.output}}",
		}

		// Creates the logger mock
		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Success", containsString("test-command")).Once()

		// Executes the test
		NewCommandExecuter(logger).executeCommand(command)

		// Asserts output file
		outputFileType := fsutility.GetFileType(outputFile)
		require.Equal(t, fsutility.Regular.String(), outputFileType.String())

		outputFileData, err := os.ReadFile(outputFile)
		fstestutility.AssertNoError(err)
		require.Equal(t, inputFileData, string(outputFileData))
	})
}

func TestFailedExecuteCommand(t *testing.T) {
	t.Run("InputFileDoenstExist", func(t *testing.T) {
		// Creates test data
		command := Command{
			Name: "test-command",
			InputPath: fstestutility.GetAvailableTempPath(),
			OutputPath: fstestutility.GetAvailableTempPath(),
			Command: "cat {{.input}} > {{.output}}",
		}

		// Creates the logger mock
		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Fail", containsString("test-command")).Once()

		// Executes the test
		NewCommandExecuter(logger).executeCommand(command)

		// Asserts output file
		outputFileType := fsutility.GetFileType(command.OutputPath)
		require.Equal(t, fsutility.Notexisting.String(), outputFileType.String())
	})

	t.Run("CommandInvalid", func(t *testing.T) {
		// Creates input file
		inputFile, cleanup := fstestutility.
			CreateTemporaryFileWithData("some data")
		defer cleanup()

		// Creates test data
		command := Command{
			Name: "test-command",
			InputPath: inputFile,
			OutputPath: fstestutility.GetAvailableTempPath(),
			Command: "cattt {{.input}} > {{.output}}",
		}

		// Creates the logger mock
		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Fail", containsString("test-command")).Once()

		// Executes the test
		NewCommandExecuter(logger).executeCommand(command)

		// Asserts output file
		outputFileType := fsutility.GetFileType(command.OutputPath)
		require.Equal(t, fsutility.Notexisting.String(), outputFileType.String())
	})

	t.Run("CommandDoesntCreateOutputFile", func(t *testing.T) {
		// Creates input file
		inputFileData := "some data"
		inputFile, cleanup := fstestutility.
			CreateTemporaryFileWithData(inputFileData)
		defer cleanup()

		// Creates test data
		command := Command{
			Name: "test-command",
			InputPath: inputFile,
			OutputPath: fstestutility.GetAvailableTempPath(),
			Command: "cat {{.input}} {{.output}}",
		}

		// Creates the logger mock
		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Fail", containsString("test-command")).Once()

		// Executes the test
		NewCommandExecuter(logger).executeCommand(command)
	})

	t.Run("UnableToCreateOutputDirectoryDueToFileInPath", func(t *testing.T) {
		// Creates input file
		inputFileData := "some data"
		inputFile, cleanup := fstestutility.
			CreateTemporaryFileWithData(inputFileData)
		defer cleanup()

		// Creates output path
		fileInPathData := "data"
		fileInPath, cleanup := fstestutility.CreateTemporaryFileWithData(
				fileInPathData)
		defer cleanup()
		outputPath := path.Join(fileInPath, "out.txt")

		// Creates test data
		command := Command{
			Name: "test-command",
			InputPath: inputFile,
			OutputPath: outputPath,
			Command: "cat {{.input}} > {{.output}}",
		}

		// Creates the logger mock
		logger := &LoggerMock{}
		defer logger.AssertExpectations(t)
		logger.On("Fail", containsString("test-command")).Once()

		// Executes the test
		NewCommandExecuter(logger).executeCommand(command)

		// Asserts that file in path wasn't changed
		fileInPathType := fsutility.GetFileType(fileInPath)
		require.Equal(t, fsutility.Regular.String(), fileInPathType.String())

		fileInPathResultData, err := os.ReadFile(fileInPath)
		fstestutility.AssertNoError(err)
		require.Equal(t, fileInPathData, string(fileInPathResultData))
	})
}

func TestSkippedExecuteCommand(t *testing.T) {
	// Creates input file
	inputFileData := "some data"
	inputFile, cleanup := fstestutility.
		CreateTemporaryFileWithData(inputFileData)
	defer cleanup()

	// Creates output file
	outputFile, cleanup := fstestutility.CreateTemporaryFileWithData("some data")
	defer cleanup()

	// Creates test data
	command := Command{
		Name: "test-command",
		InputPath: inputFile,
		OutputPath: outputFile,
		Command: "cat {{.input}} > {{.output}}",
	}

	// Creates the logger mock
	logger := &LoggerMock{}
	defer logger.AssertExpectations(t)
	logger.On("Log", containsString("test-command")).Once()

	// Executes the test
	NewCommandExecuter(logger).executeCommand(command)

	// Asserts output file
	outputFileType := fsutility.GetFileType(outputFile)
	require.Equal(t, fsutility.Regular.String(), outputFileType.String())

	outputFileData, err := os.ReadFile(outputFile)
	fstestutility.AssertNoError(err)
	require.Equal(t, inputFileData, string(outputFileData))
}
