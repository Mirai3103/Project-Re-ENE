interface SpeechBubbleProps {
  text: string;
}

export function SpeechBubble({ text }: SpeechBubbleProps) {
  if (!text) return null;

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-300 max-w-2xl mx-4">
      <div className="relative">
        {/* Cute speech bubble */}
        <div className="bg-linear-to-br from-white to-pink-50 backdrop-blur-md rounded-3xl px-4 py-2 shadow-2xl border-2 border-pink-200/50">
          <p className="text-gray-800 text-base leading-relaxed font-medium text-center">
            {text}
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
  );
}

