import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import ModelCard from "@/components/model-card";
import type { Model as Live2dModel } from "@wailsbindings/services";
import type { ILive2DModel } from "@/types/models";

interface ModelListProps {
  models?: Live2dModel[];
  selectedModel?: Live2dModel & { data: ILive2DModel };
  onSelectModel: (
    model: (Live2dModel & { data: ILive2DModel }) | undefined
  ) => void;
  isLoading?: boolean;
}

export default function ModelList({
  models,
  selectedModel,
  onSelectModel,
  isLoading,
}: ModelListProps) {
  return (
    <Card className="border-2 border-primary/20 bg-card/50 backdrop-blur-sm relative overflow-hidden">
      <div className="absolute -left-12 -bottom-12 -z-10 h-32 w-32 rounded-full bg-primary/10 blur-3xl" />

      <CardHeader>
        <CardTitle className="text-xl">Available Models</CardTitle>
        <CardDescription>
          {isLoading ? "Loading..." : `${models?.length ?? 0} models installed`}
        </CardDescription>
      </CardHeader>

      <CardContent className="space-y-3 max-h-[500px] overflow-y-auto pr-2">
        {isLoading ? (
          <div className="flex items-center justify-center py-8">
            <p className="text-muted-foreground">Loading models...</p>
          </div>
        ) : models && models.length > 0 ? (
          models.map((model) => (
            <ModelCard
              key={model.id}
              model={model}
              setSelectedModel={onSelectModel}
              selectedModel={selectedModel}
            />
          ))
        ) : (
          <div className="flex items-center justify-center py-8">
            <p className="text-muted-foreground">No models available</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
