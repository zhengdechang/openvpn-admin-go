"use client";

import React, { useState, useEffect } from "react";
import { ConfigItem } from "@/types/types";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { X, Plus, Edit2, Check, X as Cancel } from "lucide-react";

interface ConfigItemComponentProps {
  item: ConfigItem;
  onValueChange: (key: string, value: any) => void;
  isEditing: boolean;
  onEditToggle: (key: string) => void;
}

export default function ConfigItemComponent({
  item,
  onValueChange,
  isEditing,
  onEditToggle,
}: ConfigItemComponentProps) {
  const [localValue, setLocalValue] = useState(item.value);
  const [arrayItems, setArrayItems] = useState<string[]>(
    Array.isArray(item.value) ? item.value : []
  );

  useEffect(() => {
    setLocalValue(item.value);
    if (Array.isArray(item.value)) {
      setArrayItems(item.value);
    }
  }, [item.value]);

  const handleSave = () => {
    if (item.type === "array") {
      onValueChange(item.key, arrayItems);
    } else {
      onValueChange(item.key, localValue);
    }
    onEditToggle(item.key);
  };

  const handleCancel = () => {
    setLocalValue(item.value);
    if (Array.isArray(item.value)) {
      setArrayItems(item.value);
    }
    onEditToggle(item.key);
  };

  const addArrayItem = () => {
    setArrayItems([...arrayItems, ""]);
  };

  const removeArrayItem = (index: number) => {
    setArrayItems(arrayItems.filter((_, i) => i !== index));
  };

  const updateArrayItem = (index: number, value: string) => {
    const newItems = [...arrayItems];
    newItems[index] = value;
    setArrayItems(newItems);
  };

  const renderEditableValue = () => {
    if (!isEditing) {
      return renderDisplayValue();
    }

    switch (item.type) {
      case "text":
        return (
          <Input
            value={localValue || ""}
            onChange={(e) => setLocalValue(e.target.value)}
            placeholder={item.description}
            className="flex-1"
          />
        );
      case "number":
        return (
          <Input
            type="number"
            value={localValue || ""}
            onChange={(e) => setLocalValue(parseInt(e.target.value) || 0)}
            placeholder={item.description}
            className="flex-1"
          />
        );
      case "boolean":
        return (
          <Switch
            checked={localValue || false}
            onCheckedChange={setLocalValue}
          />
        );
      case "select":
        return (
          <Select value={localValue || ""} onValueChange={setLocalValue}>
            <SelectTrigger className="flex-1">
              <SelectValue placeholder={item.description} />
            </SelectTrigger>
            <SelectContent>
              {item.options?.map((option) => (
                <SelectItem key={option} value={option}>
                  {option}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        );
      case "array":
        return (
          <div className="flex-1 space-y-2">
            {arrayItems.map((arrayItem, index) => (
              <div key={index} className="flex items-center space-x-2">
                <Input
                  value={arrayItem}
                  onChange={(e) => updateArrayItem(index, e.target.value)}
                  placeholder="输入路由 (例如: 192.168.1.0 255.255.255.0)"
                  className="flex-1"
                />
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => removeArrayItem(index)}
                >
                  <X className="h-4 w-4" />
                </Button>
              </div>
            ))}
            <Button
              variant="outline"
              size="sm"
              onClick={addArrayItem}
              className="w-full"
            >
              <Plus className="h-4 w-4 mr-2" />
              添加路由
            </Button>
          </div>
        );
      default:
        return <span>不支持的类型</span>;
    }
  };

  const renderDisplayValue = () => {
    switch (item.type) {
      case "boolean":
        return (
          <Badge variant={item.value ? "default" : "secondary"}>
            {item.value ? "启用" : "禁用"}
          </Badge>
        );
      case "array":
        return (
          <div className="space-y-1">
            {Array.isArray(item.value) && item.value.length > 0 ? (
              item.value.map((arrayItem, index) => (
                <Badge key={index} variant="outline">
                  {arrayItem}
                </Badge>
              ))
            ) : (
              <span className="text-gray-500">无</span>
            )}
          </div>
        );
      default:
        return <span>{item.value || "未设置"}</span>;
    }
  };

  return (
    <div className="border rounded-lg p-4 space-y-3 hover:bg-gray-50 transition-colors">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="font-medium">{item.label}</h3>
          <p className="text-sm text-gray-600">{item.description}</p>
          {item.required && (
            <Badge variant="destructive" className="text-xs mt-1">
              必填
            </Badge>
          )}
        </div>
        <div className="flex items-center space-x-2">
          {isEditing ? (
            <>
              <Button variant="outline" size="sm" onClick={handleCancel}>
                <Cancel className="h-4 w-4" />
              </Button>
              <Button size="sm" onClick={handleSave}>
                <Check className="h-4 w-4" />
              </Button>
            </>
          ) : (
            <Button
              variant="outline"
              size="sm"
              onClick={() => onEditToggle(item.key)}
            >
              <Edit2 className="h-4 w-4" />
            </Button>
          )}
        </div>
      </div>
      <div
        className="flex items-center space-x-2 cursor-pointer"
        onDoubleClick={() => !isEditing && onEditToggle(item.key)}
        title={!isEditing ? "双击编辑" : ""}
      >
        {renderEditableValue()}
      </div>
    </div>
  );
}
