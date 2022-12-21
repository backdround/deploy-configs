package tests_test

import (
	"testing"

	"github.com/backdround/deploy-configs/tests/testcase"
)

func TestInvalidPathSubstitution(t *testing.T) {
	t.Run("Link", func(t *testing.T) {
		fileTree := `
			.git:
			link.conf:
				type: file
			deploy-configs.yaml:
				type: file
				data: |
					instances:
						pc1:
							links:
								link1:
									target: "{{.GitRut}}/link.conf"
									link: "{{.GitRoot}}/link1"
		`

		expectedGeneralMessage := `Invalid config links:`
		expectedSpecificMessage := `
			unable to expand "link1" link:
				{{.GitRut}}/link.conf
					template: path-expander:
		`
		expectedSubSpecificMessage := `map has no entry for key "GitRut"`

		c := testcase.RunCase(t, fileTree, "./run", "pc1")
		c.RequireReturnCode(t, 1)
		c.RequireFileTree(t, fileTree)
		c.RequireFailMessage(t, expectedGeneralMessage)
		c.RequireFailMessage(t, expectedSpecificMessage)
		c.RequireFailMessage(t, expectedSubSpecificMessage)
	})

	t.Run("Template", func(t *testing.T) {
		fileTree := `
			.git:
			template.conf:
				type: file
				data: var = {{.var}}
			deploy-configs.yaml:
				type: file
				data: |
					instances:
						pc1:
							templates:
								template1:
									input: "{{.GitRut}}/template.conf"
									output: "{{.GitRoot}}/template"
									data:
										var: 3
		`

		expectedGeneralMessage := `Invalid config templates:`
		expectedSpecificMessage := `
			unable to expand "template1" template:
				{{.GitRut}}/template.conf
					template: path-expander:
		`
		expectedSubSpecificMessage := `map has no entry for key "GitRut"`

		c := testcase.RunCase(t, fileTree, "./run", "pc1")
		c.RequireReturnCode(t, 1)
		c.RequireFileTree(t, fileTree)
		c.RequireFailMessage(t, expectedGeneralMessage)
		c.RequireFailMessage(t, expectedSpecificMessage)
		c.RequireFailMessage(t, expectedSubSpecificMessage)
	})

	t.Run("Command", func(t *testing.T) {
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
								command1:
									input: "{{.GitRut}}/data.txt"
									output: "{{.GitRoot}}/result-data.txt"
									command: "cat {{.Input}} > {{.Output}}"
		`

		expectedGeneralMessage := `Invalid config commands:`
		expectedSpecificMessage := `
			unable to expand "command1" command:
				{{.GitRut}}/data.txt
					template: path-expander:
		`
		expectedSubSpecificMessage := `map has no entry for key "GitRut"`

		c := testcase.RunCase(t, fileTree, "./run", "pc1")
		c.RequireReturnCode(t, 1)
		c.RequireFileTree(t, fileTree)
		c.RequireFailMessage(t, expectedGeneralMessage)
		c.RequireFailMessage(t, expectedSpecificMessage)
		c.RequireFailMessage(t, expectedSubSpecificMessage)
	})

	t.Run("CommandCommandField", func(t *testing.T) {
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
								command1:
									input: "{{.GitRoot}}/data.txt"
									output: "{{.GitRoot}}/result-data.txt"
									command: "cat {{.Inpuz}} > {{.Output}}"
		`

		expectedGeneralMessage := `
			Unable to execute "command1" command:
				input: "{Root}/data.txt"
				output: "{Root}/result-data.txt"
				command: "cat {{.Inpuz}} > {{.Output}}"
					error: template: command1:
		`
		expectedSpecificMessage := `map has no entry for key "Inpuz"`

		c := testcase.RunCase(t, fileTree, "./run", "pc1")
		c.RequireReturnCode(t, 1)
		c.RequireFileTree(t, fileTree)
		c.RequireFailMessage(t, expectedGeneralMessage)
		c.RequireFailMessage(t, expectedSpecificMessage)
	})
}
