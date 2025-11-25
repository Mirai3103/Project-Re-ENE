import { useState } from "react";
import { RippleButton } from "@/components/ui/shadcn-io/ripple-button";
import { AuroraBackground } from "@/components/ui/shadcn-io/aurora-background";
import { useLocation } from "wouter";

import { ModelService, Model as Live2dModel } from "@wailsbindings/services";
import { ArrowLeft, Sparkles } from "lucide-react";
import { useQuery } from "@/lib/query";
import ModelUpload from "@/components/ModelUpload";
import ModelList from "@/components/ModelList";
import ModelPreview from "@/components/ModelPreview";
import type { ILive2DModel } from "@/types/models";

export default function Model() {
  const [, navigate] = useLocation();
  const [selectedModel, setSelectedModel] = useState<
    (Live2dModel & { data: ILive2DModel }) | undefined
  >(undefined);

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["models"],
    queryFn: () => ModelService.GetModelList(),
  });

  const handleUploadSuccess = () => {
    refetch();
  };

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
              <Sparkles className="h-8 w-8 text-primary" />
              <h1 className="text-4xl font-bold tracking-tight text-primary drop-shadow-sm">
                Live2D Models
              </h1>
            </div>
            <p className="text-muted-foreground text-lg">
              Manage and preview your character models
            </p>
          </div>

          <div className="w-[100px]"></div>
        </div>

        {/* Main Content - Two Column Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Left Column - Upload & List */}
          <div className="space-y-6">
            <ModelUpload onUploadSuccess={handleUploadSuccess} />
            <ModelList
              models={data}
              selectedModel={selectedModel}
              onSelectModel={setSelectedModel}
              isLoading={isLoading}
            />
          </div>

          {/* Right Column - Preview */}
          <div className="space-y-6">
            <ModelPreview model={selectedModel} />
          </div>
        </div>
      </div>
    </AuroraBackground>
  );
}
