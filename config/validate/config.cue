// List of symbolic lists to create
#Links: {
	[string]: [string, string]
}

// List of commands to execute
#Commands: {
	[string]: {
		input:   string
		output:  string
		command: string
	}
}

// List of tmeplates to evaluate
#Templates: {
	[string]: {
		input:  string
		output: string
		data:   _
	}
}

// Instances of config
#Instances: [string]: {
	links?:     #Links | null
	commands?:  #Commands | null
	templates?: #Templates | null
}

// Top level dictionary of instances
instances: #Instances
