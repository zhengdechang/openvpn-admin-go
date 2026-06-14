"use client";

import React, { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import MainLayout from "@/components/layout/main-layout";
import { useTranslation } from "react-i18next";
import { serverAPI } from "@/services/api";
import type { ServerStatus } from "@/types/types";
import { UserRole } from "@/types/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import ConfigManager from "@/components/config/config-manager";
import { Button } from "@/components/ui/button";

export default function ServerPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation([]);
  const [status, setStatus] = useState<ServerStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [savedConfig, setSavedConfig] = useState("");
  const [editedConfig, setEditedConfig] = useState("");
  const [rawConfigLoading, setRawConfigLoading] = useState(false);
  const [savingConfig, setSavingConfig] = useState(false);

  const fetchStatus = async () => {
    setLoading(true);
    try {
      const data = await serverAPI.getStatus();
      setStatus(data);
    } catch (error) {
      toast.error(t("dashboard.server.fetchStatusError"));
    } finally {
      setLoading(false);
    }
  };

  const fetchRawConfig = async () => {
    setRawConfigLoading(true);
    try {
      const data = await serverAPI.getRawConfig();
      setSavedConfig(data.config);
      setEditedConfig(data.config);
    } catch {
      toast.error(t("dashboard.server.configLoadError", "Failed to load server.conf"));
    } finally {
      setRawConfigLoading(false);
    }
  };

  const handleSaveConfig = async () => {
    setSavingConfig(true);
    try {
      await serverAPI.updateConfig(editedConfig);
      setSavedConfig(editedConfig);
      toast.success(t("dashboard.server.configSaveSuccess", "server.conf saved successfully"));
    } catch {
      toast.error(t("dashboard.server.configSaveError", "Failed to save server.conf"));
    } finally {
      setSavingConfig(false);
    }
  };

  useEffect(() => {
    fetchStatus();
    fetchRawConfig();
  }, []);

  if (!currentUser || currentUser.role !== UserRole.SUPERADMIN) {
    return (
      <MainLayout className="p-4">
        <p className="text-center mt-10">
          {t("dashboard.server.noPermission")}
        </p>
      </MainLayout>
    );
  }
  return (
    <MainLayout className="p-4 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.server.statusCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p>{t("common.loading")}</p>
          ) : status ? (
            <div className="space-y-2">
              <p>
                {t("dashboard.server.labelName")}
                {status.name}
              </p>
              <p>
                {t("dashboard.server.labelStatus")}
                {status.status}
              </p>
              <p>
                {t("dashboard.server.labelUptime")}
                {status.uptime}
              </p>
              <p>
                {t("dashboard.server.labelConnected")}
                {status.connected}
              </p>
              <p>
                {t("dashboard.server.labelTotal")}
                {status.total}
              </p>
              <p>
                {t("dashboard.server.labelLastUpdated")}
                {new Date(status.lastUpdated).toLocaleString()}
              </p>
            </div>
          ) : (
            <p>{t("dashboard.server.noData")}</p>
          )}
        </CardContent>
      </Card>
      {/* 配置管理组件 */}
      <ConfigManager />

      {/* 服务器控制 */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.server.controlCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-x-2">
            <Button
              variant="default"
              onClick={async () => {
                try {
                  await serverAPI.start();
                  toast.success(t("dashboard.server.startSuccess"));
                  fetchStatus();
                } catch {
                  toast.error(t("dashboard.server.startError"));
                }
              }}
            >
              {t("dashboard.server.startButton")}
            </Button>
            <Button
              variant="default"
              onClick={async () => {
                try {
                  await serverAPI.stop();
                  toast.success(t("dashboard.server.stopSuccess"));
                  fetchStatus();
                } catch {
                  toast.error(t("dashboard.server.stopError"));
                }
              }}
            >
              {t("dashboard.server.stopButton")}
            </Button>
            <Button
              variant="default"
              onClick={async () => {
                try {
                  await serverAPI.restart();
                  toast.success(t("dashboard.server.restartSuccess"));
                  fetchStatus();
                } catch {
                  toast.error(t("dashboard.server.restartError"));
                }
              }}
            >
              {t("dashboard.server.restartButton")}
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* 原始 server.conf — 可编辑，放在最后 */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle className="font-mono text-base">server.conf</CardTitle>
          <div className="flex items-center gap-2">
            {editedConfig !== savedConfig && (
              <Button
                variant="ghost"
                size="sm"
                disabled={savingConfig}
                onClick={() => setEditedConfig(savedConfig)}
              >
                {t("common.reset", "Reset")}
              </Button>
            )}
            <Button
              size="sm"
              disabled={savingConfig || editedConfig === savedConfig || rawConfigLoading}
              onClick={handleSaveConfig}
            >
              {savingConfig ? t("common.saving", "Saving...") : t("common.save")}
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {rawConfigLoading ? (
            <p className="py-4 text-sm text-muted-foreground">{t("common.loading")}</p>
          ) : (
            <textarea
              className="w-full rounded-md bg-muted px-4 py-3 text-xs leading-relaxed font-mono resize-y min-h-[420px] border border-transparent focus:border-ring focus:outline-none focus:ring-1 focus:ring-ring"
              value={editedConfig}
              onChange={(e) => setEditedConfig(e.target.value)}
              spellCheck={false}
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
            />
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
