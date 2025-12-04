import { Route, Router, Switch } from "wouter";
import HomePage from "./pages/";
import { ThemeProvider } from "./components/theme-provider";
import CharacterPage from "./pages/Settings/Character";
import ModelPage from "./pages/Settings/Model";
import AsrPage from "./pages/Settings/Asr";
import ElevenLabsAsrPage from "./pages/Settings/Asr/elevenlabs";

import { useHashLocation } from "wouter/use-hash-location";
import SettingPage from "./pages/Settings";
import React from "react";
import { Window } from "@wailsio/runtime";
import { AuroraBackground } from "./components/ui/shadcn-io/aurora-background";
function App() {
  const location = useHashLocation();
  React.useEffect(() => {
    // Set initial window size based on route
    async function adjustWindow() {
      if (location[0].startsWith("/settings")) {
        await Window.Maximise();
      } else {
        // Check if window is currently maximized or large enough
        const size = await Window.Size();
        // If window is smaller than the breakpoint for showing chat (1280px), set to compact size
        // Otherwise, keep the current size for the chat panel layout
        if (size.width < 1280) {
          await Window.SetSize(400, 400);
          await Window.Center();
        }
        // If already wide enough, let user keep their preferred size
      }
    }
    adjustWindow();
  }, [location]);
  return (
    <ThemeProvider>
      <Router hook={useHashLocation} base="">
        <Switch>
          <Route path="/" component={HomePage} />
          <Route path="/settings" component={SettingPage} />

          <Route path="/settings/character" component={CharacterPage} />
          <Route path="/settings/model" component={ModelPage} />
          <Route path="/settings/asr" component={AsrPage} />

          {/* ASR Provider Detail Routes */}
          <Route
            path="/settings/asr/providers/elevenlabs"
            component={ElevenLabsAsrPage}
          />
          <Route path="/settings/asr/providers/:providerId">
            {(params) => (
              <AuroraBackground showRadialGradient className="overflow-auto">
                <div className="min-h-screen w-full max-w-4xl mx-auto p-4 py-10 flex items-center justify-center">
                  <div className="text-center">
                    <h1 className="text-4xl font-bold text-primary mb-4">
                      {params.providerId?.toUpperCase()} Configuration
                    </h1>
                    <p className="text-muted-foreground mb-6">
                      Provider-specific settings coming soon...
                    </p>
                    <a
                      href="#/settings/asr"
                      className="inline-block px-6 py-3 bg-primary text-primary-foreground rounded-xl hover:bg-primary/90 transition-colors"
                    >
                      Back to ASR Settings
                    </a>
                  </div>
                </div>
              </AuroraBackground>
            )}
          </Route>

          {/* Other Settings Routes */}
          <Route path="/settings/:section">
            {(params) => (
              <AuroraBackground showRadialGradient className="overflow-auto">
                <div className="min-h-screen w-full max-w-4xl mx-auto p-4 py-10 flex items-center justify-center">
                  <div className="text-center">
                    <h1 className="text-4xl font-bold text-primary mb-4">
                      {params.section?.toUpperCase()} Settings
                    </h1>
                    <p className="text-muted-foreground mb-6">Coming soon...</p>
                    <a
                      href="#/settings"
                      className="inline-block px-6 py-3 bg-primary text-primary-foreground rounded-xl hover:bg-primary/90 transition-colors"
                    >
                      Back to Settings
                    </a>
                  </div>
                </div>
              </AuroraBackground>
            )}
          </Route>

          <Route path="/users/:name">
            {(params) => <>Hello, {params.name}!</>}
          </Route>

          <Route>404: No such page!</Route>
        </Switch>
      </Router>
    </ThemeProvider>
  );
}

export default App;
