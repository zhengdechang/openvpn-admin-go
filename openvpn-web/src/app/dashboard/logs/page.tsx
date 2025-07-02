/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-07-01 14:40:50
 */
// In openvpn-web/src/app/dashboard/logs/page.tsx
"use client";

import React, { useEffect, useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/services/api";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import SimpleLogViewer from "@/components/ui/simple-log-viewer";

export default function LogsPage() {
  const { t } = useTranslation();
  const [serverLogs, setServerLogs] = useState<string>("");
  const [loadingServer, setLoadingServer] = useState(true);

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
          <SimpleLogViewer height={400} />
        </CardContent>
      </Card>
    </MainLayout>
  );
}
