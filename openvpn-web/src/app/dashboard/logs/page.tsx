"use client";

import React, { useEffect, useState } from "react";
import MainLayout from "@/components/layout/main-layout";
import { openvpnAPI } from "@/services/api";
import { useTranslation } from "react-i18next";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

export default function LogsPage() {
  const { t } = useTranslation();
  const [serverLogs, setServerLogs] = useState<string>("");
  const [clientUsername, setClientUsername] = useState<string>("");
  const [clientLogs, setClientLogs] = useState<string[] | null>(null);
  const [loadingServer, setLoadingServer] = useState(true);

  useEffect(() => {
    const fetchServerLogs = async () => {
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
  }, []);

  const handleFetchClientLogs = async () => {
    if (!clientUsername) {
      toast.error(t("dashboard.logs.usernameRequired"));
      return;
    }
    try {
      const logs = await openvpnAPI.getClientLogs(clientUsername);
      setClientLogs(logs);
    } catch (error) {
      toast.error(t("dashboard.logs.fetchClientLogsError"));
    }
  };

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

      <Card>
        <CardHeader>
          <CardTitle>{t("dashboard.logs.clientLogsCardTitle")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center space-x-2">
            <Input
              placeholder={t("dashboard.logs.usernamePlaceholder")}
              value={clientUsername}
              onChange={(e) => setClientUsername(e.target.value)}
            />
            <Button onClick={handleFetchClientLogs}>
              {t("dashboard.logs.queryButton")}
            </Button>
          </div>
          {clientLogs && (
            <pre className="whitespace-pre-wrap">
              {clientLogs.join("\n") || t("dashboard.logs.noClientLogs")}
            </pre>
          )}
        </CardContent>
      </Card>
    </MainLayout>
  );
}
