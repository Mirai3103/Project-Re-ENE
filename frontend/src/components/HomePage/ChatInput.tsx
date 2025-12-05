import { Send } from "lucide-react";
import { Input } from "@/components/ui/input";

interface ChatInputProps {
  value: string;
  onChange: (value: string) => void;
  onSubmit: () => void;
}

export function ChatInput({ value, onChange, onSubmit }: ChatInputProps) {
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!value.trim()) return;
    onSubmit();
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <Input
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Type a message..."
        className="flex-1 bg-card/70 border-border/50 focus-visible:ring-primary/50 rounded-xl"
      />
      <button
        type="submit"
        disabled={!value.trim()}
        className="
          shrink-0 w-10 h-10 rounded-xl
          bg-linear-to-br from-primary/80 to-primary/60
          hover:from-primary hover:to-primary/80
          disabled:from-muted disabled:to-muted
          disabled:cursor-not-allowed
          flex items-center justify-center
          transition-all duration-200
          hover:scale-105 active:scale-95
          disabled:hover:scale-100
        "
      >
        <Send className={`w-4 h-4 ${value.trim() ? 'text-white' : 'text-muted-foreground'}`} />
      </button>
    </form>
  );
}

