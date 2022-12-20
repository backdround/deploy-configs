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
