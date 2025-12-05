import { useState } from "react";
import {
  StartRecording,
  StopRecording,
} from "@wailsbindings/services/recorderservice";
import { InvokeWithAudio } from "@wailsbindings/services/appservice";

export function useVoiceRecording(conversationID: string) {
  const [isRecording, setIsRecording] = useState(false);

  const startRecording = () => {
    StartRecording();
    setIsRecording(true);
  };

  const stopRecording = async () => {
    const audioPath = await StopRecording();
    setIsRecording(false);
    await InvokeWithAudio(conversationID, audioPath);
  };

  return {
    isRecording,
    startRecording,
    stopRecording,
  };
}

