package config

import (
	"testing"

	"github.com/lithammer/dedent"
)

func TestAllowMissingFields(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      links:
	`)
	assertNoTab(data)

	err := cueValidate([]byte(data))

	if err != nil {
		t.Errorf("expected no errors:\n%v\ndata:%v\n", err, data)
	}
}

func TestAllowEmptyFields(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      links:
	      commands:
	      templates:
	`)
	assertNoTab(data)

	err := cueValidate([]byte(data))

	if err != nil {
		t.Errorf("expected no errors:\n%v\ndata:%v\n", err, data)
	}
}

func TestAllowNotInstancesFields(t *testing.T) {
	data := dedent.Dedent(`
	  some-dictionary:
	    data: "data"
	  instances:
	    instance1:
	      links:
	`)
	assertNoTab(data)

	err := cueValidate([]byte(data))

	if err != nil {
		t.Errorf("expected no errors:\n%v\ndata:%v\n", err, data)
	}
}

func TestInvalidData(t *testing.T) {
	t.Run("InvalidLink", func(t *testing.T) {
		data := dedent.Dedent(`
		  instances:
		    instance:
		      links:
		        - ["./file1.txt", "./link"]
		`)
		assertNoTab(data)

		err := cueValidate([]byte(data))

		if err == nil {
			t.Error("expected error on incorrect link data")
		}
	})

	t.Run("InvalidCommand", func(t *testing.T) {
		data := dedent.Dedent(`
		  instances:
		    instance:
		      commands:
		        command1:
		          output: "~/file.txt"
		          command: "cat {{input}} > {{output}}"
		`)
		assertNoTab(data)

		err := cueValidate([]byte(data))
		if err == nil {
			t.Error("expected error on incorrect command data")
		}
	})

	t.Run("InvalidTemplate", func(t *testing.T) {
		data := dedent.Dedent(`
		  instances:
		    instance:
		      templates:
		        template1:
		          output: "~/file.txt"
		          data:
		            variable1: "value1"
		            variable2: "value2"
		`)
		assertNoTab(data)

		err := cueValidate([]byte(data))
		if err == nil {
			t.Error("expected error on incorrect template data")
		}
	})
}

func TestSimpleData(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      links:
	        link1: ["./file1.txt", "./link"]
	      commands:
	        command1:
	          input: "./file.txt"
	          output: "~/file.txt"
	          command: "cat {{input}} > {{output}}"
	      templates:
	        template1:
	          input: "./file.txt"
	          output: "~/file.txt"
	          data:
	            variable1: "value1"
	            variable2: "value2"
	`)
	assertNoTab(data)

	err := cueValidate([]byte(data))

	if err != nil {
		t.Errorf("expected no errors:\n%v\ndata:%v\n", err, data)
	}
}
