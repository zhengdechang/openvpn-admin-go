"use client";

import React, { useState, useEffect } from "react";
import { ConfigItem } from "@/types/types";
import { serverAPI } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import { Save, RefreshCw } from "lucide-react";
import { useTranslation } from "react-i18next";
import ConfigItemComponent from "./config-item";

export default function ConfigManager() {
  const { t, i18n } = useTranslation();
  const [configItems, setConfigItems] = useState<ConfigItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [editingItems, setEditingItems] = useState<Set<string>>(new Set());
  const [changedItems, setChangedItems] = useState<Record<string, any>>({});

  const fetchConfigItems = async () => {
    setLoading(true);
    try {
      const data = await serverAPI.getConfigItems(i18n.language);
      setConfigItems(data.items);
      setChangedItems({});
    } catch (error) {
      toast.error(t("dashboard.server.config.fetchError"));
      console.error("Failed to fetch config items:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfigItems();
  }, []);

  // 当语言变化时重新获取配置项
  useEffect(() => {
    if (configItems.length > 0) {
      fetchConfigItems();
    }
  }, [i18n.language]);

  const handleValueChange = (key: string, value: any) => {
    setChangedItems((prev) => ({
      ...prev,
      [key]: value,
    }));

    // 更新本地显示的值
    setConfigItems((prev) =>
      prev.map((item) => (item.key === key ? { ...item, value } : item))
    );
  };

  const handleEditToggle = (key: string) => {
    setEditingItems((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(key)) {
        newSet.delete(key);
      } else {
        newSet.add(key);
      }
      return newSet;
    });
  };

  const handleSaveAll = async () => {
    if (Object.keys(changedItems).length === 0) {
      toast.info(t("dashboard.server.config.noChanges"));
      return;
    }

    setSaving(true);
    try {
      await serverAPI.updateConfigItems(changedItems);
      toast.success(t("dashboard.server.config.saveSuccess"));
      setChangedItems({});
      setEditingItems(new Set());
      // 重新获取配置以确保同步
      await fetchConfigItems();
    } catch (error) {
      toast.error(t("dashboard.server.config.saveError"));
      console.error("Failed to save config items:", error);
    } finally {
      setSaving(false);
    }
  };

  const handleRefresh = () => {
    fetchConfigItems();
    setEditingItems(new Set());
    setChangedItems({});
  };

  const hasChanges = Object.keys(changedItems).length > 0;

  if (loading) {
    return (
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-center">
            <RefreshCw className="h-6 w-6 animate-spin mr-2" />
            {t("dashboard.server.config.loading")}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>{t("dashboard.server.config.title")}</CardTitle>
          <div className="flex items-center space-x-2">
            <Button variant="outline" onClick={handleRefresh} disabled={saving}>
              <RefreshCw className="h-4 w-4 mr-2" />
              {t("dashboard.server.config.refreshButton")}
            </Button>
            <Button
              onClick={handleSaveAll}
              disabled={!hasChanges || saving}
              className={hasChanges ? "bg-green-600 hover:bg-green-700" : ""}
            >
              <Save className="h-4 w-4 mr-2" />
              {saving
                ? t("dashboard.server.config.saving")
                : `${t("dashboard.server.config.saveAllButton")}${
                    hasChanges ? ` (${Object.keys(changedItems).length})` : ""
                  }`}
            </Button>
          </div>
        </div>
        {hasChanges && (
          <div className="text-sm text-orange-600 bg-orange-50 p-2 rounded">
            {t("dashboard.server.config.unsavedChanges", {
              count: Object.keys(changedItems).length,
            })}
          </div>
        )}
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {configItems.map((item) => (
            <ConfigItemComponent
              key={item.key}
              item={item}
              onValueChange={handleValueChange}
              isEditing={editingItems.has(item.key)}
              onEditToggle={handleEditToggle}
            />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
