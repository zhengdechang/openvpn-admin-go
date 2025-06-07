"use client";

import React, { useEffect, useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/services/api";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { ClientLog, PaginatedClientLogs } from "@/types/types"; // Import new types
import {
  Table,
  TableHeader,
  TableRow,
  TableHead,
  TableBody,
  TableCell,
} from "@/components/ui/table"; // Import table components

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
  const [loadingServer, setLoadingServer] = useState(true);

  // New state for client connection logs
  const [clientApiLogs, setClientApiLogs] = useState<ClientLog[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10); // Default page size
  const [totalLogs, setTotalLogs] = useState(0);
  const [loadingClientLogs, setLoadingClientLogs] = useState(true);
  const [filterUserId, setFilterUserId] = useState<string>("");
  const [searchInputUserId, setSearchInputUserId] = useState<string>("");

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

  // Fetch client connection logs
  useEffect(() => {
    const fetchClientApiLogs = async () => {
      setLoadingClientLogs(true);
      try {
        const response = await openvpnAPI.getClientLogs(
          currentPage,
          pageSize,
          filterUserId.trim() === "" ? undefined : filterUserId.trim()
        );
        if (response.success && response.data) {
          setClientApiLogs(response.data.data || []); // Ensure we always set an array
          setTotalLogs(response.data.total || 0);
        } else {
          toast.error(response.error || t("dashboard.logs.fetchClientLogsError"));
          setClientApiLogs([]);
          setTotalLogs(0);
        }
      } catch (error) {
        toast.error(t("dashboard.logs.fetchClientLogsError"));
        setClientApiLogs([]);
        setTotalLogs(0);
      } finally {
        setLoadingClientLogs(false);
      }
    };
    fetchClientApiLogs();
  }, [currentPage, pageSize, filterUserId, t]);

  const handleSearchUserId = () => {
    setCurrentPage(1); // Reset to first page on new search
    setFilterUserId(searchInputUserId);
  };

  const totalPages = Math.ceil(totalLogs / pageSize);

  return (
    <MainLayout className="p-4 space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.logs.serverLogsCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingServer ? (
            <p>{t("common.loading")}</p>
          ) : (
            <pre className="whitespace-pre-wrap">
              {serverLogs || t("dashboard.logs.noServerLogs")}
            </pre>
          )}
        </CardContent>
      </Card>

      {/* New Client Connection Logs Card */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.logs.clientConnectionLogsTitle")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center space-x-2 mb-4">
            <Input
              placeholder={t("dashboard.logs.filterByUserIdPlaceholder")}
              value={searchInputUserId}
              onChange={(e) => setSearchInputUserId(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSearchUserId()}
            />
            <Button onClick={handleSearchUserId}>
              {t("dashboard.logs.searchButton")}
            </Button>
            <select
              value={pageSize}
              onChange={(e) => {
                setPageSize(Number(e.target.value));
                setCurrentPage(1); // Reset to page 1 on page size change
              }}
              className="border px-2 py-2 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            >
              {[10, 25, 50, 100].map(size => (
                <option key={size} value={size}>
                  {t("dashboard.logs.showEntries", { count: size })}
                </option>
              ))}
            </select>
          </div>

          {loadingClientLogs ? (
            <p>{t("common.loading")}</p>
          ) : !clientApiLogs || clientApiLogs.length === 0 ? (
            <p>{t("dashboard.logs.noClientConnectionLogs")}</p>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t("dashboard.logs.columnUserId")}</TableHead>
                    <TableHead>{t("dashboard.logs.columnOnlineStatus")}</TableHead>
                    <TableHead>{t("dashboard.logs.columnOnlineDuration")}</TableHead>
                    <TableHead>{t("dashboard.logs.columnTrafficUsage")}</TableHead>
                    <TableHead>{t("dashboard.logs.columnLastConnectionTime")}</TableHead>
                    <TableHead>{t("dashboard.logs.columnLogCreatedAt")}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {clientApiLogs.map((log) => (
                    <TableRow key={log.id}>
                      <TableCell>{log.userId}</TableCell>
                      <TableCell>
                        {log.isOnline
                          ? t("dashboard.logs.statusOnline")
                          : t("dashboard.logs.statusOffline")}
                      </TableCell>
                      <TableCell>{formatDuration(log.onlineDuration)}</TableCell>
                      <TableCell>{formatBytes(log.trafficUsage)}</TableCell>
                      <TableCell>
                        {log.lastConnectionTime
                          ? new Date(log.lastConnectionTime).toLocaleString()
                          : t("common.na")}
                      </TableCell>
                      <TableCell>{new Date(log.createdAt).toLocaleString()}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              <div className="flex justify-between items-center mt-4">
                <div>
                  <Button
                    onClick={() => setCurrentPage((prev) => Math.max(1, prev - 1))}
                    disabled={currentPage === 1 || loadingClientLogs}
                  >
                    {t("common.previous")}
                  </Button>
                  <span className="mx-2">
                    {t("common.pageInfo", { currentPage, totalPages })}
                  </span>
                  <Button
                    onClick={() => setCurrentPage((prev) => Math.min(totalPages, prev + 1))}
                    disabled={currentPage === totalPages || totalPages === 0 || loadingClientLogs}
                  >
                    {t("common.next")}
                  </Button>
                </div>
                <div className="text-sm text-gray-600">
                  {t("common.totalItems", { count: totalLogs })}
                </div>
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
