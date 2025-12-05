import type { ChatMessage } from "@/types/chat";

export const mockChatHistory: ChatMessage[] = [
  {
    id: "1",
    role: "user",
    text: "こんにちは！元気ですか？",
    timestamp: new Date(Date.now() - 1000 * 60 * 30),
  },
  {
    id: "2",
    role: "assistant",
    text: "こんにちは！元気です～ 今日も頑張りましょう！✨",
    timestamp: new Date(Date.now() - 1000 * 60 * 29),
  },
  {
    id: "3",
    role: "user",
    text: "今日の予定を教えて",
    timestamp: new Date(Date.now() - 1000 * 60 * 25),
  },
  {
    id: "4",
    role: "assistant",
    text: "はい！今日の予定ですね。午前中はミーティングがあります。午後は自由時間ですよ～ 何か特別なことをしたいですか？",
    timestamp: new Date(Date.now() - 1000 * 60 * 24),
  },
  {
    id: "5",
    role: "user",
    text: "ありがとう！助かるよ",
    timestamp: new Date(Date.now() - 1000 * 60 * 20),
  },
  {
    id: "6",
    role: "assistant",
    text: "どういたしまして！いつでも聞いてくださいね～ 💖",
    timestamp: new Date(Date.now() - 1000 * 60 * 19),
  },
  {
    id: "7",
    role: "user",
    text: "Can you speak English?",
    timestamp: new Date(Date.now() - 1000 * 60 * 15),
  },
  {
    id: "8",
    role: "assistant",
    text: "Of course! I can speak multiple languages. How can I help you today? 😊",
    timestamp: new Date(Date.now() - 1000 * 60 * 14),
  },
];

