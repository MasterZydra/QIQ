<?php

// Spec: https://www.php.net/manual/en/class.traversable.php
interface Traversable {
}

// Spec: https://www.php.net/manual/en/class.iteratoraggregate.php
interface IteratorAggregate extends Traversable {
    /* Methods */
    public function getIterator(): Traversable;
}

// Spec: https://www.php.net/manual/en/class.iterator.php
interface Iterator extends Traversable {
    /* Methods */
    public function current(): mixed;
    public function key(): mixed;
    public function next(): void;
    public function rewind(): void;
    public function valid(): bool;
}

// Spec: https://www.php.net/manual/en/class.serializable.php
interface Serializable {
    /* Methods */
    public function serialize(): ?string;
    public function unserialize(string $data): void;
}

// Spec: https://www.php.net/manual/en/class.arrayaccess.php
interface ArrayAccess {
    /* Methods */
    public function offsetExists(mixed $offset): bool;
    public function offsetGet(mixed $offset): mixed;
    public function offsetSet(mixed $offset, mixed $value): void;
    public function offsetUnset(mixed $offset): void;
}

// Spec: https://www.php.net/manual/en/class.countable.php
interface Countable {
    /* Methods */
    public function count(): int;
}

// Spec: https://www.php.net/manual/en/class.stringable.php
interface Stringable {
    /* Methods */
    public function __toString(): string;
}

// Spec: https://www.php.net/manual/en/class.throwable.php
interface Throwable extends Stringable {
    /* Methods */
    public function getMessage(): string;
    public function getCode(): int;
    public function getFile(): string;
    public function getLine(): int;
    public function getTrace(): array;
    public function getTraceAsString(): string;
    public function getPrevious(): ?Throwable;
}

// Spec: https://www.php.net/manual/en/class.unitenum.php
interface UnitEnum {
    /* Methods */
    public static function cases(): array;
}
