package tests_test

import (
	"testing"

	"github.com/backdround/deploy-configs/tests/testcase"
)

func TestCommands(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("OutputDoentExist", func(t *testing.T) {
			initialFileTree := `
				.git:
				data.txt:
					type: file
					data: some data
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								commands:
									data-rev:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/data-rev.txt"
										command: "rev {{.Input}} > {{.Output}}"
			`
			resultFileTree := initialFileTree + `
				data-rev.txt:
					type: file
					data: atad emos
			`

			expectedSuccessMessage := `
				Command "data-rev" is executed:
					input: "{Root}/data.txt"
					output: "{Root}/data-rev.txt"
					command: "rev {{.Input}} > {{.Output}}"
			`

			c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, resultFileTree)
			c.RequireSuccessMessage(t, expectedSuccessMessage)
		})

		t.Run("OutputExists", func(t *testing.T) {
			initialFileTree := `
				.git:
				data.txt:
					type: file
					data: some data
				data-rev.txt:
					type: file
					data: unknown data
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								commands:
									data-rev:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/data-rev.txt"
										command: "rev {{.Input}} > {{.Output}}"
			`
			resultFileTree := `
				.git:
				data.txt:
					type: file
					data: some data
				data-rev.txt:
					type: file
					data: atad emos
				deploy-configs.yaml:
					type: file
			`

			expectedSuccessMessage := `
				Command "data-rev" is executed:
					input: "{Root}/data.txt"
					output: "{Root}/data-rev.txt"
					command: "rev {{.Input}} > {{.Output}}"
			`

			c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, resultFileTree)
			c.RequireSuccessMessage(t, expectedSuccessMessage)
		})
	})

	t.Run("Skip", func(t *testing.T) {
		fileTree := `
			.git:
			data.txt:
				type: file
				data: some data
			data-rev.txt:
				type: file
				data: atad emos
			deploy-configs.yaml:
				type: file
				data: |
					instances:
						pc1:
							commands:
								data-rev:
									input: "{{.GitRoot}}/data.txt"
									output: "{{.GitRoot}}/data-rev.txt"
									command: "rev {{.Input}} > {{.Output}}"
		`

		c := testcase.RunCase(t, fileTree, "./run", "pc1")
		c.RequireReturnCode(t, 0)
		c.RequireFileTree(t, fileTree)
		c.RequireLogMessage(t, `Command "data-rev" is skipped`)
	})

	t.Run("Fail", func(t *testing.T) {
		t.Run("InputDoesntExist", func(t *testing.T) {
			fileTree := `
				.git:
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								commands:
									data-rev:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/data-rev.txt"
										command: "rev {{.Input}} > {{.Output}}"
			`

			expectedFailMessage := `
				Unable to execute "data-rev" command:
					input: "{Root}/data.txt"
					output: "{Root}/data-rev.txt"
					command: "rev {{.Input}} > {{.Output}}"
						error: input file doesn't exist
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedFailMessage)
		})

		t.Run("ExecutionFail", func(t *testing.T) {
			fileTree := `
				.git:
				data.txt:
					type: file
					data: some data
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								commands:
									do:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/data-rev.txt"
										command: "./do {{.Input}} > {{.Output}}"
			`

			expectedFailMessage := `
				Unable to execute "do" command:
					input: "{Root}/data.txt"
					output: "{Root}/data-rev.txt"
					command: "./do {{.Input}} > {{.Output}}"
						error: exit status 127
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedFailMessage)
		})

		t.Run("OutputPathExists", func(t *testing.T) {
			fileTree := `
				.git:
				data.txt:
					type: file
					data: some data
				sub:
					type: file
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								commands:
									data-rev:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/sub/data-rev.txt"
										command: "rev {{.Input}} > {{.Output}}"
			`

			expectedFailMessage := `
				Unable to execute "data-rev" command:
					input: "{Root}/data.txt"
					output: "{Root}/sub/data-rev.txt"
					command: "rev {{.Input}} > {{.Output}}"
						error: unable to create directory
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedFailMessage)
		})

		t.Run("OutputPathIsADirectory", func(t *testing.T) {
			fileTree := `
				.git:
				data.txt:
					type: file
					data: some data
				data-rev:
					sub:
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								commands:
									data-rev:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/data-rev"
										command: "rev {{.Input}} > {{.Output}}"
			`

			expectedGeneralFailMessage := `
				Unable to execute "data-rev" command:
					input: "{Root}/data.txt"
					output: "{Root}/data-rev"
					command: "rev {{.Input}} > {{.Output}}"
						error: unable to replace output path
			`
			expectedSpecificFailMessage := "data-rev: directory not empty"

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedGeneralFailMessage)
			c.RequireFailMessage(t, expectedSpecificFailMessage)
		})
	})
}
