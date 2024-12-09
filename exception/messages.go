package exception

var Messages = map[string]string{
	"RT000": "unknown declaration: %s",
	"RT001": "variable %s already exists",
	"RT002": "undefined variable: %s",
	"RT003": "invalid index: %v",
	"RT004": "array index out of bounds: %v",
	"RT005": "cannot assign to property of type: %s",
	"RT006": "cannot assign property to %s type",
	"RT007": "cannot sum value of type %s with a type %s",
	"RT008": "invalid operation: division by zero",
	"RT009": "cannot apply modulo operator of type %s and type %s",
	"RT010": "unknown binary expression: %s",
	"RT011": "cannot get length of type %s",
	"RT012": "cannot get range of type %s",
	"RT013": "unknown unary expression: %s",
	"RT014": "invalid function call %v",
	"RT015": "unknown logical operator: %s",
	"RT016": "cannot access property of type %s",
	"RT017": "unknown expression: %s",
	"RT018": "unknown namespace: %s",
	"RT019": "invalid namespace: %s",
	"RT020": "invalid member expression property",
	"RT021": "control variable not found in environment: %s",
	"RT022": "loop stop should be a <num> got %s instead",
	"RT023": "loop step should be a <num> got %s instead",
	"RT024": "loop condition should be a <bool> got %s instead",
	"RT025": "cannot make %s public. identifier does not exits",
	"RT026": "cannot operate comparison with type %s",
	"RT027": "cannot subtract value of type %s with a type %s",
	"RT028": "cannot convert value of type %s to type %s",
	"RT029": "<str> should have only one char to be converted to type <char>",
	"RT030": "convert into char failed",
	"RT031": "invalid internal function call",
	"RT032": "cannot %s file: %s",
	"RT033": "all paths must be strings",
	"RT034": "enum member '%s' does not exits in '%s'",
	"RT035": "cannot use '%s' as a type",
	"RT036": "enum member expects to be called with arguments",
	"RT037": "left operand in enumof should be an enum member",
	"RT038": "right operand in enumof should be an enum member",
	"RT039": "cannot destruct type %s",
	"GN001": "cannot find module '%s'",
	"GN002": "unable to load file '%s'. module does not exits. path: %s",
}
