// In openvpn-web/src/app/dashboard/logs/page.tsx
"use client";

import React, { useEffect, useState, useCallback } from "react";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/services/api";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { LiveClientConnection } from "@/types/types";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";

// Helper function to format duration from seconds to a readable string
const formatDuration = (totalSeconds: number): string => {
  if (totalSeconds < 0) return "N/A";
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;
  let result = "";
  if (hours > 0) result += `${hours}h `;
  if (minutes > 0) result += `${minutes}m `;
  if (seconds > 0 || result === "") result += `${seconds}s`;
  return result.trim() || "0s";
};

// Helper function to format bytes to a readable string
const formatBytes = (bytes: number, decimals = 2): string => {
  if (bytes < 0) return "N/A";
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
};

export default function LogsPage() {
  const { t } = useTranslation();
  const [serverLogs, setServerLogs] = useState<string>("");
  const [clientLogs, setClientLogs] = useState<string>("");
  const [loadingServer, setLoadingServer] = useState(true);
  const [loadingClient, setLoadingClient] = useState(true);

  // Fetch server logs
  useEffect(() => {
    const fetchServerLogs = async () => {
      setLoadingServer(true);
      try {
        const logs = await openvpnAPI.getServerLogs();
        setServerLogs(logs);
      } catch (error) {
        toast.error(t("dashboard.logs.fetchServerLogsError"));
      } finally {
        setLoadingServer(false);
      }
    };
    fetchServerLogs();
  }, [t]);

  // Fetch client logs
  useEffect(() => {
    const fetchClientLogs = async () => {
      setLoadingClient(true);
      try {
        const logs = await openvpnAPI.getClientLogs();
        setClientLogs(logs);
      } catch (error) {
        toast.error(t("dashboard.logs.fetchClientLogsError"));
      } finally {
        setLoadingClient(false);
      }
    };
    fetchClientLogs();
  }, [t]);

  return (
    <MainLayout className="p-4 space-y-6">
      {/* Server Logs Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.logs.serverLogsCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingServer ? (
            <p>{t("common.loading")}</p>
          ) : (
            <div className="max-h-[400px] overflow-y-auto bg-gray-100 dark:bg-gray-800 rounded-md p-4">
              <pre className="whitespace-pre-wrap break-all">
                {serverLogs || t("dashboard.logs.noServerLogs")}
              </pre>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Client Logs Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.logs.clientLogsCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingClient ? (
            <p>{t("common.loading")}</p>
          ) : (
            <div className="max-h-[400px] overflow-y-auto bg-gray-100 dark:bg-gray-800 rounded-md p-4">
              <pre className="whitespace-pre-wrap break-all">
                {clientLogs || t("dashboard.logs.noClientLogs")}
              </pre>
            </div>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
