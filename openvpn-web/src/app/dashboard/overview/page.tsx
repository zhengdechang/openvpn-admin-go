"use client";

import React, { useEffect, useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { useTranslation } from "react-i18next";
import { serverAPI } from "@/services/api";
import type { ServerStatus, SystemInfo } from "@/types/types";
import { CbiSection, CbiValue } from "@/components/ui/cbi-form";
import { Button } from "@/components/ui/button";

// LuCI 风格药丸进度条：浅灰轨道 + 渐变填充 + 居中叠加文字
function LuciProgress({ percent, text }: { percent: number; text: string }) {
  const p = Math.max(0, Math.min(100, percent || 0));
  return (
    <div className="luci-progress">
      <div className="luci-progress-fill" style={{ width: `${p}%` }} />
      <span className="luci-progress-text">{text}</span>
    </div>
  );
}

function formatBytes(n: number): string {
  if (!n || n <= 0) return "0 B";
  const units = ["B", "KiB", "MiB", "GiB", "TiB"];
  let v = n;
  let i = 0;
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024;
    i++;
  }
  return `${v.toFixed(2)} ${units[i]}`;
}

function pct(part: number, total: number): number {
  if (!total || total <= 0) return 0;
  return (part / total) * 100;
}

export default function OverviewPage() {
  const { t } = useTranslation();
  const [status, setStatus] = useState<ServerStatus | null>(null);
  const [system, setSystem] = useState<SystemInfo | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchStatus = async () => {
    setLoading(true);
    const [statusRes, systemRes] = await Promise.allSettled([
      serverAPI.getStatus(),
      serverAPI.getSystemInfo(),
    ]);
    setStatus(statusRes.status === "fulfilled" ? statusRes.value : null);
    setSystem(systemRes.status === "fulfilled" ? systemRes.value : null);
    setLoading(false);
  };

  useEffect(() => {
    fetchStatus();
  }, []);

  return (
    <MainLayout className="p-6 space-y-6">
      {/* 服务状态 */}
      <CbiSection
        title={t("dashboard.overview.statusCardTitle")}
        actions={
          <Button variant="outline" size="sm" onClick={fetchStatus} disabled={loading}>
            {t("dashboard.overview.refreshButton")}
          </Button>
        }
      >
        {loading ? (
          <p className="text-sm text-muted-foreground">{t("common.loading")}</p>
        ) : status ? (
          <>
            <CbiValue title={t("dashboard.overview.labelName")}>{status.name}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelStatus")}>{status.status}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelUptime")}>{status.uptime}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelConnected")}>{status.connected}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelTotal")}>{status.total}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelLastUpdated")}>
              {status.lastUpdated ? new Date(status.lastUpdated).toLocaleString() : "—"}
            </CbiValue>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">
            {t("dashboard.overview.unavailable")}
          </p>
        )}
      </CbiSection>

      {system && (
        <>
          {/* 系统 */}
          <CbiSection title={t("dashboard.overview.sysCardTitle")}>
            <CbiValue title={t("dashboard.overview.labelVersion")}>{system.version}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelHostname")}>{system.hostname || "—"}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelArch")}>{`${system.os} / ${system.arch}`}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelKernel")}>{system.kernelVersion || "—"}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelGoVersion")}>{system.goVersion}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelLocalTime")}>{system.localTime}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelHostUptime")}>{system.uptime || "—"}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelCpuCores")}>{system.numCpu}</CbiValue>
            <CbiValue title={t("dashboard.overview.labelLoadAvg")}>
              {system.loadAvg && system.loadAvg.length
                ? system.loadAvg.map((v) => v.toFixed(2)).join(", ")
                : "—"}
            </CbiValue>
          </CbiSection>

          {/* CPU */}
          <CbiSection title={t("dashboard.overview.cpuCardTitle")}>
            <CbiValue title={t("dashboard.overview.labelCpuUsage")}>
              <LuciProgress
                percent={system.cpuUsagePercent}
                text={`${system.cpuUsagePercent.toFixed(0)}% / 100%`}
              />
            </CbiValue>
          </CbiSection>

          {/* 内存 */}
          <CbiSection title={t("dashboard.overview.memCardTitle")}>
            <CbiValue title={t("dashboard.overview.labelMemUsed")}>
              <LuciProgress
                percent={system.memory.usedPercent}
                text={`${formatBytes(system.memory.used)} / ${formatBytes(system.memory.total)} (${system.memory.usedPercent.toFixed(0)}%)`}
              />
            </CbiValue>
            <CbiValue title={t("dashboard.overview.labelMemBuffers")}>
              <LuciProgress
                percent={pct(system.memory.buffers, system.memory.total)}
                text={`${formatBytes(system.memory.buffers)} / ${formatBytes(system.memory.total)} (${pct(system.memory.buffers, system.memory.total).toFixed(0)}%)`}
              />
            </CbiValue>
            <CbiValue title={t("dashboard.overview.labelMemCached")}>
              <LuciProgress
                percent={pct(system.memory.cached, system.memory.total)}
                text={`${formatBytes(system.memory.cached)} / ${formatBytes(system.memory.total)} (${pct(system.memory.cached, system.memory.total).toFixed(0)}%)`}
              />
            </CbiValue>
            <CbiValue title={t("dashboard.overview.labelSwap")}>
              <LuciProgress
                percent={pct(system.memory.swapUsed, system.memory.swapTotal)}
                text={`${formatBytes(system.memory.swapUsed)} / ${formatBytes(system.memory.swapTotal)} (${pct(system.memory.swapUsed, system.memory.swapTotal).toFixed(0)}%)`}
              />
            </CbiValue>
          </CbiSection>

          {/* 存储空间使用 */}
          <CbiSection title={t("dashboard.overview.storageCardTitle")}>
            <CbiValue title={`${t("dashboard.overview.labelDisk")} (${system.disk.path})`}>
              <LuciProgress
                percent={system.disk.usedPercent}
                text={`${formatBytes(system.disk.used)} / ${formatBytes(system.disk.total)} (${system.disk.usedPercent.toFixed(0)}%)`}
              />
            </CbiValue>
          </CbiSection>
        </>
      )}
    </MainLayout>
  );
}
