package tests_test

import (
	"testing"

	"github.com/backdround/deploy-configs/tests/testcase"
)

func TestLinks(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Run("LinkDoenstExist", func(t *testing.T) {
			initialFileTree := `
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
										target: "{{.GitRoot}}/link.conf"
										link: "{{.GitRoot}}/link1"
			`
			resultFileTree := initialFileTree + `
				link1:
					type: link
					path: ./link.conf
			`

			expectedSuccessMessage := `
				Link "link1" created:
					target: "{Root}/link.conf"
					link: "{Root}/link1"
			`

			c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, resultFileTree)
			c.RequireSuccessMessage(t, expectedSuccessMessage)
		})

		t.Run("LinkPointsToDifferentDestination", func(t *testing.T) {
			initialFileTree := `
				.git:
				link.conf:
					type: file
				link1:
					type: link
					path: ./another.conf
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									link1:
										target: "{{.GitRoot}}/link.conf"
										link: "{{.GitRoot}}/link1"
			`
			resultFileTree := `
				.git:
				link.conf:
					type: file
				link1:
					type: link
					path: ./link.conf
				deploy-configs.yaml:
					type: file
			`

			expectedSuccessMessage := `
				Link "link1" created:
					target: "{Root}/link.conf"
					link: "{Root}/link1"
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
				link.conf:
					type: file
				link1:
					type: link
					path: ./link.conf
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									link1:
										target: "./link.conf"
										link: "{{.GitRoot}}/link1"
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 0)
			c.RequireFileTree(t, fileTree)
			c.RequireLogMessage(t, `Link "link1" skipped`)
	})

	t.Run("Fail", func(t *testing.T) {
		t.Run("TargetDoesntExist", func(t *testing.T) {
			fileTree := `
				.git:
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									link1:
										target: "{{.GitRoot}}/link.conf"
										link: "{{.GitRoot}}/link1"
			`

			expectedSuccessMessage := `
				Unable to create "link1" link:
					target: "{Root}/link.conf"
					link: "{Root}/link1"
						error: Target file isn't exist
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedSuccessMessage)
		})

		t.Run("DestinationExists", func(t *testing.T) {
			fileTree := `
				.git:
				link.conf:
					type: file
				link1:
					type: file
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									link1:
										target: "{{.GitRoot}}/link.conf"
										link: "{{.GitRoot}}/link1"
			`

			expectedSuccessMessage := `
				Unable to create "link1" link:
					target: "{Root}/link.conf"
					link: "{Root}/link1"
						error: Link file already exists
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedSuccessMessage)
		})

		t.Run("DestinationPathExists", func(t *testing.T) {
			fileTree := `
				.git:
				link.conf:
					type: file
				link-path:
					type: file
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									link1:
										target: "{{.GitRoot}}/link.conf"
										link: "{{.GitRoot}}/link-path/link1"
			`
			expectedSuccessMessage := `
				Unable to create "link1" link:
					target: "{Root}/link.conf"
					link: "{Root}/link-path/link1"
						error: Link path already exists
			`
			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedSuccessMessage)
		})
	})
}