export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  text: string;
  timestamp: Date;
}

export interface Live2DModelRef {
  speak: (url: string, options: {
    volume: number;
    onFinish?: () => void;
    onError?: () => void;
  }) => void;
  motion: (group: string, index: number) => void;
  scale: {
    set: (scale: number) => void;
  };
  x: number;
  y: number;
  anchor: {
    set: (x: number, y: number) => void;
  };
}

