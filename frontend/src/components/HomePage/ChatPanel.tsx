import { Badge } from "@/components/ui/badge";
import { ChatMessage } from "./ChatMessage";
import { ChatInput } from "./ChatInput";
import type { ChatMessage as ChatMessageType } from "@/types/chat";

interface ChatPanelProps {
  messages: ChatMessageType[];
  inputValue: string;
  onInputChange: (value: string) => void;
  onSendMessage: () => void;
  streamingMessage: string|null;
}

export function ChatPanel({
  messages,
  inputValue,
  onInputChange,
  onSendMessage,
  streamingMessage,
}: ChatPanelProps) {
  return (
    <div className="hidden xl:flex flex-col w-96 rounded-2xl  bg-transparent py-10 h-full px-2 h-screen  ">
      {/* Chat Header */}
      <div className="border-l border-border/50 flex-1 flex flex-col bg-card/30 backdrop-blur-md h-full rounded-2xl">
      <div className="p-4 border-b border-border/50 bg-card/50 backdrop-blur-sm rounded-2xl">
        <div className="flex items-center justify-between">
          <div className="flex-1">
            <h3 className="font-semibold text-foreground">Chat History</h3>
            <p className="text-xs text-muted-foreground">
              {messages.length} messages
            </p>
          </div>
          <Badge variant="secondary" className="bg-primary/10 text-primary">
            Live
          </Badge>
        </div>
      </div>

      {/* Chat Messages */}
      <div className="flex-1 overflow-y-auto p-4 space-y-3 grow ">
        {messages.map((message) => (
          <ChatMessage key={message.id} message={message} />
        ))}
        {streamingMessage && (
          <ChatMessage message={{ id: crypto.randomUUID(), role: "assistant", text: streamingMessage, timestamp: new Date() }} />
        )}
      </div>

      {/* Chat Input */}
      <div className="p-4 border-t border-border/50 bg-card/50 backdrop-blur-sm rounded-2xl">
        <ChatInput
          value={inputValue}
          onChange={onInputChange}
          onSubmit={onSendMessage}
        />
      </div>
        </div>
    </div>
  );
}

