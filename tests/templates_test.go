package tests_test

import (
	"testing"

	"github.com/backdround/deploy-configs/tests/testcase"
)

func TestTemplates(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("OutputDoentExist", func(t *testing.T) {
			initialFileTree := `
				.git:
				config.temp:
					type: file
					data: var = {{.var}}
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								templates:
									config:
										input: "{{.GitRoot}}/config.temp"
										output: "{{.GitRoot}}/config"
										data:
											var: 3
			`
			resultFileTree := initialFileTree + `
				config:
					type: file
					data: var = 3
			`

			expectedSuccessMessage := `
				Template "config" expanded:
					input: "{Root}/config.temp"
					output: "{Root}/config"
					data: map["var":'\x03']
			`

			c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, resultFileTree)
			c.RequireSuccessMessage(t, expectedSuccessMessage)
		})

		t.Run("OutputExists", func(t *testing.T) {
			initialFileTree := `
				.git:
				config.temp:
					type: file
					data: var = {{.var}}
				config:
					type: file
					data: old data
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								templates:
									config:
										input: "{{.GitRoot}}/config.temp"
										output: "{{.GitRoot}}/config"
										data:
											var: 3
			`
			resultFileTree := `
				.git:
				config.temp:
					type: file
					data: var = {{.var}}
				config:
					type: file
					data: var = 3
				deploy-configs.yaml:
					type: file
			`

			expectedSuccessMessage := `
				Template "config" expanded:
					input: "{Root}/config.temp"
					output: "{Root}/config"
					data: map["var":'\x03']
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
			config.temp:
				type: file
				data: var = {{.var}}
			config:
				type: file
				data: var = 3
			deploy-configs.yaml:
				type: file
				data: |
					instances:
						pc1:
							templates:
								config:
									input: "{{.GitRoot}}/config.temp"
									output: "{{.GitRoot}}/config"
									data:
										var: 3
		`

		c := testcase.RunCase(t, fileTree, "./run", "pc1")
		c.RequireReturnCode(t, 0)
		c.RequireFileTree(t, fileTree)
		c.RequireLogMessage(t, `Template "config" skipped`)
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
								templates:
									config:
										input: "{{.GitRoot}}/config.temp"
										output: "{{.GitRoot}}/config"
										data:
											var: 3
			`

			expectedFailMessage := `
				Unable to expand "config" template:
					input: "{Root}/config.temp"
					output: "{Root}/config"
					data: map["var":'\x03']
						error: input file doesn't exist
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedFailMessage)
		})

		t.Run("InvalidTemplate", func(t *testing.T) {
			fileTree := `
				.git:
				config.temp:
					type: file
					data: var = {{.var}
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								templates:
									config:
										input: "{{.GitRoot}}/config.temp"
										output: "{{.GitRoot}}/config"
										data:
											var: 3
			`

			expectedFailMessage := `
				Unable to expand "config" template:
					input: "{Root}/config.temp"
					output: "{Root}/config"
					data: map["var":'\x03']
						error: template: config.temp:1: bad character
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedFailMessage)
		})

		t.Run("MisspellingData", func(t *testing.T) {
			fileTree := `
				.git:
				config.temp:
					type: file
					data: var = {{.vvvar}}
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								templates:
									config:
										input: "{{.GitRoot}}/config.temp"
										output: "{{.GitRoot}}/config"
										data:
											var: 3
			`

			expectedGeneralFailMessage := `
				Unable to expand "config" template:
					input: "{Root}/config.temp"
					output: "{Root}/config"
					data: map["var":'\x03']
						error: template: config.temp:1:
			`
			expectedSpecificFailMessage := `
				map has no entry for key "vvvar"
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedGeneralFailMessage)
			c.RequireFailMessage(t, expectedSpecificFailMessage)
		})

		t.Run("OutputPathIsUnreachable", func(t *testing.T) {
			fileTree := `
				.git:
				config.temp:
					type: file
					data: var = {{.var}}
				sub:
					type: file
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								templates:
									config:
										input: "{{.GitRoot}}/config.temp"
										output: "{{.GitRoot}}/sub/config"
										data:
											var: 3
			`

			expectedGeneralFailMessage := `
				Unable to expand "config" template:
					input: "{Root}/config.temp"
					output: "{Root}/sub/config"
					data: map["var":'\x03']
						error: unable to create directory
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedGeneralFailMessage)
		})
	})
}
