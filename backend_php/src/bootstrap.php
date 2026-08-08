<?php

declare(strict_types=1);

spl_autoload_register(function (string $class): void {
    $prefix = 'Quokka\\';
    if (!str_starts_with($class, $prefix)) {
        return;
    }

    $relative = str_replace('\\', DIRECTORY_SEPARATOR, substr($class, strlen($prefix))) . '.php';
    $path = __DIR__ . DIRECTORY_SEPARATOR . $relative;

    if (is_file($path)) {
        require $path;
    }
});

if (!is_dir(__DIR__ . '/../public/uploads')) {
    mkdir(__DIR__ . '/../public/uploads', 0775, true);
}
