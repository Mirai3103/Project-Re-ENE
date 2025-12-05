import { Mic } from "lucide-react";

interface MicButtonProps {
  isRecording: boolean;
  onStart: () => void;
  onStop: () => void;
}

export function MicButton({ isRecording, onStart, onStop }: MicButtonProps) {
  return (
    <button
      onClick={isRecording ? onStop : onStart}
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
  );
}

