import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
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
import { Textarea } from "@/components/ui/textarea";
import { RippleButton } from "@/components/ui/shadcn-io/ripple-button";
import { AuroraBackground } from "@/components/ui/shadcn-io/aurora-background";
import { useLocation } from "wouter";

import { ArrowLeft, User, Sparkles } from "lucide-react";

// Mock data for Live2D models
const live2dModels = [
  { id: "hiyori_pro_t11", name: "Hiyori Pro T11", preview: "Current model" },
  { id: "shizuku_model", name: "Shizuku", preview: "Alternative character" },
  { id: "miku_model", name: "Miku", preview: "Vocaloid style" },
  { id: "ene_model", name: "Ene", preview: "Default character" },
];

// Zod schema for form validation
const characterConfigSchema = z.object({
  live2d_model_name: z.string().min(1, "Please select a Live2D model"),
  character_name: z
    .string()
    .min(1, "Character name is required")
    .max(50, "Character name must be less than 50 characters"),
  user_name: z
    .string()
    .min(1, "User name is required")
    .max(50, "User name must be less than 50 characters"),
  persona_prompt: z
    .string()
    .min(10, "Persona prompt must be at least 10 characters")
    .max(1000, "Persona prompt must be less than 1000 characters"),
});

type CharacterConfigFormValues = z.infer<typeof characterConfigSchema>;

// Mock default values
const defaultValues: CharacterConfigFormValues = {
  live2d_model_name: "hiyori_pro_t11",
  character_name: "Ene",
  user_name: "Master",
  persona_prompt:
    "You are Ene, a cheerful and energetic AI assistant with a playful personality. You love helping your master with tasks and enjoy casual conversations. You occasionally use cute expressions and are always enthusiastic about learning new things.",
};

export default function Character() {
  const [, navigate] = useLocation();

  const form = useForm<CharacterConfigFormValues>({
    resolver: zodResolver(characterConfigSchema),
    defaultValues,
  });

  const onSubmit = (data: CharacterConfigFormValues) => {
    console.log("Form submitted:", data);
    // TODO: Integrate with backend API to save configuration
    alert("Configuration saved! (Mock implementation)");
  };

  return (
    <AuroraBackground showRadialGradient className="overflow-auto">
      <div className="min-h-screen w-full max-w-5xl p-4 py-10 pb-20 font-sans selection:bg-primary/20">
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
              <User className="h-8 w-8 text-primary" />
              <h1 className="text-4xl font-bold tracking-tight text-primary drop-shadow-sm">
                Character Configuration
              </h1>
            </div>
            <p className="text-muted-foreground text-lg">
              Customize your AI companion's appearance and personality
            </p>
          </div>

          <div className="w-[100px]"></div>
        </div>

        {/* Form Card */}
        <div className="mx-auto max-w-3xl">
          <Card className="group relative overflow-hidden border-2 border-primary/20 bg-card/50 backdrop-blur-sm transition-all">
            <div className="absolute -left-12 -top-12 -z-10 h-32 w-32 rounded-full bg-primary/20 blur-3xl transition-all duration-1000 ease-out pointer-events-none" />
            <div className="absolute -right-12 -bottom-12 -z-10 h-32 w-32 rounded-full bg-primary/10 blur-3xl transition-all duration-1000 ease-out pointer-events-none" />

            <CardHeader>
              <CardTitle className="text-2xl flex items-center gap-2">
                <Sparkles className="h-6 w-6 text-primary" />
                Character Settings
              </CardTitle>
              <CardDescription>
                Configure your character's model, identity, and personality
                traits
              </CardDescription>
            </CardHeader>

            <CardContent>
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="space-y-6"
                >
                  {/* Live2D Model Selection */}
                  <FormField
                    control={form.control}
                    name="live2d_model_name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold">
                          Live2D Model
                        </FormLabel>
                        <Select
                          onValueChange={field.onChange}
                          defaultValue={field.value}
                        >
                          <FormControl>
                            <SelectTrigger className="bg-background/50 backdrop-blur-sm border-2 w-full  hover:border-primary/50 transition-colors">
                              <SelectValue placeholder="Select a Live2D model" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {live2dModels.map((model) => (
                              <SelectItem key={model.id} value={model.id}>
                                <div className="flex flex-col">
                                  <span className="font-medium">
                                    {model.name}
                                  </span>
                                  {/* <span className="text-xs text-muted-foreground">
                                    {model.preview}
                                  </span> */}
                                </div>
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        <FormDescription>
                          Choose the 3D avatar model for your character
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Character Name */}
                  <FormField
                    control={form.control}
                    name="character_name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold">
                          Character Name
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder="Enter character name"
                            className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 focus:border-primary transition-colors"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          The name of your AI companion
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* User Name */}
                  <FormField
                    control={form.control}
                    name="user_name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold">
                          Your Name
                        </FormLabel>
                        <FormControl>
                          <Input
                            placeholder="Enter your name"
                            className="bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 focus:border-primary transition-colors"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          How the character will address you
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Persona Prompt */}
                  <FormField
                    control={form.control}
                    name="persona_prompt"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-base font-semibold">
                          Persona Prompt
                        </FormLabel>
                        <FormControl>
                          <Textarea
                            placeholder="Describe your character's personality, traits, and behavior..."
                            className="min-h-[150px] bg-background/50 backdrop-blur-sm border-2 hover:border-primary/50 focus:border-primary transition-colors resize-none"
                            {...field}
                          />
                        </FormControl>
                        <FormDescription>
                          Define the character's personality and behavior
                          (10-1000 characters)
                        </FormDescription>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* Submit Button */}
                  <div className="flex gap-3 pt-4">
                    <RippleButton
                      type="submit"
                      className="flex-1 bg-primary text-primary-foreground hover:bg-primary/90 h-12 text-base font-semibold"
                    >
                      <Sparkles className="h-5 w-5 mr-2" />
                      Save Configuration
                    </RippleButton>
                    <RippleButton
                      type="button"
                      variant="outline"
                      onClick={() => form.reset()}
                      className="bg-secondary/50 hover:bg-secondary"
                    >
                      Reset
                    </RippleButton>
                  </div>
                </form>
              </Form>
            </CardContent>
          </Card>

          {/* Preview Card */}
          <Card className="mt-6 border-2 border-primary/10 bg-card/30 backdrop-blur-sm">
            <CardHeader>
              <CardTitle className="text-lg">Current Configuration</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <span className="text-muted-foreground">Model:</span>
                  <p className="font-medium text-foreground mt-1">
                    {live2dModels.find(
                      (m) => m.id === form.watch("live2d_model_name")
                    )?.name || "Not selected"}
                  </p>
                </div>
                <div>
                  <span className="text-muted-foreground">Character:</span>
                  <p className="font-medium text-foreground mt-1">
                    {form.watch("character_name") || "Not set"}
                  </p>
                </div>
                <div>
                  <span className="text-muted-foreground">User:</span>
                  <p className="font-medium text-foreground mt-1">
                    {form.watch("user_name") || "Not set"}
                  </p>
                </div>
                <div>
                  <span className="text-muted-foreground">Persona Length:</span>
                  <p className="font-medium text-foreground mt-1">
                    {form.watch("persona_prompt")?.length || 0} characters
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </AuroraBackground>
  );
}
