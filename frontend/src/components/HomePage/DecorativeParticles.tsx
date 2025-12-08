export function DecorativeParticles() {
  return (
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
  );
}

