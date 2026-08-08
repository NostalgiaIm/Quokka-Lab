<?php

declare(strict_types=1);

use Quokka\Controllers\CompositionController;
use Quokka\Controllers\HealthController;
use Quokka\Routing\Router;

require __DIR__ . '/../src/bootstrap.php';

$router = new Router();
$router->get('/health', [HealthController::class, 'show']);
$router->get('/api/v1/compositions', [CompositionController::class, 'index']);
$router->post('/api/v1/compositions', [CompositionController::class, 'store']);

$router->dispatch($_SERVER['REQUEST_METHOD'], parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH) ?: '/');
