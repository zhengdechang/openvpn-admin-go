// In openvpn-web/src/app/dashboard/logs/page.tsx
"use client";

import React, { useEffect, useState, useCallback } from "react"; // Added useCallback
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/services/api";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
// Ensure LiveClientConnection is imported
import { LiveClientConnection } from "@/types/types";
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table";

// Helper function to format duration from seconds to a readable string (already exists)
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

// Helper function to format bytes to a readable string (already exists)
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
  const [loadingServer, setLoadingServer] = useState(true);

  // New state for Live Client Connections
  const [liveConnections, setLiveConnections] = useState<LiveClientConnection[]>([]);
  const [loadingLiveConnections, setLoadingLiveConnections] = useState(true);
  const [errorLiveConnections, setErrorLiveConnections] = useState<string | null>(null);

  // Fetch server logs (existing useEffect)
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

  // New useEffect for fetching live client connections
  const fetchLiveConnections = useCallback(async () => {
    setLoadingLiveConnections(true);
    setErrorLiveConnections(null);
    try {
      const data = await openvpnAPI.getLiveClientConnections();
      setLiveConnections(data || []);
    } catch (error) {
      console.error("Failed to fetch live connections:", error);
      setErrorLiveConnections(t("dashboard.logs.fetchLiveConnectionsError"));
      // toast.error(t("dashboard.logs.fetchLiveConnectionsError")); // Can be too noisy with polling
      setLiveConnections([]); // Clear data on error
    } finally {
      setLoadingLiveConnections(false);
    }
  }, [t]); // Added t to dependency array

  // useEffect(() => {
  //    fetchLiveConnections(); // Initial fetch
  //    const intervalId = setInterval(fetchLiveConnections, 10000); // Poll every 10 seconds
  //    return () => clearInterval(intervalId); // Cleanup on unmount
  // }, [fetchLiveConnections]);


  return (
    <MainLayout className="p-4 space-y-6">
      {/* Server Logs Card (existing) */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.logs.serverLogsCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingServer ? (
            <p>{t("common.loading")}</p>
          ) : (
            <pre className="whitespace-pre-wrap break-all">
              {serverLogs || t("dashboard.logs.noServerLogs")}
            </pre>
          )}
        </CardContent>
      </Card>

      {/* New Live Client Connections Card */}
      {/* <Card>
        <CardHeader className="flex flex-row items-center justify-between">
         <CardTitle>{t("dashboard.logs.liveConnectionsCardTitle")}</CardTitle>
         <Button onClick={fetchLiveConnections} disabled={loadingLiveConnections} size="sm">
             {loadingLiveConnections ? t("common.refreshing") : t("common.refresh")}
         </Button>
        </CardHeader>
        <CardContent>
          {loadingLiveConnections && liveConnections.length === 0 ? ( // Show loading only on initial load or if data is empty
            <p>{t("common.loading")}</p>
          ) : errorLiveConnections ? (
            <p className="text-red-500">{errorLiveConnections}</p>
          ) : liveConnections.length === 0 ? (
            <p>{t("dashboard.logs.noLiveConnections")}</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>{t("dashboard.logs.liveColumnUserId")}</TableHead>
                  <TableHead>{t("dashboard.logs.liveColumnConnectionIp")}</TableHead>
                  <TableHead>{t("dashboard.logs.liveColumnVpnIp")}</TableHead>
                  <TableHead>{t("dashboard.logs.liveColumnOnlineDuration")}</TableHead>
                  <TableHead>{t("dashboard.logs.liveColumnBytesSent")}</TableHead>
                  <TableHead>{t("dashboard.logs.liveColumnBytesReceived")}</TableHead>
                  <TableHead>{t("dashboard.logs.liveColumnConnectedSince")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {liveConnections.map((conn) => (
                  <TableRow key={conn.commonName + conn.realAddress}>
                    <TableCell>{conn.commonName}</TableCell>
                    <TableCell>{conn.realAddress}</TableCell>
                    <TableCell>{conn.virtualAddress || t("common.na")}</TableCell>
                    <TableCell>{formatDuration(conn.onlineDurationSeconds)}</TableCell>
                    <TableCell>{formatBytes(conn.bytesSent)}</TableCell>
                    <TableCell>{formatBytes(conn.bytesReceived)}</TableCell>
                    <TableCell>
                      {conn.connectedSince
                        ? new Date(conn.connectedSince).toLocaleString()
                        : t("common.na")}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card> */}
    </MainLayout>
  );
}
