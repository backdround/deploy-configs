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
						error: target path isn't exist
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedSuccessMessage)
		})

		t.Run("LinkPathIsAFile", func(t *testing.T) {
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
						error: link path is occupied
			`

			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedSuccessMessage)
		})

		t.Run("LinkPathIsUnreachable", func(t *testing.T) {
			fileTree := `
				.git:
				link.conf:
					type: file
				sub:
					type: file
				deploy-configs.yaml:
					type: file
					data: |
						instances:
							pc1:
								links:
									link1:
										target: "{{.GitRoot}}/link.conf"
										link: "{{.GitRoot}}/sub/link1"
			`
			expectedSuccessMessage := `
				Unable to create "link1" link:
					target: "{Root}/link.conf"
					link: "{Root}/sub/link1"
						error: unable to create directory
			`
			c := testcase.RunCase(t, fileTree, "./run", "pc1")
			c.RequireReturnCode(t, 1)
			c.RequireFileTree(t, fileTree)
			c.RequireFailMessage(t, expectedSuccessMessage)
		})
	})
}

func TestLinkDirectory(t *testing.T) {
	initialFileTree := `
		.git:
		services:
			service1:
				type: file
			service2:
				type: file
			service3:
				type: file
		deploy-configs.yaml:
			type: file
			data: |
				instances:
					pc1:
						links:
							services:
								target: "{{.GitRoot}}/services"
								link: "{{.GitRoot}}/deploy/services"
	`
	resultFileTree := initialFileTree + `
		deploy:
			services:
				service1:
					type: link
					path: ../../services/service1
				service2:
					type: link
					path: ../../services/service2
				service3:
					type: link
					path: ../../services/service3
	`

	expectedService1Message := `
		Link "services/service1" created:
			target: "{Root}/services/service1"
			link: "{Root}/deploy/services/service1"
	`
	expectedService2Message := `
		Link "services/service2" created:
			target: "{Root}/services/service2"
			link: "{Root}/deploy/services/service2"
	`
	expectedService3Message := `
		Link "services/service3" created:
			target: "{Root}/services/service3"
			link: "{Root}/deploy/services/service3"
	`

	c := testcase.RunCase(t, initialFileTree, "./run", "pc1")
	c.RequireReturnCode(t, 0)
	c.RequireFileTree(t, resultFileTree)

	c.RequireSuccessMessage(t, expectedService1Message)
	c.RequireSuccessMessage(t, expectedService2Message)
	c.RequireSuccessMessage(t, expectedService3Message)
}
