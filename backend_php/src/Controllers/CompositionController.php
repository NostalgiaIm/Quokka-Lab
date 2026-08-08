<?php

declare(strict_types=1);

namespace Quokka\Controllers;

final class CompositionController
{
    private string $databasePath;

    public function __construct()
    {
        $this->databasePath = __DIR__ . '/../../storage/compositions.json';
    }

    public function index(): void
    {
        $this->json($this->readAll());
    }

    public function store(): void
    {
        $payload = $_POST['composition'] ?? $_POST;
        $title = $payload['title'] ?? 'Untitled melody';
        $composition = [
            'id' => time(),
            'title' => $title,
            'audio_url' => $this->storeAudioUpload(),
            'midi_data' => json_decode($payload['midi_data'] ?? '{}', true) ?: [],
            'bpm' => (int) ($payload['bpm'] ?? 120),
            'key_signature' => $payload['key_signature'] ?? 'C',
            'is_public' => filter_var($payload['is_public'] ?? true, FILTER_VALIDATE_BOOLEAN),
            'created_at' => gmdate(DATE_ATOM),
            'updated_at' => gmdate(DATE_ATOM),
        ];

        $items = $this->readAll();
        array_unshift($items, $composition);
        file_put_contents($this->databasePath, json_encode($items, JSON_PRETTY_PRINT | JSON_THROW_ON_ERROR));

        $this->json($composition, 201);
    }

    private function storeAudioUpload(): ?string
    {
        if (!isset($_FILES['audio']) && !isset($_FILES['composition']['tmp_name']['audio'])) {
            return null;
        }

        $file = $_FILES['audio'] ?? [
            'tmp_name' => $_FILES['composition']['tmp_name']['audio'],
            'name' => $_FILES['composition']['name']['audio'],
        ];

        $safeName = preg_replace('/[^a-zA-Z0-9._-]/', '-', basename($file['name']));
        $target = __DIR__ . '/../../public/uploads/' . time() . '-' . $safeName;
        move_uploaded_file($file['tmp_name'], $target);

        return '/uploads/' . basename($target);
    }

    /** @return array<int, array<string, mixed>> 返回本地 JSON 文件中的作品列表。 */
    private function readAll(): array
    {
        if (!is_file($this->databasePath)) {
            return [];
        }

        return json_decode((string) file_get_contents($this->databasePath), true) ?: [];
    }

    private function json(array $payload, int $status = 200): void
    {
        http_response_code($status);
        header('Content-Type: application/json');
        echo json_encode($payload, JSON_THROW_ON_ERROR);
    }
}
