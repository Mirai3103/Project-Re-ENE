import React, { useCallback, useRef, useState } from "react";
import * as PIXI from "pixi.js";
import { useLocation } from "wouter";
import type { InternalModel, Live2DModel } from "@laffy1309/pixi-live2d-lipsyncpatch/cubism4";
import { Live2DCanvas } from "@/components/HomePage/Live2DCanvas";
import { ChatPanel } from "@/components/HomePage/ChatPanel";
import { useLive2DAudio } from "@/hooks/useLive2DAudio";
import { useVoiceRecording } from "@/hooks/useVoiceRecording";
import { InvokeWithText } from "@wailsbindings/services/appservice";
import { GetChatHistory } from "@wailsbindings/services/chatservice";
import type { ChatMessage } from "@/types/chat";
import { useQuery } from "@/lib/query";

window.PIXI = PIXI;

declare global {
  interface Window {
    PIXI: typeof PIXI;
    playListUrl: (urls: string[]) => Promise<void>;
  }
}

/**
 * HomePage Component
 * Main page with Live2D character and chat interface
 */
type MessageContent = {
   content: {
    text: string;
    // ToolResponse *ToolResponse  `json:"toolResponse,omitempty"` // valid for kind==partToolResponse
    toolResponse?: {
      name: string;
      output: any;
    };
   }[];
   role: string;
}
function base64ToUtf8(base64: string) {
  const binary = atob(base64);
  const bytes = Uint8Array.from(binary, c => c.charCodeAt(0));
  return new TextDecoder("utf-8").decode(bytes);
}
export default function HomePage() {
  const modelRef = useRef<Live2DModel<InternalModel> | null>(null);
  const [speakingText, setSpeakingText] = useState<string>("");
  const [inputMessage, setInputMessage] = useState("");
  const [conversationID] = useState<string>(crypto.randomUUID());
  const [, navigate] = useLocation();
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [streamingMessage, setStreamingMessage] = useState<string|null>(null);
  const {data: chatHistory,refetch} = useQuery({
    queryKey: ["messages", conversationID],
    queryFn: () => GetChatHistory(conversationID),
  })
  React.useEffect(() => {
    if(chatHistory){
        console.log(chatHistory);
       const newMessages = chatHistory.map((item)=>{
        const content = JSON.parse(base64ToUtf8(item!.Content)) as MessageContent;
        console.log(content);
        if (content.role == 'tool'){
          const name = content.content.find((c)=>c.toolResponse?.name)?.toolResponse?.name
          const res = content.content.find((c)=>c.toolResponse?.output)
          console.log(res)
          return  {
            id: crypto.randomUUID(),
            role: 'tool',
            text: 'Gọi công cụ '+name,
            toolName: name,
            toolOutput: res?.toolResponse?.output,
            timestamp: new Date(item!.CreatedAt),
          }
        }
        return {
          id: crypto.randomUUID(),
          role: item!.Role,
          text: content.content.map(item => item.text).join(" "),
          timestamp: new Date(item!.CreatedAt),
        }
       });
       setMessages(newMessages as any);
    }
  }, [chatHistory]);
  // Custom hooks
  const { isRecording, startRecording, stopRecording } = useVoiceRecording(conversationID);
  const onSpeakingTextChange = useCallback((text: string, isDone?: boolean) => {
    if(isDone){
      console.log("isDone");
      setTimeout(() => {
        refetch();
        setStreamingMessage(null);
      }, 1000);
    }
    setSpeakingText(text);
    setStreamingMessage(prev => prev ? prev + text : text);
    
    
  }, [refetch]);
  useLive2DAudio(modelRef, onSpeakingTextChange);

  const handleModelReady = (model: Live2DModel<InternalModel>) => {
    modelRef.current = model;
  };

  const handleSendMessage = () => {
    if (!inputMessage.trim()) return;
    
    // TODO: Implement actual message sending logic
    InvokeWithText(conversationID, inputMessage);
    setInputMessage("");
    setMessages([...messages, { id: crypto.randomUUID(), role: "user", text: inputMessage, timestamp: new Date() }]);
  };

  const handleSettingsClick = () => {
    navigate("/settings");
  };

  return (
    <div className="w-full h-full flex bg-transparent relative overflow-hidden min-h-screen">
      <Live2DCanvas
        speakingText={speakingText}
        isRecording={isRecording}
        onStartRecording={startRecording}
        onStopRecording={stopRecording}
        onSettingsClick={handleSettingsClick}
        onModelReady={handleModelReady}
      />

      <ChatPanel
        messages={messages}
        inputValue={inputMessage}
        onInputChange={setInputMessage}
        onSendMessage={handleSendMessage}
        streamingMessage={streamingMessage}
      />
    </div>
  );
}
