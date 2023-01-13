package tests_test

import (
	"testing"

	"github.com/backdround/deploy-configs/tests/testcase"
	"github.com/backdround/go-indent"
)

func TestInvalidArguments(t *testing.T) {
	c := testcase.RunCase(t, "", "./run")

	c.RequireReturnCode(t, 1)
	c.RequireFailMessage(t, "Expected config instance as argument")
}

func TestConfigNotFound(t *testing.T) {
	c := testcase.RunCase(t, "", "./run", "pc1")

	c.RequireReturnCode(t, 1)
	c.RequireFailMessage(t, "Error occurs while config searching")
}

func TestInvalidConfig(t *testing.T) {
	initialFileTree := `
		deploy-configs.yaml:
			type: file
			data: "\t"
	`
	c := testcase.RunCase(t, initialFileTree, "./run", "pc1")

	c.RequireReturnCode(t, 1)
	c.RequireFailMessage(t, "Fail to parse config data")
}

func TestComplex(t *testing.T) {
	deployConfigsYaml := `
		instances:
			pc1:
				links:
					link1:
						target: "{{.GitRoot}}/configs/link.conf"
						link: "{{.GitRoot}}/deploy/link1"
				templates:
					template1:
						input: "{{.GitRoot}}/configs/template.conf"
						output: "{{.GitRoot}}/deploy/template1"
						data:
							var: 3
				commands:
					command1:
						input: "{{.GitRoot}}/configs/command.conf"
						output: "{{.GitRoot}}/deploy/command1"
						command: "cat {{.Input}} > {{.Output}}"
	`

	initialFileTree := `
		.git:
		configs:
			link.conf:
				type: file
			template.conf:
				type: file
				data: "var = {{.var}}"
			command.conf:
				type: file
				data: "some data"

		deploy-configs.yaml:
			type: file
			data: |
	` + indent.Indent(deployConfigsYaml, "\t", 2)

	expectedFileTree := initialFileTree + `
		deploy:
			link1:
				type: link
				path: ../configs/link.conf
			template1:
				type: file
				data: "var = 3"
			command1:
				type: file
				data: "some data"
	`

	expectedLinkMessage := `
		Link "link1" created:
			target: "{Root}/configs/link.conf"
			link: "{Root}/deploy/link1"
	`

	expectedTemplateMessage := `
		Template "template1" expanded:
			input: "{Root}/configs/template.conf"
			output: "{Root}/deploy/template1"
	`

	expectedCommandMessage := `
		Command "command1" is executed:
			input: "{Root}/configs/command.conf"
			output: "{Root}/deploy/command1"
			command: "cat {{.Input}} > {{.Output}}"
	`

	c := testcase.RunCase(t, initialFileTree, "./run", "pc1")

	c.RequireReturnCode(t, 0)
	c.RequireFileTree(t, expectedFileTree)
	c.RequireSuccessMessage(t, expectedLinkMessage)
	c.RequireSuccessMessage(t, expectedTemplateMessage)
	c.RequireSuccessMessage(t, expectedCommandMessage)
}

func TestReplaceExistingOutputLink(t *testing.T) {
	t.Run("Links", func(t *testing.T) {
		initialFileTree := `
			.git:
			sources:
				file1:
					type: file
			deploy:
				link1:
					type: link
					path: ../file.txt
			deploy-configs.yaml:
				type: file
				data: |
					instances:
						pc1:
							links:
								link1:
									target: "{{.GitRoot}}/sources/file1"
									link: "{{.GitRoot}}/deploy/link1"
		`

		expectedFileTree := `
			.git:
			sources:
				file1:
					type: file
			deploy:
				link1:
					type: link
					path: ../sources/file1
			deploy-configs.yaml:
				type: file
		`

		expectedLinkMessage := `
			Link "link1" created:
				target: "{Root}/sources/file1"
				link: "{Root}/deploy/link1"
		`

		c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
		c.RequireReturnCode(t, 0)
		c.RequireFileTree(t, expectedFileTree)
		c.RequireSuccessMessage(t, expectedLinkMessage)
	})

	t.Run("Commands", func(t *testing.T) {
		initialFileTree := `
			.git:
			sources:
				file1:
					type: file
					data: "some data"
			deploy:
				command1:
					type: link
					path: ../file.txt
			deploy-configs.yaml:
				type: file
				data: |
					instances:
						pc1:
							commands:
								command1:
									input: "{{.GitRoot}}/sources/file1"
									output: "{{.GitRoot}}/deploy/command1"
									command: "cat {{.Input}} > {{.Output}}"
		`

		expectedFileTree := `
			.git:
			sources:
				file1:
					type: file
					data: "some data"
			deploy:
				command1:
					type: file
					data: "some data"
			deploy-configs.yaml:
				type: file
		`

		expectedCommandMessage := `
			Command "command1" is executed:
				input: "{Root}/sources/file1"
				output: "{Root}/deploy/command1"
				command: "cat {{.Input}} > {{.Output}}"
		`

		c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
		c.RequireReturnCode(t, 0)
		c.RequireFileTree(t, expectedFileTree)
		c.RequireSuccessMessage(t, expectedCommandMessage)
	})

	t.Run("Templates", func(t *testing.T) {
		initialFileTree := `
			.git:
			sources:
				file1:
					type: file
					data: "var = {{ .var }}"
			deploy:
				template1:
					type: link
					path: ../file.txt
			deploy-configs.yaml:
				type: file
				data: |
					instances:
						pc1:
							templates:
								template1:
									input: "{{.GitRoot}}/sources/file1"
									output: "{{.GitRoot}}/deploy/template1"
									data:
										var: 3
		`

		expectedFileTree := `
			.git:
			sources:
				file1:
					type: file
					data: "var = {{ .var }}"
			deploy:
				template1:
					type: file
					data: "var = 3"
			deploy-configs.yaml:
				type: file
		`

		expectedTemplateMessage := `
			Template "template1" expanded:
				input: "{Root}/sources/file1"
				output: "{Root}/deploy/template1"
		`

		c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
		c.RequireReturnCode(t, 0)
		c.RequireFileTree(t, expectedFileTree)
		c.RequireSuccessMessage(t, expectedTemplateMessage)
	})
}

func TestAlphabeticalExecutionSequence(t *testing.T) {
	t.Run("links", func(t *testing.T) {
			fileTree := `
				.git:
				sources:
					link1:
						type: file
					link2:
						type: file
					link3:
						type: file
				deploy:
					link1:
						type: link
						path: ../sources/link1
					link2:
						type: link
						path: ../sources/link2
					link3:
						type: link
						path: ../sources/link3
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									link3:
										target: "{{.GitRoot}}/sources/link3"
										link: "{{.GitRoot}}/deploy/link3"
									link2:
										target: "{{.GitRoot}}/sources/link2"
										link: "{{.GitRoot}}/deploy/link2"
									link1:
										target: "{{.GitRoot}}/sources/link1"
										link: "{{.GitRoot}}/deploy/link1"
			`
			expectedMessages := []string{
				`Link "link1" is skipped`,
				`Link "link2" is skipped`,
				`Link "link3" is skipped`,
			}
			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, fileTree)
			c.RequireLogMessages(t, expectedMessages, 2)
	})

	t.Run("linkDirectory", func(t *testing.T) {
			fileTree := `
				.git:
				sources:
					link1:
						type: file
					link2:
						type: file
					link3:
						type: file
				deploy:
					link1:
						type: link
						path: ../sources/link1
					link2:
						type: link
						path: ../sources/link2
					link3:
						type: link
						path: ../sources/link3
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									sources:
										target: "{{.GitRoot}}/sources"
										link: "{{.GitRoot}}/deploy"
			`
			expectedMessages := []string{
				`Link "sources/link1" is skipped`,
				`Link "sources/link2" is skipped`,
				`Link "sources/link3" is skipped`,
			}
			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, fileTree)
			c.RequireLogMessages(t, expectedMessages, 2)
	})

	t.Run("Commands", func(t *testing.T) {
			fileTree := `
				.git:
				data.txt:
					type: file
					data: some data
				rev1.txt:
					type: file
					data: atad emos
				rev2.txt:
					type: file
					data: atad emos
				rev3.txt:
					type: file
					data: atad emos
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								commands:
									reverse2:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/rev2.txt"
										command: "rev {{.Input}} > {{.Output}}"
									reverse1:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/rev1.txt"
										command: "rev {{.Input}} > {{.Output}}"
									reverse3:
										input: "{{.GitRoot}}/data.txt"
										output: "{{.GitRoot}}/rev3.txt"
										command: "rev {{.Input}} > {{.Output}}"
			`
			expectedMessages := []string{
				`Command "reverse1" is skipped`,
				`Command "reverse2" is skipped`,
				`Command "reverse3" is skipped`,
			}
			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, fileTree)
			c.RequireLogMessages(t, expectedMessages, 2)
	})

	t.Run("Templates", func(t *testing.T) {
			fileTree := `
				.git:
				template.txt:
					type: file
					data: var = {{.var}}
				result1.txt:
					type: file
					data: var = 1
				result2.txt:
					type: file
					data: var = 2
				result3.txt:
					type: file
					data: var = 3
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								templates:
									template3:
										input: "{{.GitRoot}}/template.txt"
										output: "{{.GitRoot}}/result3.txt"
										data:
											var: 3
									template1:
										input: "{{.GitRoot}}/template.txt"
										output: "{{.GitRoot}}/result1.txt"
										data:
											var: 1
									template2:
										input: "{{.GitRoot}}/template.txt"
										output: "{{.GitRoot}}/result2.txt"
										data:
											var: 2
			`
			expectedMessages := []string{
				`Template "template1" is skipped`,
				`Template "template2" is skipped`,
				`Template "template3" is skipped`,
			}
			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, fileTree)
			c.RequireLogMessages(t, expectedMessages, 2)
	})
}
