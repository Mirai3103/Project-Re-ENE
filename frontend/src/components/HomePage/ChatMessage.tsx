import { Clock } from "lucide-react";
import type { ChatMessage as ChatMessageType } from "@/types/chat";
import { formatTime } from "@/utils/time";

interface ChatMessageProps {
  message: ChatMessageType;
}

export function ChatMessage({ message }: ChatMessageProps) {
  return (
    <div
      className={`flex animate-in fade-in slide-in-from-bottom-2 ${
        message.role === "user" ? "justify-end" : "justify-start"
      }`}
    >
      {/* Message Bubble */}
      <div
        className={`group relative max-w-[85%] ${
          message.role === "user"
            ? "bg-linear-to-br from-primary/30 to-primary/20 border border-primary/40"
            : "bg-card/70 border border-border/50"
        } rounded-2xl px-4 py-2.5 backdrop-blur-sm transition-all hover:scale-[1.01]`}
      >
        <p className="text-sm text-foreground leading-relaxed wrap-break-word">
          {message.text}
        </p>
        
        {/* Timestamp tooltip */}
        <div
          className={`absolute -bottom-6 ${
            message.role === "user" ? "right-0" : "left-0"
          } opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none`}
        >
          <div className="flex items-center gap-1 bg-muted/90 text-muted-foreground text-xs px-2 py-0.5 rounded-full backdrop-blur-sm">
            <Clock className="w-3 h-3" />
            {formatTime(message.timestamp)}
          </div>
        </div>
      </div>
    </div>
  );
}

