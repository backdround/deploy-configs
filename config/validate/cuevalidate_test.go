package validate

import (
	"testing"

	"fmt"
	"github.com/lithammer/dedent"
	"strings"
)

func assertNoTab(data string) {
	if strings.Contains(data, "\t") {
		message := fmt.Sprintf("data contains tab: \n%v", data)
		panic(message)
	}
}

func TestAllowMissingFields(t *testing.T) {
	data := dedent.Dedent(`
	  instances:
	    instance1:
	      links:
	`)
	assertNoTab(data)

	err := Validate([]byte(data))

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

	err := Validate([]byte(data))

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

	err := Validate([]byte(data))

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
		        target: ./file1.txt
		        link: ./link
		`)
		assertNoTab(data)

		err := Validate([]byte(data))

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

		err := Validate([]byte(data))
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

		err := Validate([]byte(data))
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
	        link1:
	          target: ./file1.txt
	          link: ./link
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
	            simple-variable: "value"
	            nested-variable:
	              nested1: "value"
	              nested2: "value"
	`)
	assertNoTab(data)

	err := Validate([]byte(data))

	if err != nil {
		t.Errorf("expected no errors:\n%v\ndata:%v\n", err, data)
	}
}
