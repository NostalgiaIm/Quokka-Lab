import type { NoteEvent } from '@/types/audio';

export interface CollaborationSocket {
  close: () => void;
  sendNoteEvent: (event: NoteEvent) => void;
}

/**
 * 用于协作房间的最小 Action Cable 协议客户端。
 * 后续如果需要更完整的 Rails 生态能力，可以替换成 @rails/actioncable。
 */
export function useWebSocket(roomId: string, onNoteEvent: (event: NoteEvent) => void): CollaborationSocket {
  const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  const socket = new WebSocket(`${wsProtocol}//${window.location.host}/cable`);
  const identifier = JSON.stringify({ channel: 'CollaborationChannel', room_id: roomId });

  socket.addEventListener('open', () => {
    socket.send(JSON.stringify({ command: 'subscribe', identifier }));
  });

  socket.addEventListener('message', (message) => {
    const payload = JSON.parse(message.data);
    if (payload.type || !payload.message?.event) return;
    onNoteEvent(payload.message.event as NoteEvent);
  });

  function sendNoteEvent(event: NoteEvent): void {
    if (socket.readyState !== WebSocket.OPEN) return;

    socket.send(
      JSON.stringify({
        command: 'message',
        identifier,
        data: JSON.stringify({ action: 'receive', event }),
      }),
    );
  }

  return {
    close: () => socket.close(),
    sendNoteEvent,
  };
}
