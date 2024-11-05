package exception

var Messages = map[string]string{
	"RT000": "unknown declaration: %T",
	"RT001": "variable %s already exists",
	"RT002": "undefined variable: %s",
	"RT003": "invalid array index: %v",
	"RT004": "array index out of bounds: %v",
	"RT005": "cannot assign to property of type: %T",
	"RT006": "invalid assignment target: %T",
	"RT007": "type mismatch: %T + %T",
	"RT008": "invalid operation: division by zero",
	"RT009": "invalid operand type for modulus: %T",
	"RT010": "unknown binary expression: %s",
	"RT011": "cannot get length of: %T",
	"RT012": "illegal use of range. type %T is invalid",
	"RT013": "unknown unary expression: %s",
	"RT014": "invalid function call %v",
	"RT015": "unknown logical operator: %s",
	"RT016": "cannot access property of type: %T",
	"RT017": "unknown expression: %T",
	"RT018": "unknown namespace: %s",
	"RT019": "invalid namespace: %T",
	"RT020": "invalid member expression property",
	"RT021": "control variable not found in environment: %s",
	"RT022": "loop stop should be a number, got: %T",
	"RT023": "loop step should be a number, got: %T",
	"RT024": "loop condition should be a boolean, got: %T",
	"RT025": "cannot make %s public. identifier does not exits",
	"GN001": "cannot find module '%s'",
}
