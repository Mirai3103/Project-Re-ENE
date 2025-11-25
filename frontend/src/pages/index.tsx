import {
  InternalModel,
  Live2DModel,
} from "@laffy1309/pixi-live2d-lipsyncpatch/cubism4";
import { Application, Ticker } from "pixi.js";
import React, { useRef, useState } from "react";
import * as PIXI from "pixi.js";
import { Events, Screens, Window } from "@wailsio/runtime";
import {
  StartRecording,
  StopRecording,
} from "@wailsbindings/services/recorderservice";
import {
PlayAudioData,
} from "@wailsbindings/services";
import { InvokeWithAudio } from "@wailsbindings/services/appservice";
import PQueue from "p-queue";
import { Mic, Settings } from "lucide-react";
import { useLocation } from "wouter";
window.PIXI = PIXI;
declare global {
  interface Window {
    PIXI: typeof PIXI;
    playListUrl: (urls: string[]) => Promise<void>;
  }
}

const speakQueue = new PQueue({ concurrency: 1 });

function base64ToBlobUrl(base64: string): string {
  // Nếu base64 có dạng dataURI → tách header + mime
  let mime = "audio/wav";
  let pureBase64 = base64;

  if (base64.startsWith("data:")) {
    const split = base64.split(",");
    const header = split[0];
    pureBase64 = split[1];
    mime = header.match(/data:(.*?);base64/)?.[1] || mime;
  }

  const byteCharacters = atob(pureBase64);
  const byteNumbers = new Array(byteCharacters.length);

  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }

  const byteArray = new Uint8Array(byteNumbers);
  const blob = new Blob([byteArray], { type: mime });

  return URL.createObjectURL(blob);
}
type SpeakData = {
  text: string;
  base64: string;
};
export default function HomePage() {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const modelRef = useRef<Live2DModel<InternalModel>>(null);
  const [speakingText, setSpeakingText] = useState<string>("");
  const [conversationID] = useState<string>(crypto.randomUUID());
  function resizeModel() {
    if (!modelRef.current) return;

    const model = modelRef.current;
    const w = window.innerWidth;
    const h = window.innerHeight;

    // scale dựa vào chiều cao hoặc chiều rộng (tùy bạn muốn ưu tiên)
    const scale = Math.min(w / 2000, h / 2000) * 1.2; // 1.2 để model to hơn một chút

    model.scale.set(scale);

    // luôn giữ model ở giữa màn hình
    model.x = w / 2;
    model.y = h;
    model.anchor.set(0.5, 0.4);
    model.motion("idle", 1);
  }
  React.useEffect(() => {
    async function init() {
      const app = new Application({
        view: canvasRef.current!,
        autoStart: true,
        resizeTo: window,
        backgroundAlpha: 0,
      });

      const model = await Live2DModel.from(
        "/models/hiyori/hiyori_pro_t11.model3.json",
        {
          ticker: Ticker.shared,
        }
      );
      modelRef.current = model;
      app.stage.addChild(model);

      model.x = 0;
      model.y = 0;
      model.scale.set(0.35, 0.35);
      console.log(modelRef.current.internalModel);
      resizeModel();
      window.playListUrl = async (urls: string[]) => {
        for await (const url of urls) {
          await new Promise((resolve) => {
            console.log("play-audio", url);
            modelRef.current?.speak(url, {
              volume: 1,
              onFinish() {
                resolve(true);
              },
            });
          });
        }
      };
    }

    init();
    window.addEventListener("resize", resizeModel);

    return () => {
      window.removeEventListener("resize", resizeModel);
    };
  }, []);

  // Chạy motion
  function playMotion(group: string, index: number) {
    console.log("Play motion:", group, index);
    modelRef.current?.motion(group, index);
  }

  React.useEffect(() => {
    Events.On("live2d:play-audio", ({ data }) => {
      playMotion(data.group, data.index);
    });
    Events.On("live2d:play-audio", ({ data }: { data: PlayAudioData }) => {
      speakQueue.add(async () => {
        console.log("play-audio", data.Text);
        const url = base64ToBlobUrl(data.Base64);
        setSpeakingText(data.Text);

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
        setSpeakingText("");
      });
    });

    return () => {
      Events.Off("live2d:set-motion");
      Events.Off("live2d:play-audio");
    };
  }, []);
  const [isRecording, setIsRecording] = useState(false);
  function startRecording() {
    StartRecording();
    setIsRecording(true);
  }
  async function stopRecording() {
    const audioPath = await StopRecording();
    setIsRecording(false);
    await InvokeWithAudio(conversationID, audioPath);
  
  }
  const [, navigate] = useLocation();

  function openSettingsPage() {
    navigate("/settings");
  }
  return (
    <div className="w-full h-full bg-transparent relative overflow-hidden rounded-lg">
      <canvas ref={canvasRef} className="w-full h-full"></canvas>
      <div className="absolute inset-0 pointer-events-none overflow-hidden">
        <div className="absolute top-20 left-10 w-3 h-3 bg-pink-300 rounded-full animate-pulse opacity-60"></div>
        <div
          className="absolute top-40 right-20 w-2 h-2 bg-purple-300 rounded-full animate-pulse opacity-50"
          style={{ animationDelay: "0.5s" }}
        ></div>
        <div
          className="absolute bottom-32 left-16 w-2 h-2 bg-blue-300 rounded-full animate-pulse opacity-40"
          style={{ animationDelay: "1s" }}
        ></div>
        <div
          className="absolute top-60 right-40 w-3 h-3 bg-pink-200 rounded-full animate-pulse opacity-50"
          style={{ animationDelay: "1.5s" }}
        ></div>
        <div
          className="absolute bottom-60 right-10 w-2 h-2 bg-purple-200 rounded-full animate-pulse opacity-60"
          style={{ animationDelay: "0.8s" }}
        ></div>
      </div>

      {/* Settings Button */}
      <button
        onClick={openSettingsPage}
        className="
          absolute top-6 right-6
          group
          w-12 h-12 rounded-full
          bg-linear-to-br from-pink-300 via-purple-300 to-blue-300
          hover:from-pink-400 hover:via-purple-400 hover:to-blue-400
          flex items-center justify-center
          transition-all duration-300 ease-out
          transform hover:scale-110 hover:rotate-90 active:scale-95
          shadow-lg hover:shadow-pink-300/50
          backdrop-blur-sm
        "
      >
        <Settings
          className="w-5 h-5 text-white drop-shadow-lg"
          strokeWidth={2.5}
        />

        {/* Tooltip */}
        <div className="absolute -bottom-12 right-0 opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none">
          <div className="bg-gray-800/90 text-white text-sm px-3 py-1.5 rounded-full whitespace-nowrap backdrop-blur-sm">
            ⚙️ Settings
          </div>
        </div>
      </button>

      {/* Bottom UI Container */}
      <div className="absolute bottom-0 left-0 right-0 flex flex-col items-center pb-1 gap-4">
        {/* Speaking Text Bubble */}
        {speakingText && (
          <div className="animate-in fade-in slide-in-from-bottom-4 duration-300 max-w-2xl mx-4">
            <div className="relative">
              {/* Cute speech bubble */}
              <div className="bg-linear-to-br from-white to-pink-50 backdrop-blur-md rounded-3xl px-4 py-2 shadow-2xl border-2 border-pink-200/50">
                <p className="text-gray-800 text-base leading-relaxed font-medium text-center">
                  {speakingText}
                </p>
                {/* Small decorative hearts */}
                <div className="absolute -top-2 -right-2 text-pink-400 text-xl animate-bounce">
                  ♡
                </div>
              </div>
              {/* Speech bubble tail */}
              <div className="absolute -bottom-3 left-1/2 transform -translate-x-1/2 w-6 h-6 bg-linear-to-br from-white to-pink-50 rotate-45 border-r-2 border-b-2 border-pink-200/50"></div>
            </div>
          </div>
        )}

        {/* Microphone Button */}
        <button
          onClick={isRecording ? stopRecording : startRecording}
          className={`
            group relative
            w-16 h-16 rounded-full
            flex items-center justify-center
            transition-all duration-300 ease-out
            transform hover:scale-110 active:scale-95
            shadow-2xl hover:shadow-pink-300/50
            ${
              isRecording
                ? "bg-linear-to-br from-red-400 to-pink-500 animate-pulse"
                : "bg-linear-to-br from-pink-300 via-purple-300 to-blue-300 hover:from-pink-400 hover:via-purple-400 hover:to-blue-400"
            }
          `}
        >
          {/* Glow effect when recording */}
          {isRecording && (
            <div className="absolute inset-0 rounded-full bg-red-400 animate-ping opacity-75"></div>
          )}

          {/* Button content */}
          <div className="relative z-10 flex items-center justify-center w-full h-full">
            <Mic
              className={`w-6 h-6 transition-colors ${
                isRecording ? "text-white" : "text-white drop-shadow-lg"
              }`}
              strokeWidth={2.5}
            />
          </div>

          {/* Cute ring decoration */}
          <div
            className={`
            absolute inset-0 rounded-full border-4 border-white/40
            transition-all duration-300
            ${isRecording ? "scale-110 opacity-0" : "scale-100 opacity-100"}
          `}
          ></div>

          {/* Tooltip */}
          <div className="absolute -top-12 left-1/2 transform -translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none">
            <div className="bg-gray-800/90 text-white text-sm px-3 py-1.5 rounded-full whitespace-nowrap backdrop-blur-sm">
              {isRecording ? "✨ Recording..." : "🎤 Click to speak"}
            </div>
          </div>
        </button>
      </div>
    </div>
  );
}
