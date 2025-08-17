package common

import (
	"slices"
	"strings"
)

func IsReservedName(token string) bool {
	return IsKeyword(token) || IsContextDependentConstant(token) || IsCorePredefinedConstant(token) || slices.Contains([]string{"parent"}, strings.ToLower(token))
}

// MARK: Keywords

// Spec: https://www.php.net/manual/en/reserved.keywords.php
// TODO handle __halt_compiler()

// Spec: https://phplang.org/spec/19-grammar.html#grammar-keyword
var keywords = []string{
	"abstract", "and", "array", "as", "break", "callable", "case", "catch", "class", "clone",
	"const", "continue", "declare", "default", "die", "do", "echo", "else", "elseif", "empty",
	"enddeclare", "endfor", "endforeach", "endif", "endswitch", "endwhile", "eval", "exit",
	"extends", "final", "finally", "fn", "for", "foreach", "function", "global",
	"goto", "if", "implements", "include", "include_once", "instanceof",
	"insteadof", "interface", "isset", "list", "match", "namespace", "new", "or", "print", "private",
	"protected", "public", "readonly", "require", "require_once", "return", "static", "switch",
	"throw", "trait", "try", "unset", "use", "var", "while", "xor", "yield", "yield from",
	// Non-spec:
	"mixed", "void",
}

func IsKeyword(token string) bool {
	// Spec: https://phplang.org/spec/19-grammar.html#grammar-keyword

	// keyword:: one of
	//    abstract   and   array   as   break   callable   case   catch   class   clone
	//    const   continue   declare   default   die   do   echo   else   elseif   empty
	//    enddeclare   endfor   endforeach   endif   endswitch   endwhile   eval   exit
	//    extends   final   finally   for   foreach   function   global
	//    goto   if   implements   include   include_once   instanceof
	//    insteadof   interface   isset   list   namespace   new   or   print   private
	//    protected   public   require   require_once   return   static   switch
	//    throw   trait   try   unset   use   var   while   xor   yield   yield from

	// Spec: https://phplang.org/spec/09-lexical-structure.html#keywords
	// Keywords are not case-sensitive.
	token = strings.ToLower(token)

	return slices.Contains(keywords, token)
}

// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type
var castTypeKeywords = []string{
	"array", "binary", "bool", "boolean", "double", "int", "integer", "float", "object", "real", "string",
}

func IsCastTypeKeyword(token string) bool {
	// Spec: https://phplang.org/spec/10-expressions.html#grammar-cast-type

	// Spec: https://phplang.org/spec/09-lexical-structure.html#keywords
	// Keywords are not case-sensitive.
	token = strings.ToLower(token)

	return slices.Contains(castTypeKeywords, token)

}

// Spec: https://phplang.org/spec/13-functions.html#grammar-base-type-declaration
var paramTypeKeywords = []string{
	"mixed", "array", "bool", "float", "int", "null", "string",
}

func IsParamTypeKeyword(token string) bool {
	// Spec: https://phplang.org/spec/13-functions.html#grammar-base-type-declaration

	// Spec: https://phplang.org/spec/09-lexical-structure.html#keywords
	// Keywords are not case-sensitive.
	token = strings.ToLower(token)

	return slices.Contains(paramTypeKeywords, token)
}

func IsReturnTypeKeyword(token string) bool {
	// Spec: https://phplang.org/spec/13-functions.html#grammar-return-type

	// Spec: https://phplang.org/spec/09-lexical-structure.html#keywords
	// Keywords are not case-sensitive.
	token = strings.ToLower(token)

	return token == "void" || slices.Contains(paramTypeKeywords, token)
}

// Spec: https://phplang.org/spec/14-classes.html#grammar-visibility-modifier
var visibilityModifierKeywords = []string{"public", "protected", "private"}

func IsVisibilitModifierKeyword(token string) bool {
	// Spec: https://phplang.org/spec/14-classes.html#grammar-visibility-modifier

	// Spec: https://phplang.org/spec/09-lexical-structure.html#keywords
	// Keywords are not case-sensitive.
	token = strings.ToLower(token)

	return slices.Contains(visibilityModifierKeywords, token)
}

// Spec: https://phplang.org/spec/14-classes.html#grammar-class-modifier
var classModifierKeywords = []string{"abstract", "final"}

func IsClassModifierKeyword(token string) bool {
	// Spec: https://phplang.org/spec/14-classes.html#grammar-class-modifier

	// Spec: https://phplang.org/spec/09-lexical-structure.html#keywords
	// Keywords are not case-sensitive.
	token = strings.ToLower(token)

	return slices.Contains(classModifierKeywords, token)
}

// MARK: Constants

// Spec: https://phplang.org/spec/06-constants.html#context-dependent-constants
var contextDependentConstants = []string{
	"__CLASS__", "__COMPILER_HALT_OFFSET__", "__DIR__", "__FILE__", "__FUNCTION__", "__LINE__",
	"__METHOD__", "__NAMESPACE__", "__TRAIT__",
}

func IsContextDependentConstant(token string) bool {
	// Spec: https://phplang.org/spec/06-constants.html#context-dependent-constants

	// The following constants—sometimes referred to as magic constants—are automatically available to all scripts;
	// their values are not fixed and they are case-insensitive:
	token = strings.ToLower(token)

	return slices.Contains(contextDependentConstants, token)
}

// Spec: https://phplang.org/spec/06-constants.html#core-predefined-constants
var corePredefinedConstants = []string{
	"DEFAULT_INCLUDE_PATH", "E_ALL", "E_COMPILE_ERROR", "E_COMPILE_WARNING", "E_CORE_ERROR", "E_CORE_WARNING", "E_DEPRECATED",
	"E_ERROR", "E_NOTICE", "E_PARSE", "E_RECOVERABLE_ERROR", "E_STRICT", "E_USER_DEPRECATED", "E_USER_ERROR", "E_USER_NOTICE",
	"E_USER_WARNING", "E_WARNING", "FALSE", "INF", "M_1_PI", "M_2_PI", "M_2_SQRTPI", "M_E", "M_EULER", "M_LN10", "M_LN2",
	"M_LNPI", "M_LOG10E", "M_LOG2E", "M_PI", "M_PI_2", "M_PI_4", "M_SQRT1_2", "M_SQRT2", "M_SQRT3", "M_SQRTPI", "NAN", "NULL",
	"PHP_BINARY", "PHP_BINDIR", "PHP_CONFIG_FILE_PATH", "PHP_CONFIG_FILE_SCAN_DIR", "PHP_DEBUG", "PHP_EOL", "PHP_EXTENSION_DIR",
	"PHP_EXTRA_VERSION", "PHP_INT_MAX", "PHP_INT_MIN", "PHP_INT_SIZE", "PHP_FLOAT_DIG", "PHP_FLOAT_EPSILON", "PHP_FLOAT_MIN",
	"PHP_FLOAT_MAX", "PHP_MAJOR_VERSION", "PHP_MANDIR", "PHP_MAXPATHLEN", "PHP_MINOR_VERSION", "PHP_OS", "PHP_OS_FAMILY",
	"PHP_PREFIX", "PHP_RELEASE_VERSION", "PHP_ROUND_HALF_DOWN", "PHP_ROUND_HALF_EVEN", "PHP_ROUND_HALF_ODD", "PHP_ROUND_HALF_UP",
	"PHP_SAPI", "PHP_SHLIB_SUFFIX", "PHP_SYSCONFDIR", "PHP_VERSION", "PHP_VERSION_ID", "PHP_ZTS", "STDIN", "STDOUT", "STDERR",
	"TRUE",
}

func IsCorePredefinedConstant(token string) bool {
	// Spec: https://phplang.org/spec/06-constants.html#core-predefined-constants

	// The following constants are automatically available to all scripts;
	// they are case-sensitive with the exception of NULL, TRUE and FALSE

	if slices.Contains([]string{"NULL", "TRUE", "FALSE"}, strings.ToUpper(token)) {
		return true
	}
	return slices.Contains(corePredefinedConstants, token)
}
