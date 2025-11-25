import React, { useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { RippleButton } from "@/components/ui/shadcn-io/ripple-button";
import { AuroraBackground } from "@/components/ui/shadcn-io/aurora-background";
import { useLocation } from "wouter";
import {
  ArrowLeft,
  Mic,
  Settings2,
  CheckCircle,
  Zap,
  Cloud,
  Sparkles,
  Radio,
} from "lucide-react";
import { useQuery } from "@/lib/query";
import { RecorderService, ConfigService } from "@wailsbindings/services";
// Mock data for providers
const asrProviders = [
  {
    id: "elevenlabs",
    name: "ElevenLabs",
    description: "High-quality speech recognition with multilingual support",
    icon: "🎙️",
    status: "active",
    features: ["Multilingual", "Real-time", "High Accuracy"],
    color: "from-purple-500 to-pink-500",
  },
  {
    id: "google",
    name: "Google Cloud",
    description: "Powerful speech-to-text with extensive language support",
    icon: "🌐",
    status: "available",
    features: ["120+ Languages", "Cloud-based", "Punctuation"],
    color: "from-blue-500 to-cyan-500",
  },
  {
    id: "azure",
    name: "Azure Speech",
    description: "Microsoft's enterprise-grade speech recognition service",
    icon: "☁️",
    status: "available",
    features: ["Custom Models", "Diarization", "Real-time"],
    color: "from-blue-600 to-indigo-600",
  },
  {
    id: "whisper",
    name: "OpenAI Whisper",
    description: "Open-source automatic speech recognition system",
    icon: "🤖",
    status: "available",
    features: ["Open Source", "Offline", "99 Languages"],
    color: "from-green-500 to-emerald-500",
  },
  {
    id: "assembly",
    name: "AssemblyAI",
    description: "AI-powered transcription and speech understanding",
    icon: "⚡",
    status: "available",
    features: ["Speaker Labels", "Sentiment", "Auto Chapters"],
    color: "from-orange-500 to-red-500",
  },
  {
    id: "deepgram",
    name: "Deepgram",
    description: "Fast and accurate speech recognition API",
    icon: "🔊",
    status: "available",
    features: ["Ultra Fast", "WebSocket", "Custom Vocab"],
    color: "from-violet-500 to-purple-500",
  },
];

// Mock input devices
const inputDevices = [
  { id: "default", name: "Default Microphone" },
  { id: "device-1", name: "Built-in Microphone" },
  { id: "device-2", name: "USB Microphone (Blue Yeti)" },
  { id: "device-3", name: "Wireless Headset" },
];

export default function AsrPage() {
  const [, navigate] = useLocation();
  const [selectedProvider, setSelectedProvider] = useState("elevenlabs");
  const [selectedDevice, setSelectedDevice] = useState("default");
  const {
    data: availableInputDevices,
    error: availableInputDevicesError,
    isLoading: availableInputDevicesIsLoading,
  } = useQuery({
    queryKey: ["available-input-devices"],
    queryFn: () => RecorderService.GetAvailableInputDevices(),
  });
  const {
    data: config,
    error: configError,
    isLoading: configIsLoading,
  } = useQuery({
    queryKey: ["configs"],
    queryFn: () => ConfigService.GetConfig(),
  });
  const defaultInputDevice = React.useMemo(() => {
    if (configIsLoading || availableInputDevicesIsLoading) return null;
    const name = config?.ASRConfig.InputDevice;
    const defaultName = availableInputDevices?.find(
      (device) => device.IsDefault
    )?.Name;
    if (!name) return defaultName;
    return (
      availableInputDevices?.find((device) => device.Name === name)?.Name ||
      defaultName
    );
  }, [availableInputDevices, config]);
  console.log(defaultInputDevice);
  React.useEffect(() => {
    if (defaultInputDevice) {
      setSelectedDevice(defaultInputDevice);
    }
  }, [defaultInputDevice]);

  return (
    <AuroraBackground showRadialGradient className="overflow-auto">
      <div className="min-h-screen w-full max-w-[1600px] p-4 py-10 font-sans selection:bg-primary/20">
        {/* Header */}
        <div className="relative space-y-2 pb-8 flex items-center justify-between">
          <RippleButton
            onClick={() => {
              navigate("/settings");
            }}
            variant="outline"
            className="bg-primary/10 text-primary hover:bg-primary hover:text-primary-foreground"
          >
            <ArrowLeft className="h-5 w-5 mr-2" />
            Back
          </RippleButton>

          <div className="text-center flex-1">
            <div className="flex items-center justify-center gap-3 mb-2">
              <Mic className="h-8 w-8 text-primary" />
              <h1 className="text-4xl font-bold tracking-tight text-primary drop-shadow-sm">
                ASR Configuration
              </h1>
            </div>
            <p className="text-muted-foreground text-lg">
              Configure Automatic Speech Recognition settings
            </p>
          </div>

          <div className="w-[100px]"></div>
        </div>

        {/* Configuration Section */}
        <div className="mx-auto max-w-4xl space-y-6 mb-12">
          <Card className="border-2 border-primary/20 bg-card/50 backdrop-blur-sm relative overflow-hidden">
            <div className="absolute -right-12 -top-12 -z-10 h-32 w-32 rounded-full bg-primary/20 blur-3xl" />

            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-2xl">
                <Settings2 className="h-6 w-6 text-primary" />
                Active Configuration
              </CardTitle>
              <CardDescription>
                Select your speech recognition provider and input device
              </CardDescription>
            </CardHeader>

            <CardContent className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Provider Selection */}
                <div className="space-y-3">
                  <label className="text-sm font-semibold text-foreground flex items-center gap-2">
                    <Sparkles className="h-4 w-4 text-primary" />
                    Speech Provider
                  </label>
                  <Select
                    value={selectedProvider}
                    onValueChange={setSelectedProvider}
                  >
                    <SelectTrigger className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 transition-colors h-12 w-full">
                      <SelectValue placeholder="Select provider" />
                    </SelectTrigger>
                    <SelectContent>
                      {asrProviders.map((provider) => (
                        <SelectItem key={provider.id} value={provider.id}>
                          <div className="flex items-center gap-2">
                            <div className="font-medium">{provider.name}</div>
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <p className="text-xs text-muted-foreground">
                    Choose the ASR service for speech recognition
                  </p>
                </div>

                {/* Input Device Selection */}
                <div className="space-y-3">
                  <label className="text-sm font-semibold text-foreground flex items-center gap-2">
                    <Radio className="h-4 w-4 text-primary" />
                    Input Device
                  </label>
                  <Select
                    value={selectedDevice}
                    onValueChange={setSelectedDevice}
                  >
                    <SelectTrigger className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 transition-colors h-12 w-full">
                      <SelectValue placeholder="Select device" />
                    </SelectTrigger>
                    <SelectContent>
                      {availableInputDevices?.map((device) => (
                        <SelectItem key={device.Name} value={device.Name}>
                          {device.Name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <p className="text-xs text-muted-foreground">
                    Select the microphone for audio input
                  </p>
                </div>
              </div>

              {/* Save Button */}
              <div className="flex gap-3 pt-4 border-t border-border/50">
                <RippleButton
                  onClick={() => {
                    console.log("Save config:", {
                      selectedProvider,
                      selectedDevice,
                    });
                    alert("Configuration saved! (Mock implementation)");
                  }}
                  className="flex-1 md:flex-initial bg-primary text-primary-foreground hover:bg-primary/90"
                >
                  <CheckCircle className="h-5 w-5 mr-2" />
                  Save Configuration
                </RippleButton>
                <RippleButton
                  variant="outline"
                  onClick={() => {
                    setSelectedProvider("elevenlabs");
                    setSelectedDevice("default");
                  }}
                  className="bg-secondary/50 hover:bg-secondary"
                >
                  Reset
                </RippleButton>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Available Providers Section */}
        <div className="mx-auto max-w-6xl space-y-6">
          <div className="flex items-center gap-3">
            <Cloud className="h-6 w-6 text-primary" />
            <h2 className="text-2xl font-bold text-foreground">
              Available Providers
            </h2>
            <Badge variant="secondary" className="bg-primary/10 text-primary">
              {asrProviders.length} Services
            </Badge>
          </div>

          {/* Providers Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {asrProviders.map((provider) => (
              <Card
                key={provider.id}
                onClick={() => {
                  navigate(`/settings/asr/providers/${provider.id}`);
                }}
                className="group relative overflow-hidden cursor-pointer border-2 border-transparent bg-card/50 backdrop-blur-sm transition-all hover:border-primary/50 hover:shadow-lg hover:shadow-primary/20"
              >
                <div className="absolute -left-12 -top-12 -z-10 h-24 w-24 rounded-full bg-primary/10 blur-3xl transition-all duration-1000 ease-out pointer-events-none group-hover:bg-primary/30 group-hover:h-48 group-hover:w-48" />

                <CardHeader className="space-y-3">
                  <div className="flex items-start justify-between">
                    <div
                      className={`w-14 h-14 rounded-2xl bg-linear-to-br ${provider.color} flex items-center justify-center text-3xl shadow-lg`}
                    >
                      {provider.icon}
                    </div>
                    {provider.status === "active" && (
                      <Badge className="bg-primary/20 text-primary border-primary/30">
                        <CheckCircle className="h-3 w-3 mr-1" />
                        Active
                      </Badge>
                    )}
                  </div>

                  <div>
                    <CardTitle className="text-xl mb-1">
                      {provider.name}
                    </CardTitle>
                    <CardDescription className="line-clamp-2">
                      {provider.description}
                    </CardDescription>
                  </div>
                </CardHeader>

                <CardContent className="space-y-4">
                  {/* Features */}
                  <div className="flex flex-wrap gap-2">
                    {provider.features.map((feature) => (
                      <Badge
                        key={feature}
                        variant="secondary"
                        className="bg-secondary/50 text-xs font-normal"
                      >
                        {feature}
                      </Badge>
                    ))}
                  </div>

                  {/* Configure Button */}
                  <RippleButton
                    onClick={(e) => {
                      e.stopPropagation();
                      navigate(`/settings/asr/providers/${provider.id}`);
                    }}
                    className="w-full bg-primary/10 text-primary hover:bg-primary hover:text-primary-foreground"
                  >
                    <Settings2 className="h-4 w-4 mr-2" />
                    Configure
                  </RippleButton>
                </CardContent>

                {/* Hover Indicator */}
                <div className="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity">
                  <Zap className="h-5 w-5 text-primary animate-pulse" />
                </div>
              </Card>
            ))}
          </div>
        </div>
      </div>
    </AuroraBackground>
  );
}
