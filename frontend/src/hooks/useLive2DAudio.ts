import { useEffect } from "react";
import { Events } from "@wailsio/runtime";
import PQueue from "p-queue";
import type { PlayAudioData } from "@wailsbindings/services";
import type { Live2DModelRef } from "@/types/chat";
import { base64ToBlobUrl } from "@/utils/audio";

const speakQueue = new PQueue({ concurrency: 1 });

export function useLive2DAudio(
  modelRef: React.MutableRefObject<Live2DModelRef | null>,
  onSpeakingTextChange: (text: string, isDone?: boolean) => void
) {
  useEffect(() => {
    const unsubscribe = Events.On("live2d:play-audio", ({ data }: { data: PlayAudioData }) => {
     

      speakQueue.add(async () => {
        if (data.IsDone) {
          onSpeakingTextChange("", true);
          return;
        }
        console.log("play-audio", data.Text);
        const url = base64ToBlobUrl(data.Base64);
        onSpeakingTextChange(data.Text);

        await new Promise((resolve) => {
          modelRef.current?.speak(url, {
            volume: 1,
            onFinish() {
              resolve(true);
            },
            onError() {
              resolve(false);
            },
          });
        });

        await new Promise((resolve) => setTimeout(resolve, 500));
        console.log("play-audio finished");
        onSpeakingTextChange("");
      });
    });

    return () => {
      Events.Off("live2d:play-audio");
      if (unsubscribe) unsubscribe();
    };
  }, [modelRef, onSpeakingTextChange]);
}

