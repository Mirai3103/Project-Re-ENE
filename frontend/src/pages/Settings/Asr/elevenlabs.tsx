import { useQuery } from "@/lib/query";
import { ConfigService } from "@wailsbindings/services";
import React from "react";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { RippleButton } from "@/components/ui/shadcn-io/ripple-button";
import { AuroraBackground } from "@/components/ui/shadcn-io/aurora-background";
import { useLocation } from "wouter";
import {
  ArrowLeft,
  Sparkles,
  Key,
  Settings2,
  CheckCircle,
  Info,
} from "lucide-react";

const elevenLabsConfigSchema = z.object({
  name: z.string().optional(),
  apiKey: z.string().min(1, "API Key is required"),
  modelId: z.string().min(1, "Model ID is required"),
  languageCode: z.string(),
});

type ElevenLabsConfigFormValues = z.infer<typeof elevenLabsConfigSchema>;

const LANGUAGE_CODES = [
  {
    value: "auto",
    label: "Auto Detect",
    description: "Automatically detect language",
  },
  { value: "en", label: "English", description: "English (US/UK)" },
  { value: "vi", label: "Vietnamese", description: "Tiếng Việt" },
  { value: "zh", label: "Chinese", description: "中文" },
  { value: "ja", label: "Japanese", description: "日本語" },
  { value: "ko", label: "Korean", description: "한국어" },
  { value: "es", label: "Spanish", description: "Español" },
  { value: "fr", label: "French", description: "Français" },
  { value: "de", label: "German", description: "Deutsch" },
];

const MODELS = [
  {
    value: "scribe_v1",
    label: "Scribe v1",
    description: "Original model with stable performance",
  },
  {
    value: "scribe_v2",
    label: "Scribe v2",
    description: "Improved accuracy and speed",
  },
];

export default function ElevenLabsConfig() {
  const [, navigate] = useLocation();
  const [isSaving, setIsSaving] = React.useState(false);

  const { data: config, refetch } = useQuery({
    queryKey: ["configs"],
    queryFn: () => ConfigService.GetConfig(),
  });

  const elevenLabsConfig = config?.ASRConfig.ElevenLabsConfig;

  const form = useForm<ElevenLabsConfigFormValues>({
    resolver: zodResolver(elevenLabsConfigSchema),
    defaultValues: {
      name: elevenLabsConfig?.Name || "",
      apiKey: elevenLabsConfig?.APIKey || "",
      modelId: elevenLabsConfig?.ModelID || "scribe_v1",
      languageCode: elevenLabsConfig?.LanguageCode || "auto",
    },
  });

  React.useEffect(() => {
    if (elevenLabsConfig) {
      form.reset({
        name: elevenLabsConfig.Name || "",
        apiKey: elevenLabsConfig.APIKey || "",
        modelId: elevenLabsConfig.ModelID || "scribe_v1",
        languageCode: elevenLabsConfig.LanguageCode || "auto",
      });
    }
  }, [elevenLabsConfig, form]);

  async function onSubmit(data: ElevenLabsConfigFormValues) {
    try {
      setIsSaving(true);
      await ConfigService.PatchConfig({
        ASRConfig: {
          ElevenLabsConfig: {
            Name: data.name || "",
            APIKey: data.apiKey,
            ModelID: data.modelId,
            LanguageCode: data.languageCode || "auto",
          },
        },
      } as any);

      await refetch();
      alert("Configuration saved successfully!");
    } catch (error) {
      console.error("Failed to save config:", error);
      alert("Failed to save configuration. Please try again.");
    } finally {
      setIsSaving(false);
    }
  }

  return (
    <AuroraBackground showRadialGradient className="overflow-auto">
      <div className="min-h-screen w-full max-w-5xl p-4 py-10 font-sans selection:bg-primary/20">
        {/* Header */}
        <div className="relative space-y-2 pb-8 flex items-center justify-between">
          <RippleButton
            onClick={() => {
              navigate("/settings/asr");
            }}
            variant="outline"
            className="bg-primary/10 text-primary hover:bg-primary hover:text-primary-foreground"
          >
            <ArrowLeft className="h-5 w-5 mr-2" />
            Back
          </RippleButton>

          <div className="text-center flex-1">
            <div className="flex items-center justify-center gap-3 mb-2">
              <div className="w-12 h-12 rounded-2xl bg-linear-to-br from-purple-500 to-pink-500 flex items-center justify-center text-2xl shadow-lg">
                🎙️
              </div>
              <h1 className="text-4xl font-bold tracking-tight text-primary drop-shadow-sm">
                ElevenLabs ASR
              </h1>
            </div>
            <p className="text-muted-foreground text-lg">
              Configure ElevenLabs Speech Recognition
            </p>
          </div>

          <div className="w-[100px]"></div>
        </div>

        {/* Info Banner */}
        <div className="mx-auto max-w-3xl mb-6">
          <div className="bg-primary/10 border-2 border-primary/20 rounded-2xl p-4 flex items-start gap-3">
            <Info className="h-5 w-5 text-primary shrink-0 mt-0.5" />
            <div className="text-sm">
              <p className="font-semibold text-foreground mb-1">
                Get your API Key
              </p>
              <p className="text-muted-foreground">
                Visit{" "}
                <a
                  href="https://elevenlabs.io/api"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-primary hover:underline font-medium"
                >
                  elevenlabs.io/api
                </a>{" "}
                to get your API key for speech recognition services.
              </p>
            </div>
          </div>
        </div>

        {/* Form Card */}
        <div className="mx-auto max-w-3xl">
          <Card className="group relative overflow-hidden border-2 border-primary/20 bg-card/50 backdrop-blur-sm transition-all">
            <div className="absolute -left-12 -top-12 -z-10 h-32 w-32 rounded-full bg-primary/20 blur-3xl transition-all duration-1000 ease-out pointer-events-none" />
            <div className="absolute -right-12 -bottom-12 -z-10 h-32 w-32 rounded-full bg-primary/10 blur-3xl transition-all duration-1000 ease-out pointer-events-none" />

            <CardHeader>
              <CardTitle className="text-2xl flex items-center gap-2">
                <Settings2 className="h-6 w-6 text-primary" />
                Configuration Settings
              </CardTitle>
              <CardDescription>
                Configure your ElevenLabs ASR provider settings
              </CardDescription>
            </CardHeader>

            <CardContent>
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="space-y-6"
                >
                  {/* Name Field */}
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold">
                          Configuration Name
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder="My ElevenLabs Config (optional)"
                            className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 focus:border-primary transition-colors"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          Optional name to identify this configuration
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* API Key Field */}
                  <FormField
                    control={form.control}
                    name="apiKey"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold flex items-center gap-2">
                          <Key className="h-4 w-4 text-primary" />
                          API Key
                          <Badge variant="destructive" className="text-xs">
                            Required
                          </Badge>
                        </FormLabel>
                        <FormControl>
                          <Input
                            type="password"
                            placeholder="sk_..."
                            className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 focus:border-primary transition-colors font-mono"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          Your ElevenLabs API key for authentication
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Model ID Selection */}
                  <FormField
                    control={form.control}
                    name="modelId"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold">
                          Model
                        </FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                        >
                          <FormControl>
                            <SelectTrigger className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 transition-colors w-full">
                              <SelectValue placeholder="Select a model" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {MODELS.map((model) => (
                              <SelectItem key={model.value} value={model.value}>
                                <div className="flex flex-col">
                                  <span className="font-medium">
                                    {model.label}
                                  </span>
                                </div>
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormDescription>
                          Choose the speech recognition model version
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Language Code Selection */}
                  <FormField
                    control={form.control}
                    name="languageCode"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold">
                          Language
                        </FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                        >
                          <FormControl>
                            <SelectTrigger className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 transition-colors w-full">
                              <SelectValue placeholder="Select a language" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {LANGUAGE_CODES.map((lang) => (
                              <SelectItem key={lang.value} value={lang.value}>
                                <div className="flex flex-col">
                                  <span className="font-medium">
                                    {lang.label}
                                  </span>
                                </div>
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormDescription>
                          Target language for speech recognition
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Submit Buttons */}
                  <div className="flex gap-3 pt-4">
                    <RippleButton
                      type="submit"
                      disabled={isSaving}
                      className="flex-1 bg-primary text-primary-foreground hover:bg-primary/90 h-12 text-base font-semibold"
                    >
                      {isSaving ? (
                        <>
                          <div className="h-5 w-5 mr-2 animate-spin rounded-full border-2 border-primary-foreground border-t-transparent" />
                          Saving...
                        </>
                      ) : (
                        <>
                          <CheckCircle className="h-5 w-5 mr-2" />
                          Save Configuration
                        </>
                      )}
                    </RippleButton>
                    <RippleButton
                      type="button"
                      variant="outline"
                      onClick={() => form.reset()}
                      className="bg-secondary/50 hover:bg-secondary"
                      disabled={isSaving}
                    >
                      Reset
                    </RippleButton>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>
        </div>
      </div>
    </AuroraBackground>
  );
}
