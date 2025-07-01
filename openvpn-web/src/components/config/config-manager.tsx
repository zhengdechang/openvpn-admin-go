"use client";

import React, { useState, useEffect } from "react";
import { ConfigItem } from "@/types/types";
import { serverAPI } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { toast } from "sonner";
import { Save, RefreshCw } from "lucide-react";
import ConfigItemComponent from "./config-item";

export default function ConfigManager() {
  const [configItems, setConfigItems] = useState<ConfigItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [editingItems, setEditingItems] = useState<Set<string>>(new Set());
  const [changedItems, setChangedItems] = useState<Record<string, any>>({});

  const fetchConfigItems = async () => {
    setLoading(true);
    try {
      const data = await serverAPI.getConfigItems();
      setConfigItems(data.items);
      setChangedItems({});
    } catch (error) {
      toast.error("获取配置项失败");
      console.error("Failed to fetch config items:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchConfigItems();
  }, []);

  const handleValueChange = (key: string, value: any) => {
    setChangedItems(prev => ({
      ...prev,
      [key]: value
    }));
    
    // 更新本地显示的值
    setConfigItems(prev => 
      prev.map(item => 
        item.key === key ? { ...item, value } : item
      )
    );
  };

  const handleEditToggle = (key: string) => {
    setEditingItems(prev => {
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
      toast.info("没有需要保存的更改");
      return;
    }

    setSaving(true);
    try {
      await serverAPI.updateConfigItems(changedItems);
      toast.success("配置保存成功");
      setChangedItems({});
      setEditingItems(new Set());
      // 重新获取配置以确保同步
      await fetchConfigItems();
    } catch (error) {
      toast.error("配置保存失败");
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
            加载配置项...
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle>服务器配置管理</CardTitle>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              onClick={handleRefresh}
              disabled={saving}
            >
              <RefreshCw className="h-4 w-4 mr-2" />
              刷新
            </Button>
            <Button
              onClick={handleSaveAll}
              disabled={!hasChanges || saving}
              className={hasChanges ? "bg-green-600 hover:bg-green-700" : ""}
            >
              <Save className="h-4 w-4 mr-2" />
              {saving ? "保存中..." : `保存所有更改${hasChanges ? ` (${Object.keys(changedItems).length})` : ""}`}
            </Button>
          </div>
        </div>
        {hasChanges && (
          <div className="text-sm text-orange-600 bg-orange-50 p-2 rounded">
            您有 {Object.keys(changedItems).length} 项未保存的更改
          </div>
        )}
      </CardHeader>
      <CardContent className="space-y-4">
        {configItems.map((item) => (
          <ConfigItemComponent
            key={item.key}
            item={item}
            onValueChange={handleValueChange}
            isEditing={editingItems.has(item.key)}
            onEditToggle={handleEditToggle}
          />
        ))}
      </CardContent>
    </Card>
  );
}
