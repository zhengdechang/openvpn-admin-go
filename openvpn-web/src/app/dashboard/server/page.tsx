"use client";

import React, { useEffect, useState } from "react";
import { useAuth } from "@/lib/auth-context";
import MainLayout from "@/components/layout/main-layout";
import { useTranslation } from "react-i18next";
import { serverAPI } from "@/services/api";
import type { ServerStatus } from "@/types/types";
import { UserRole } from "@/types/types";
import { Button } from "@/components/ui/button";
import { CbiSection, CbiValue } from "@/components/ui/cbi-form";
import { toast } from "sonner";
import ConfigManager from "@/components/config/config-manager";

export default function ServerPage() {
  const { user: currentUser } = useAuth();
  const { t } = useTranslation([]);
  const [status, setStatus] = useState<ServerStatus | null>(null);
  const [loading, setLoading] = useState(true);
  // acting 防止启停期间重复点击（supervisorctl 操作有短暂窗口，状态未刷新前先锁住按钮）
  const [acting, setActing] = useState(false);

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
  useEffect(() => {
    fetchStatus();
  }, []);

  if (!currentUser || currentUser.role !== UserRole.SUPERADMIN) {
    return (
      <MainLayout className="p-6">
        <p className="text-center mt-10 text-muted-foreground">
          {t("dashboard.server.noPermission")}
        </p>
      </MainLayout>
    );
  }
  // 后端 GetServerStatus 用 "active"/"inactive" 表示 supervisor 的 RUNNING/STOPPED
  const isActive = status?.status === "active";

  return (
    <MainLayout className="p-6 space-y-6">
      {/* 服务器状态 */}
      <CbiSection title={t("dashboard.server.statusCardTitle")}>
        {loading ? (
          <p className="text-sm text-muted-foreground">{t("common.loading")}</p>
        ) : status ? (
          <>
            <CbiValue title={t("dashboard.server.labelName")}>{status.name}</CbiValue>
            <CbiValue title={t("dashboard.server.labelStatus")}>{status.status || "—"}</CbiValue>
            <CbiValue title={t("dashboard.server.labelUptime")}>{status.uptime || "—"}</CbiValue>
            <CbiValue title={t("dashboard.server.labelConnected")}>{status.connected}</CbiValue>
            <CbiValue title={t("dashboard.server.labelTotal")}>{status.total}</CbiValue>
            <CbiValue title={t("dashboard.server.labelLastUpdated")}>
              {status.lastUpdated ? new Date(status.lastUpdated).toLocaleString() : "—"}
            </CbiValue>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">{t("dashboard.server.noData")}</p>
        )}
      </CbiSection>

      {/* 配置管理组件 */}
      <ConfigManager />

      {/* 服务器控制：按运行状态启用/禁用。
          运行中 → 禁用「启动」(只能重启/停止)；已停止 → 禁用「停止」「重启」(只能启动)。 */}
      <CbiSection title={t("dashboard.server.controlCardTitle")}>
        <div className="flex flex-wrap gap-2">
          <Button
            disabled={loading || acting || isActive}
            onClick={async () => {
              setActing(true);
              try {
                await serverAPI.start();
                toast.success(t("dashboard.server.startSuccess"));
                await fetchStatus();
              } catch {
                toast.error(t("dashboard.server.startError"));
              } finally {
                setActing(false);
              }
            }}
          >
            {t("dashboard.server.startButton")}
          </Button>
          <Button
            variant="outline"
            disabled={loading || acting || !isActive}
            onClick={async () => {
              setActing(true);
              try {
                await serverAPI.stop();
                toast.success(t("dashboard.server.stopSuccess"));
                await fetchStatus();
              } catch {
                toast.error(t("dashboard.server.stopError"));
              } finally {
                setActing(false);
              }
            }}
          >
            {t("dashboard.server.stopButton")}
          </Button>
          <Button
            variant="outline"
            disabled={loading || acting || !isActive}
            onClick={async () => {
              setActing(true);
              try {
                await serverAPI.restart();
                toast.success(t("dashboard.server.restartSuccess"));
                await fetchStatus();
              } catch {
                toast.error(t("dashboard.server.restartError"));
              } finally {
                setActing(false);
              }
            }}
          >
            {t("dashboard.server.restartButton")}
          </Button>
        </div>
      </CbiSection>
    </MainLayout>
  );
}
