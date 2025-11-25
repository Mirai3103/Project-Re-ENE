import React, { useRef, useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Sparkles, Play, Smile } from "lucide-react";
import { cn } from "@/lib/utils";
import { Application, Ticker } from "pixi.js";
import { Live2DModel } from "@laffy1309/pixi-live2d-lipsyncpatch";
import type { Model as Live2dModel } from "@wailsbindings/services";

import type { ILive2DModel } from "@/types/models";
import { RippleButton } from "./ui/shadcn-io/ripple-button";

interface ModelPreviewProps {
  model?: Live2dModel & { data: ILive2DModel };
}

export default function ModelPreview({ model }: ModelPreviewProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const modelRef = useRef<Live2DModel | null>(null);

  const validExtensions = React.useMemo(() => {
    if (!model?.data) return { expressions: [], motions: {} };
    const exps =
      model.data.FileReferences.Expressions?.map((exp) => exp.Name) ?? [];
    const motionGroups = model.data.FileReferences.Motions ?? {};
    return {
      expressions: exps,
      motions: motionGroups,
    };
  }, [model?.data]);

  function setMotion(group: string, index: number) {
    if (!modelRef.current) return;
    modelRef.current.motion(group, index);
  }

  function setExpression(expression: string) {
    if (!modelRef.current) return;
    modelRef.current.expression(expression);
  }
  useEffect(() => {
    if (!model?.id || !containerRef.current) return;

    let app: Application | null = null;
    let live2dModel: Live2DModel | null = null;

    async function initializeModel() {
      try {
        if (!containerRef.current) return;

        // Xóa PIXI canvas cũ nếu tồn tại
        containerRef.current.innerHTML = "";

        // Tạo PIXI App (tự tạo <canvas>)
        app = new Application({
          autoStart: true,
          backgroundAlpha: 0,
          resizeTo: containerRef.current, // tự resize theo div
        });

        // Gắn PIXI canvas vào div
        containerRef.current.appendChild(app.view as any);

        // Tải model
        live2dModel = await Live2DModel.from(`/models${model!.path}`, {
          ticker: Ticker.shared,
        });
        modelRef.current = live2dModel;
        app.stage.addChild(live2dModel);
        live2dModel.anchor.set(0.5);

        // Căn giữa
        live2dModel.x = app.renderer.width / 2;
        live2dModel.y = app.renderer.height / 2;

        // Scale vừa khung
        const scale =
          Math.min(
            app.renderer.width / live2dModel.width,
            app.renderer.height / live2dModel.height
          ) * 0.9;
        live2dModel.scale.set(scale);
      } catch (error) {
        console.error("Failed to load Live2D model:", error);
      }
    }

    initializeModel();

    return () => {
      if (live2dModel) {
        live2dModel.destroy({
          baseTexture: true,
          children: true,
          texture: true,
        });
      }

      if (app) {
        app.destroy(true);
      }
    };
  }, [model?.id, model?.path]);

  return (
    <Card className="border-2 border-primary/20 bg-card/50 backdrop-blur-sm relative overflow-hidden lg:sticky lg:top-4">
      <div className="absolute inset-0 -z-10 bg-linear-to-br from-primary/5 via-transparent to-purple-500/5" />

      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-xl">
          <Sparkles className="h-5 w-5 text-primary" />
          Model Preview
        </CardTitle>
      </CardHeader>

      <CardContent className="space-y-6">
        {/* Preview Wrapper */}
        <div className="aspect-square rounded-2xl overflow-hidden bg-linear-to-br from-pink-100/50 via-purple-100/50 to-blue-100/50 dark:from-pink-900/20 dark:via-purple-900/20 dark:to-blue-900/20 border-2 border-primary/20 flex items-center justify-center relative">
          {/* Container để PIXI tự tạo canvas */}
          <div
            ref={containerRef}
            className={cn("w-full h-full", model?.id ? "block" : "hidden")}
          />

          {!model?.id && (
            <div className="absolute inset-0 flex flex-col items-center justify-center text-muted-foreground">
              <Sparkles className="h-16 w-16 mb-4 opacity-30" />
              <p className="text-sm opacity-50">Select a model to preview</p>
            </div>
          )}
        </div>

        {/* Motion & Expression Controls */}
        {model?.id && (
          <div className="space-y-4">
            {/* Expressions Section */}
            {validExtensions.expressions.length > 0 && (
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <Smile className="h-4 w-4 text-primary" />
                  <h3 className="font-semibold text-sm text-foreground">
                    Expressions
                  </h3>
                  <Badge variant="secondary" className="text-xs">
                    {validExtensions.expressions.length}
                  </Badge>
                </div>
                <div className="flex flex-wrap gap-2">
                  {validExtensions.expressions.map((exp) => (
                    <RippleButton
                      variant="outline"
                      key={exp}
                      onClick={() => setExpression(exp)}
                      className={cn(
                        "px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200"
                      )}
                    >
                      {exp}
                    </RippleButton>
                  ))}
                </div>
              </div>
            )}

            {/* Motions Section */}
            {Object.keys(validExtensions.motions).length > 0 && (
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <Play className="h-4 w-4 text-primary" />
                  <h3 className="font-semibold text-sm text-foreground">
                    Motions
                  </h3>
                  <Badge variant="secondary" className="text-xs">
                    {Object.keys(validExtensions.motions).length} groups
                  </Badge>
                </div>
                <div className="space-y-3 max-h-[300px] overflow-y-auto pr-2">
                  {Object.entries(validExtensions.motions).map(
                    ([group, motions]) => (
                      <div key={group} className="space-y-2">
                        <p className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
                          {group}
                        </p>
                        <div className="flex flex-wrap gap-2">
                          {motions.map((motion, index) => {
                            return (
                              <RippleButton
                                key={`${group}-${index}`}
                                variant="outline"
                                onClick={() => setMotion(group, index)}
                                className={cn(
                                  "px-3 py-1.5 rounded-lg text-sm font-medium transition-all duration-200",
                                  "border-2 hover:scale-105 active:scale-95"
                                )}
                                title={motion.File}
                              >
                                {motion.File.split("/")
                                  .pop()
                                  ?.replace(".motion3.json", "") ||
                                  `Motion ${index + 1}`}
                              </RippleButton>
                            );
                          })}
                        </div>
                      </div>
                    )
                  )}
                </div>
              </div>
            )}

            {/* No animations available */}
            {validExtensions.expressions.length === 0 &&
              Object.keys(validExtensions.motions).length === 0 && (
                <div className="text-center py-8 text-muted-foreground">
                  <p className="text-sm">
                    No animations available for this model
                  </p>
                </div>
              )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
