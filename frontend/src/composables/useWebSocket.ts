import type { NoteEvent } from '@/types/audio';

export interface CollaborationSocket {
  close: () => void;
  sendNoteEvent: (event: NoteEvent) => void;
}

/**
 * Creates a lightweight collaboration socket backed by the Go gateway.
 * The payload is intentionally plain JSON so future native/mobile clients can reuse it.
 */
export function useWebSocket(roomId: string, onNoteEvent: (event: NoteEvent) => void): CollaborationSocket {
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const socket = new WebSocket(`${wsProtocol}//${window.location.host}/ws/collaboration/${roomId}`);

  socket.addEventListener('message', (message) => {
    const payload = JSON.parse(message.data);
    if (!payload.event) return;
    onNoteEvent(payload.event as NoteEvent);
  });

  function sendNoteEvent(event: NoteEvent): void {
    if (socket.readyState !== WebSocket.OPEN) return;

    socket.send(JSON.stringify({ event }));
  }

  return {
    close: () => socket.close(),
    sendNoteEvent,
  };
}
