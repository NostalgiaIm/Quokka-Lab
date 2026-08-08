<?php

declare(strict_types=1);

namespace Quokka\Controllers;

final class HealthController
{
    public function show(): void
    {
        $this->json(['status' => 'ok', 'service' => 'quokka-php-api']);
    }

    private function json(array $payload, int $status = 200): void
    {
        http_response_code($status);
        header('Content-Type: application/json');
        echo json_encode($payload, JSON_THROW_ON_ERROR);
    }
}
