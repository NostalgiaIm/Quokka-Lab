<?php

declare(strict_types=1);

namespace Quokka\Routing;

final class Router
{
    /** @var array<string, array<string, callable|array{class-string, string}>> 简易路由表。 */
    private array $routes = [];

    public function get(string $path, callable|array $handler): void
    {
        $this->routes['GET'][$path] = $handler;
    }

    public function post(string $path, callable|array $handler): void
    {
        $this->routes['POST'][$path] = $handler;
    }

    public function dispatch(string $method, string $path): void
    {
        $handler = $this->routes[$method][$path] ?? null;
        if ($handler === null) {
            $this->json(['error' => 'Not found'], 404);
            return;
        }

        if (is_array($handler)) {
            [$class, $methodName] = $handler;
            $handler = [new $class(), $methodName];
        }

        $handler();
    }

    private function json(array $payload, int $status): void
    {
        http_response_code($status);
        header('Content-Type: application/json');
        echo json_encode($payload, JSON_THROW_ON_ERROR);
    }
}
