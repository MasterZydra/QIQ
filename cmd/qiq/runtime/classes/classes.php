<?php

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

// Spec: https://www.php.net/manual/en/class.stdclass.php
class stdClass {
}
