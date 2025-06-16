// In openvpn-web/src/app/dashboard/logs/page.tsx
"use client";

import React, { useEffect, useState, useCallback } from "react";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/services/api";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input"; // May not be needed if using select for page size
import { Label } from "@/components/ui/label"; // For pagination UI
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

  // Server Logs State
  const [serverLogLines, setServerLogLines] = useState<string[]>([]);
  const [serverCurrentPage, setServerCurrentPage] = useState(1);
  const [serverTotalPages, setServerTotalPages] = useState(1);
  const [serverTotalItems, setServerTotalItems] = useState(0);
  const [serverPageSize, setServerPageSize] = useState(100); // Default page size for server logs
  const [loadingServer, setLoadingServer] = useState(true);

  // Client Logs State
  const [clientLogLines, setClientLogLines] = useState<string[]>([]);
  const [clientCurrentPage, setClientCurrentPage] = useState(1);
  const [clientTotalPages, setClientTotalPages] = useState(1);
  const [clientTotalItems, setClientTotalItems] = useState(0);
  const [clientPageSize, setClientPageSize] = useState(100); // Default page size for client logs
  const [loadingClient, setLoadingClient] = useState(true);


  // Fetch server logs (parameterized function)
  const fetchServerLogsData = useCallback(async (page: number, size: number) => {
    setLoadingServer(true);
    try {
      const response = await openvpnAPI.getServerLogs(page, size);
      setServerLogLines(response.logs);
      setServerTotalItems(response.totalItems);
      setServerTotalPages(response.totalPages > 0 ? response.totalPages : 1);
      setServerCurrentPage(response.currentPage);
      // serverPageSize is already managed by its own state
    } catch (error) {
      toast.error(t("dashboard.logs.fetchServerLogsError"));
      setServerLogLines([]); // Clear logs on error
      setServerTotalItems(0);
      setServerTotalPages(1);
    } finally {
      setLoadingServer(false);
    }
  }, [t]); // Added t to dependencies for toast translation

  useEffect(() => {
    fetchServerLogsData(serverCurrentPage, serverPageSize);
  }, [serverCurrentPage, serverPageSize, fetchServerLogsData]);

  // Fetch client logs (parameterized function)
  const fetchClientLogsData = useCallback(async (page: number, size: number) => {
    setLoadingClient(true);
    try {
      const response = await openvpnAPI.getClientLogs(page, size);
      setClientLogLines(response.logs);
      setClientTotalItems(response.totalItems);
      setClientTotalPages(response.totalPages > 0 ? response.totalPages : 1);
      setClientCurrentPage(response.currentPage);
      // clientPageSize is already managed by its own state
    } catch (error) {
      toast.error(t("dashboard.logs.fetchClientLogsError"));
      setClientLogLines([]); // Clear logs on error
      setClientTotalItems(0);
      setClientTotalPages(1);
    } finally {
      setLoadingClient(false);
    }
  }, [t]); // Added t for toast

  useEffect(() => {
    fetchClientLogsData(clientCurrentPage, clientPageSize);
  }, [clientCurrentPage, clientPageSize, fetchClientLogsData]);

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
              <pre className="whitespace-pre-wrap break-all text-xs">
                {serverLogLines.length > 0 ? serverLogLines.join("\n") : t("dashboard.logs.noServerLogs")}
              </pre>
            </div>
            {/* Server Logs Pagination UI */}
            {!loadingServer && serverTotalItems > 0 && (
              <div className="flex items-center justify-between mt-4 pt-2 border-t">
                <div className="text-sm text-muted-foreground space-x-2">
                  <span>{t("dashboard.pagination.totalItems", { count: serverTotalItems })}</span>
                  <span>|</span>
                  <span>{t("dashboard.pagination.pageInfo", { currentPage: serverCurrentPage, totalPages: serverTotalPages })}</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Label htmlFor="serverLogsPageSizeSelect" className="sr-only">{t("dashboard.pagination.pageSizeLabel")}</Label>
                  <select
                    id="serverLogsPageSizeSelect"
                    value={serverPageSize}
                    onChange={(e) => {
                      setServerPageSize(Number(e.target.value));
                      setServerCurrentPage(1); // Reset to first page
                    }}
                    className="border px-2 py-1.5 rounded-md text-sm h-8 bg-background focus:ring-ring focus:border-input"
                  >
                    {[50, 100, 200, 500, 1000].map(size => (
                      <option key={size} value={size}>{t("dashboard.pagination.pageSizeOption", { count: size })}</option>
                    ))}
                  </select>
                  <Button
                    variant="outline" size="sm"
                    onClick={() => setServerCurrentPage(p => Math.max(1, p - 1))}
                    disabled={serverCurrentPage === 1}
                    className="h-8 px-2"
                  >
                    {t("dashboard.pagination.previousButton", "Previous")}
                  </Button>
                  <Button
                    variant="outline" size="sm"
                    onClick={() => setServerCurrentPage(p => Math.min(serverTotalPages, p + 1))}
                    disabled={serverCurrentPage === serverTotalPages || serverTotalPages === 0}
                    className="h-8 px-2"
                  >
                    {t("dashboard.pagination.nextButton", "Next")}
                  </Button>
                </div>
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
            <>
              <div className="max-h-[400px] overflow-y-auto bg-gray-100 dark:bg-gray-800 rounded-md p-4">
                <pre className="whitespace-pre-wrap break-all text-xs">
                  {clientLogLines.length > 0 ? clientLogLines.join("\n") : t("dashboard.logs.noClientLogs")}
                </pre>
              </div>
              {/* Client Logs Pagination UI */}
              {!loadingClient && clientTotalItems > 0 && (
                <div className="flex items-center justify-between mt-4 pt-2 border-t">
                  <div className="text-sm text-muted-foreground space-x-2">
                    <span>{t("dashboard.pagination.totalItems", { count: clientTotalItems })}</span>
                    <span>|</span>
                    <span>{t("dashboard.pagination.pageInfo", { currentPage: clientCurrentPage, totalPages: clientTotalPages })}</span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Label htmlFor="clientLogsPageSizeSelect" className="sr-only">{t("dashboard.pagination.pageSizeLabel")}</Label>
                    <select
                      id="clientLogsPageSizeSelect"
                      value={clientPageSize}
                      onChange={(e) => {
                        setClientPageSize(Number(e.target.value));
                        setClientCurrentPage(1); // Reset to first page
                      }}
                      className="border px-2 py-1.5 rounded-md text-sm h-8 bg-background focus:ring-ring focus:border-input"
                    >
                      {[50, 100, 200, 500, 1000].map(size => (
                        <option key={size} value={size}>{t("dashboard.pagination.pageSizeOption", { count: size })}</option>
                      ))}
                    </select>
                    <Button
                      variant="outline" size="sm"
                      onClick={() => setClientCurrentPage(p => Math.max(1, p - 1))}
                      disabled={clientCurrentPage === 1}
                      className="h-8 px-2"
                    >
                      {t("dashboard.pagination.previousButton", "Previous")}
                    </Button>
                    <Button
                      variant="outline" size="sm"
                      onClick={() => setClientCurrentPage(p => Math.min(clientTotalPages, p + 1))}
                      disabled={clientCurrentPage === clientTotalPages || clientTotalPages === 0}
                      className="h-8 px-2"
                    >
                      {t("dashboard.pagination.nextButton", "Next")}
                    </Button>
                  </div>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
