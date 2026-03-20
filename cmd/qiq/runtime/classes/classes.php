<?php

// TODO InternalIterator

// -------------------------------------- Exception -------------------------------------- MARK: Exception

// Spec: https://www.php.net/manual/en/class.exception.php
class Exception implements Throwable {
    /* Properties */
    protected string $message = "";
    private string $string = "";
    protected int $code;
    protected string $file = "";
    protected int $line;
    private array $trace = [];
    private ?Throwable $previous = null;

    /* Methods */
    public function __construct(string $message = "", int $code = 0, ?Throwable $previous = null) {
        $this->message = $message;
        $this->code = $code;
        $this->previous = $previous;
    }

    final public function getMessage(): string { return $this->message; }

    final public function getPrevious(): ?Throwable { return $this->previous; }

    final public function getCode(): int { return $this->code; }

    final public function getFile(): string { return $this->file; }

    final public function getLine(): int { return $this->line; }

    final public function getTrace(): array { return $this->trace; }

    final public function getTraceAsString(): string { return implode(PHP_EOL, $this->trace); }

    public function __toString(): string { return $this->string; }

    private function __clone(): void {}
}

// -------------------------------------- ErrorException -------------------------------------- MARK: ErrorException

// Spec: https://www.php.net/manual/en/class.errorexception.php
class ErrorException extends Exception {
    /* Properties */
    protected int $severity = E_ERROR;

    /* Methods */
    public function __construct(
    string $message = "",
    int $code = 0,
    int $severity = E_ERROR,
    ?string $filename = null,
    ?int $line = null,
    ?Throwable $previous = null
    ) {
        parent::__construct($message, $code, $previous);
        $this->severity = $severity;
        $this->file = $filename ?? "";
        $this->line = $line ?? 0;
    }

    final public function getSeverity(): int {
        return $this->severity;
    }
}

// -------------------------------------- Error -------------------------------------- MARK: Error

// Spec: https://www.php.net/manual/en/class.error.php
class Error implements Throwable {
    /* Properties */
    protected string $message = "";
    private string $string = "";
    protected int $code;
    protected string $file = "";
    protected int $line;
    private array $trace = [];
    private ?Throwable $previous = null;

    /* Methods */
    public function __construct(string $message = "", int $code = 0, ?Throwable $previous = null) {
        $this->message = $message;
        $this->code = $code;
        $this->previous = $previous;
    }

    final public function getMessage(): string { return $this->message; }

    final public function getPrevious(): ?Throwable { return $this->previous; }

    final public function getCode(): int { return $this->code; }

    final public function getFile(): string { return $this->file; }

    final public function getLine(): int { return $this->line; }

    final public function getTrace(): array { return $this->trace; }

    final public function getTraceAsString(): string { return implode(PHP_EOL, $this->trace); }

    public function __toString(): string { return $this->string; }

    private function __clone(): void {}
}

// TODO CompileError
// TODO ParseError
// TODO TypeError

// -------------------------------------- TypeError -------------------------------------- MARK: TypeError

// Spec: https://www.php.net/manual/en/class.typeerror.php
class TypeError extends Error {}

// -------------------------------------- ArgumentCountError -------------------------------------- MARK: ArgumentCountError

// Spec: https://www.php.net/manual/en/class.argumentcounterror.php
class ArgumentCountError extends TypeError {}

// -------------------------------------- ValueError -------------------------------------- MARK: ValueError

// Spec: https://www.php.net/manual/en/class.valueerror.php
class ValueError extends Error {}

// -------------------------------------- ArithmeticError -------------------------------------- MARK: ArithmeticError

// Spec: https://www.php.net/manual/en/class.arithmeticerror.php
class ArithmeticError extends Error {}

// TODO DivisionByZeroError
// TODO UnhandledMatchError
// TODO RequestParseBodyException
// TODO Closure
// TODO Generator

// -------------------------------------- ClosedGeneratorException -------------------------------------- MARK: ClosedGeneratorException

// Spec: https://www.php.net/manual/en/class.closedgeneratorexception.php
class ClosedGeneratorException extends Exception {}

// TODO WeakReference
// TODO WeakMap
// TODO Attribute
// TODO ReturnTypeWillChange
// TODO AllowDynamicProperties
// TODO SensitiveParameter
// TODO SensitiveParameterValue
// TODO Override
// TODO Deprecated
// TODO NoDiscard
// TODO DelayedTargetValidation
// TODO Fiber
// TODO FiberError

// -------------------------------------- stdClass -------------------------------------- MARK: stdClass

// Spec: https://www.php.net/manual/en/class.stdclass.php
class stdClass {
}

// TODO DateTime
// TODO DateTimeImmutable
// TODO DateTimeZone
// TODO DateInterval
// TODO DatePeriod
// TODO DateError
// TODO DateObjectError
// TODO DateRangeError
// TODO DateException
// TODO DateInvalidTimeZoneException
// TODO DateInvalidOperationException
// TODO DateMalformedStringException
// TODO DateMalformedIntervalStringException
// TODO DateMalformedPeriodStringException
// TODO BcMath\Number
// TODO Filter\FilterException
// TODO Filter\FilterFailedException
// TODO HashContext
// TODO JsonException
// TODO Random\RandomError
// TODO Random\BrokenRandomEngineError
// TODO Random\RandomException
// TODO Random\Engine\Mt19937
// TODO Random\Engine\PcgOneseq128XslRr64
// TODO Random\Engine\Xoshiro256StarStar
// TODO Random\Engine\Secure
// TODO Random\Randomizer
// TODO Random\IntervalBoundary
// TODO ReflectionException
// TODO Reflection
// TODO ReflectionFunctionAbstract
// TODO ReflectionFunction
// TODO ReflectionGenerator
// TODO ReflectionParameter
// TODO ReflectionType
// TODO ReflectionNamedType
// TODO ReflectionUnionType
// TODO ReflectionIntersectionType
// TODO ReflectionMethod
// TODO ReflectionClass
// TODO ReflectionObject
// TODO ReflectionProperty
// TODO ReflectionClassConstant
// TODO ReflectionExtension
// TODO ReflectionZendExtension
// TODO ReflectionReference
// TODO ReflectionAttribute
// TODO ReflectionEnum
// TODO ReflectionEnumUnitCase
// TODO ReflectionEnumBackedCase
// TODO ReflectionFiber
// TODO ReflectionConstant
// TODO PropertyHookType
// TODO Uri\Rfc3986\Uri
// TODO Uri\WhatWg\Url
// TODO Uri\UriComparisonMode
// TODO Uri\UriException
// TODO Uri\UriError
// TODO Uri\InvalidUriException
// TODO Uri\WhatWg\InvalidUrlException
// TODO Uri\WhatWg\UrlValidationError
// TODO Uri\WhatWg\UrlValidationErrorType
// TODO LogicException
// TODO BadFunctionCallException
// TODO BadMethodCallException
// TODO DomainException
// TODO InvalidArgumentException
// TODO LengthException
// TODO OutOfRangeException
// TODO RuntimeException
// TODO OutOfBoundsException
// TODO OverflowException
// TODO RangeException
// TODO UnderflowException
// TODO UnexpectedValueException
// TODO RecursiveIteratorIterator
// TODO IteratorIterator
// TODO FilterIterator
// TODO RecursiveFilterIterator
// TODO CallbackFilterIterator
// TODO RecursiveCallbackFilterIterator
// TODO ParentIterator
// TODO LimitIterator
// TODO CachingIterator
// TODO RecursiveCachingIterator
// TODO NoRewindIterator
// TODO AppendIterator
// TODO InfiniteIterator
// TODO RegexIterator
// TODO RecursiveRegexIterator
// TODO EmptyIterator
// TODO RecursiveTreeIterator
// TODO ArrayObject
// TODO ArrayIterator
// TODO RecursiveArrayIterator
// TODO SplFileInfo
// TODO DirectoryIterator
// TODO FilesystemIterator
// TODO RecursiveDirectoryIterator
// TODO GlobIterator
// TODO SplFileObject
// TODO SplTempFileObject
// TODO SplDoublyLinkedList
// TODO SplQueue
// TODO SplStack
// TODO SplHeap
// TODO SplMinHeap
// TODO SplMaxHeap
// TODO SplPriorityQueue
// TODO SplFixedArray
// TODO SplObjectStorage
// TODO MultipleIterator
// TODO SessionHandler
// TODO PhpToken
// TODO __PHP_Incomplete_Class
// TODO AssertionError
// TODO RoundingMode
// TODO php_user_filter
// TODO StreamBucket
// TODO Directory
