package errors

package: "errors"

#Error: {
	code:        string & =~"^ERR-[0-9]{3,4}$"
	name:        string & =~"^Err[A-Z][a-zA-Z0-9]+$"
	message:     string & !=""
	httpStatus:  *500 | int & >=100 & <=599
	description: *message | string
	severity:    *"medium" | "critical" | "high" | "low"
	parameters:  *[] | [...string]
}

// Category-specific error types
#AuthError: #Error & {
	httpStatus: *401 | int
	severity:   *"high" | string
}

#UserError: #Error & {
	httpStatus: *404 | int
	severity:   *"medium" | string
}

