import { Settings } from "lucide-react";

interface SettingsButtonProps {
  onClick: () => void;
}

export function SettingsButton({ onClick }: SettingsButtonProps) {
  return (
    <button
      onClick={onClick}
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
  );
}

