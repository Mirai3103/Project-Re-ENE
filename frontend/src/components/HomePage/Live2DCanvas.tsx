import { useRef, useEffect } from "react";
import { Application, Ticker } from "pixi.js";
import {
  InternalModel,
  Live2DModel,
} from "@laffy1309/pixi-live2d-lipsyncpatch/cubism4";
import { DecorativeParticles } from "./DecorativeParticles";
import { SettingsButton } from "./SettingsButton";
import { MicButton } from "./MicButton";
import { SpeechBubble } from "./SpeechBubble";

interface Live2DCanvasProps {
  speakingText: string;
  isRecording: boolean;
  onStartRecording: () => void;
  onStopRecording: () => void;
  onSettingsClick: () => void;
  onModelReady: (model: Live2DModel<InternalModel>) => void;
}

export function Live2DCanvas({
  speakingText,
  isRecording,
  onStartRecording,
  onStopRecording,
  onSettingsClick,
  onModelReady,
}: Live2DCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const modelRef = useRef<Live2DModel<InternalModel> | null>(null);

  const resizeModel = () => {
    if (!modelRef.current) return;

    const model = modelRef.current;
    const w = window.innerWidth;
    const h = window.innerHeight;

    const scale = Math.min(w / 2000, h / 2000) * 1.2;

    model.scale.set(scale);
    model.x = w / 2;
    model.y = h;
    model.anchor.set(0.5, 0.4);
    model.motion("idle", 1);
  };

  useEffect(() => {
    let app: Application | null = null;
    let mounted = true;

    async function init() {
      // Wait for next frame to ensure canvas is fully rendered
      await new Promise(resolve => requestAnimationFrame(resolve));
      
      if (!canvasRef.current || !mounted) return;

      // Additional check to ensure canvas has dimensions
      const canvas = canvasRef.current;
      if (canvas.width === 0 || canvas.height === 0) {
        console.warn("Canvas not ready, retrying...");
        // Retry after a short delay
        setTimeout(() => {
          if (mounted) init();
        }, 100);
        return;
      }

      try {
        // Create Pixi Application (v6 synchronous pattern)
        app = new Application({
          view: canvas,
          autoStart: true,
          resizeTo: window,
          backgroundAlpha: 0,
          // Explicitly set renderer options to avoid shader issues
          antialias: true,
          resolution: window.devicePixelRatio || 1,
          autoDensity: true,
        });

        if (!mounted || !app || !app.stage) return;

        // Load the Live2D model
        const model = await Live2DModel.from(
          "/models/hiyori/hiyori_pro_t11.model3.json",
          {
            ticker: Ticker.shared,
          }
        );

        if (!mounted || !app || !app.stage) return;

        modelRef.current = model;
        app.stage.addChild(model);

        model.x = 0;
        model.y = 0;
        model.scale.set(0.35, 0.35);
        
        resizeModel();
        onModelReady(model);

        window.playListUrl = async (urls: string[]) => {
          for await (const url of urls) {
            if (!mounted) break;
            await new Promise((resolve) => {
              modelRef.current?.speak(url, {
                volume: 1,
                onFinish() {
                  resolve(true);
                },
              });
            });
          }
        };
      } catch (error) {
        console.error("Failed to initialize Live2D:", error);
      }
    }

    init();
    window.addEventListener("resize", resizeModel);

    return () => {
      mounted = false;
      window.removeEventListener("resize", resizeModel);
      if (app) {
        try {
          app.destroy(true, { children: true, texture: true, baseTexture: true });
        } catch (error) {
          console.error("Error destroying app:", error);
        }
      }
    };
  }, [onModelReady]);

  return (
    <div className="flex-1 relative overflow-hidden rounded-lg">
      <canvas ref={canvasRef} className="w-full h-full"></canvas>
      
      <DecorativeParticles />
      <SettingsButton onClick={onSettingsClick} />

      {/* Bottom UI Container */}
      <div className="absolute bottom-0 left-0 right-0 flex flex-col items-center pb-1 gap-4">
        <SpeechBubble text={speakingText} />
        <MicButton
          isRecording={isRecording}
          onStart={onStartRecording}
          onStop={onStopRecording}
        />
      </div>
    </div>
  );
}

