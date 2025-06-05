"use client";

"use client";
import React, { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import MainLayout from "@/components/layout/main-layout";
import { useTranslation } from "react-i18next";
import { serverAPI } from "@/services/api";
import type { ServerStatus } from "@/types/types";
import { UserRole } from "@/types/types";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

export default function ServerPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation([]);
  const [status, setStatus] = useState<ServerStatus | null>(null);
  const [config, setConfig] = useState<string>("");
  const [loading, setLoading] = useState(true);

  const fetchStatus = async () => {
    setLoading(true);
    try {
      const data = await serverAPI.getStatus();
      setStatus(data);
      const tpl = await serverAPI.getConfigTemplate();
      setConfig(tpl.template);
    } catch (error) {
      toast.error(t("dashboard.server.fetchStatusError"));
    } finally {
      setLoading(false);
    }
  };
  useEffect(() => { fetchStatus(); }, []);

  if (!currentUser || currentUser.role !== UserRole.SUPERADMIN) {
    return (
      <MainLayout className="p-4">
        <p className="text-center mt-10">{t("dashboard.server.noPermission")}</p>
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
              <p>{t("dashboard.server.labelName")}{status.name}</p>
              <p>{t("dashboard.server.labelStatus")}{status.status}</p>
              <p>{t("dashboard.server.labelUptime")}{status.uptime}</p>
              <p>{t("dashboard.server.labelConnected")}{status.connected}</p>
              <p>{t("dashboard.server.labelTotal")}{status.total}</p>
              <p>{t("dashboard.server.labelLastUpdated")}{new Date(status.lastUpdated).toLocaleString()}</p>
            </div>
          ) : (
            <p>{t("dashboard.server.noData")}</p>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.server.configCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          <textarea
            className="w-full h-64 border p-2"
            value={config}
            onChange={(e) => setConfig(e.target.value)}
          />
          <div className="mt-4 space-x-2">
            <Button onClick={async () => {
              try {
                await serverAPI.updateConfig(config);
                toast.success(t("dashboard.server.updateConfigSuccess"));
              } catch {
                toast.error(t("dashboard.server.updateConfigError"));
              }
            }}>
              {t("dashboard.server.saveConfigButton")}
            </Button>
            <Button onClick={async () => {
              try { await serverAPI.start(); toast.success(t("dashboard.server.startSuccess")); fetchStatus(); } catch { toast.error(t("dashboard.server.startError")); }
            }}>
              {t("dashboard.server.startButton")}
            </Button>
            <Button onClick={async () => {
              try { await serverAPI.stop(); toast.success(t("dashboard.server.stopSuccess")); fetchStatus(); } catch { toast.error(t("dashboard.server.stopError")); }
            }}>
              {t("dashboard.server.stopButton")}
            </Button>
            <Button onClick={async () => {
              try { await serverAPI.restart(); toast.success(t("dashboard.server.restartSuccess")); fetchStatus(); } catch { toast.error(t("dashboard.server.restartError")); }
            }}>
              {t("dashboard.server.restartButton")}
            </Button>
          </div>
        </CardContent>
      </Card>
    </MainLayout>
  );
}
