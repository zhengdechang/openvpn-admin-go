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
import { ClientLog, PaginatedClientLogs, LiveClientConnection } from "@/types/types";
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

  const [clientApiLogs, setClientApiLogs] = useState<ClientLog[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [totalLogs, setTotalLogs] = useState(0);
  const [loadingClientLogs, setLoadingClientLogs] = useState(true);
  const [filterUsername, setFilterUsername] = useState<string>("");
  const [searchInputUsername, setSearchInputUsername] = useState<string>("");

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

  // Fetch client API logs (existing useEffect)
  useEffect(() => {
    const fetchClientApiLogs = async () => {
      setLoadingClientLogs(true);
      try {
        const response = await openvpnAPI.getClientLogs(
          currentPage,
          pageSize,
          filterUsername.trim() === "" ? undefined : filterUsername.trim()
        );
        if (response.success && response.data) {
          setClientApiLogs(response.data.data || []);
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
  }, [currentPage, pageSize, filterUsername, t]);

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


  const handleSearchUsername = () => {
    setCurrentPage(1);
    setFilterUsername(searchInputUsername);
  };

  const totalPages = Math.ceil(totalLogs / pageSize);

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

      {/* Client API Logs Card (existing) */}
      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.logs.clientConnectionLogsTitle")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* ... existing content for client API logs ... */}
           <div className="flex items-center space-x-2 mb-4">
             <Input
               placeholder={t("dashboard.logs.filterByUsernamePlaceholder")}
               value={searchInputUsername}
               onChange={(e) => setSearchInputUsername(e.target.value)}
               onKeyPress={(e) => e.key === 'Enter' && handleSearchUsername()}
             />
             <Button onClick={handleSearchUsername}>
               {t("dashboard.logs.searchButton")}
             </Button>
             <select
               value={pageSize}
               onChange={(e) => {
                 setPageSize(Number(e.target.value));
                 setCurrentPage(1);
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
