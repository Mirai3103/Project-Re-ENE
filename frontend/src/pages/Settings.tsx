import {
  Card,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { RippleButton } from "@/components/ui/shadcn-io/ripple-button";
import {
  Mic,
  User,
  Bot,
  FileText,
  Volume2,
  Brain,
  ArrowLeft,
  PersonStanding,
} from "lucide-react";
import { ModeToggle } from "@/components/mode-toggle";
import { AuroraBackground } from "@/components/ui/shadcn-io/aurora-background";
import { Link, useLocation } from "wouter";
import { Window, Screens } from "@wailsio/runtime";
// Mock Data based on comments
const configSections = [
  {
    id: "asr",
    title: "ASR Config",
    description: "Automatic Speech Recognition settings",
    icon: Mic,
    details: [
      { label: "Provider", value: "elevenlabs" },
      { label: "Model ID", value: "scribe_v1" },
      { label: "Language", value: "vi" },
    ],
    href: "/settings/asr",
  },
  {
    id: "live2d",
    title: "Live2D Models Config",
    description: "Live2D models settings",
    icon: PersonStanding,
    href: "/settings/model",
  },
  {
    id: "character",
    title: "Character Config",
    description: "Avatar and persona settings",
    href: "/settings/character",
    icon: User,
    details: [
      { label: "Name", value: "Ene" },
      { label: "Avatar", value: "ene.png" },
      { label: "Live2D", value: "ene.model" },
    ],
  },
  {
    id: "llm",
    title: "LLM Config",
    description: "Language Model settings",
    icon: Brain,
    details: [
      { label: "Provider", value: "gemini" },
      { label: "Model", value: "gemini-2.0-flash" },
      { label: "Temperature", value: "0.7" },
    ],
    href: "/settings/llm",
  },
  {
    id: "logger",
    title: "Logger Config",
    description: "System logging preferences",
    icon: FileText,
    details: [
      { label: "Level", value: "info" },
      { label: "Mode", value: "console" },
      { label: "Path", value: "logs/app.log" },
    ],
    href: "/settings/logger",
  },
  {
    id: "tts",
    title: "TTS Config",
    description: "Text-to-Speech settings",
    icon: Volume2,
    details: [
      { label: "Provider", value: "elevenlabs" },
      { label: "Model", value: "eleven_flash_v2_5" },
      { label: "Voice ID", value: "4zQx..." },
    ],
    href: "/settings/tts",
  },
  {
    id: "agent",
    title: "Agent Config",
    description: "Memory and behavior settings",
    icon: Bot,
    details: [
      { label: "Memory Window", value: "10" },
      { label: "Storage", value: ".data/conversations" },
    ],
    href: "/settings/agent",
  },
];

export default function SettingPage() {
  const [, navigate] = useLocation();
  return (
    <AuroraBackground showRadialGradient className="overflow-auto">
      <div className="min-h-screen w-full max-w-7xl p-4 py-10 font-sans selection:bg-primary/20">
        <div className="relative space-y-2 text-center pb-6 flex items-center justify-around">
          <RippleButton
            onClick={() => {
              navigate("/");
            }}
            variant="outline"
            className="bg-primary/10 text-primary hover:bg-primary hover:text-primary-foreground"
          >
            <ArrowLeft className="h-6 w-6" />
            <Link to="/">Back</Link>
          </RippleButton>
          <div>
            <h1 className="text-4xl font-bold tracking-tight text-primary drop-shadow-sm">
              Settings
            </h1>
            <p className="text-muted-foreground text-lg">
              Configure your AI companion
            </p>
          </div>
          <div className=" ">
            <ModeToggle />
          </div>
        </div>
        <div className="mx-auto max-w-3xl w-full space-y-8 pb-12">
          <div className="space-y-4">
            {configSections.map((section) => (
              <Card
                key={section.id}
                onClick={() => {
                  navigate(section.href || `/settings/${section.id}`);
                }}
                className="group relative overflow-hidden cursor-pointer border-2 border-transparent bg-card/50 backdrop-blur-sm transition-all hover:border-primary/50"
              >
                <CardContent className="relative z-10 flex flex-col gap-4 px-4 py-2 sm:flex-row sm:items-center sm:justify-between">
                  <div className="absolute -left-12 -top-12 -z-10 h-16 w-16 rounded-full bg-primary/20 blur-3xl transition-all duration-1000 ease-out pointer-events-none group-hover:bg-primary/40 group-hover:h-64 group-hover:w-64" />

                  <div className="flex items-start gap-4">
                    <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-primary/10 text-primary transition-colors group-hover:bg-primary group-hover:text-primary-foreground">
                      <section.icon className="h-6 w-6" />
                    </div>
                    <div className="space-y-1">
                      <CardTitle className="text-xl">{section.title}</CardTitle>
                      <CardDescription>{section.description}</CardDescription>
                      <div className="flex flex-wrap gap-2 pt-2">
                        {section.details &&
                          section.details.map((detail) => (
                            <Badge
                              key={detail.label}
                              variant="secondary"
                              className="bg-secondary/50 font-normal"
                            >
                              {detail.label}: {detail.value}
                            </Badge>
                          ))}
                      </div>
                    </div>
                  </div>

                  <div className="shrink-0 pt-4 sm:pt-0">
                    <RippleButton
                      onClick={() => {
                        navigate(section.href);
                      }}
                      className="w-full sm:w-auto bg-primary/10 text-primary hover:bg-primary hover:text-primary-foreground"
                    >
                      Configure
                    </RippleButton>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    </AuroraBackground>
  );
}
