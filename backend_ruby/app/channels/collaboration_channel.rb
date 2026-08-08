class CollaborationChannel < ApplicationCable::Channel
  def subscribed
    stream_from stream_name
  end

  def receive(data)
    # 广播 note_on / note_off 事件，让同一房间的用户同步键盘动作。
    ActionCable.server.broadcast(stream_name, { event: data["event"] })
  end

  private

  def stream_name
    "composition_room_#{params[:room_id]}"
  end
end
