"use client";

import React, { useState, useEffect } from "react";
import { ConfigItem } from "@/types/types";
import { X, Plus, Edit2, Check, X as Cancel } from "lucide-react";
import { useTranslation } from "react-i18next";
import MuiButton from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Chip from "@mui/material/Chip";
import Switch from "@mui/material/Switch";
import FormControlLabel from "@mui/material/FormControlLabel";
import { FormControl, InputLabel, Select as MuiSelect, MenuItem } from "@mui/material";

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
  const { t } = useTranslation();
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
          <TextField
            value={localValue || ""}
            onChange={(e) => setLocalValue(e.target.value)}
            placeholder={item.description}
            fullWidth
            size="small"
          />
        );
      case "number":
        return (
          <TextField
            type="number"
            value={localValue || ""}
            onChange={(e) => setLocalValue(parseInt(e.target.value) || 0)}
            placeholder={item.description}
            fullWidth
            size="small"
          />
        );
      case "boolean":
        return (
          <FormControlLabel
            control={
              <Switch
                checked={localValue || false}
                onChange={(e) => setLocalValue(e.target.checked)}
              />
            }
            label={localValue ? "启用" : "禁用"}
          />
        );
      case "select":
        return (
          <FormControl fullWidth size="small">
            <InputLabel>{item.description}</InputLabel>
            <MuiSelect
              value={localValue || ""}
              label={item.description}
              onChange={(e) => setLocalValue(e.target.value)}
            >
              {item.options?.map((option) => (
                <MenuItem key={option} value={option}>
                  {option}
                </MenuItem>
              ))}
            </MuiSelect>
          </FormControl>
        );
      case "array":
        return (
          <div className="w-full h-full flex flex-col">
            <div className="flex-1 overflow-y-auto space-y-1 max-h-16">
              {arrayItems.map((arrayItem, index) => (
                <div key={index} className="flex items-center space-x-1">
                  <TextField
                    value={arrayItem}
                    onChange={(e) => updateArrayItem(index, e.target.value)}
                    placeholder="输入路由"
                    size="small"
                    sx={{ flex: 1, "& .MuiInputBase-input": { py: 0.5, fontSize: "0.75rem" } }}
                  />
                  <MuiButton
                    variant="outlined"
                    size="small"
                    onClick={() => removeArrayItem(index)}
                    sx={{ minWidth: "28px", p: "2px", height: "28px" }}
                  >
                    <X className="h-3 w-3" />
                  </MuiButton>
                </div>
              ))}
            </div>
            <MuiButton
              variant="outlined"
              size="small"
              onClick={addArrayItem}
              fullWidth
              startIcon={<Plus className="h-3 w-3" />}
              sx={{ mt: 0.5, fontSize: "0.75rem", height: "28px" }}
            >
              {t("dashboard.server.config.addButton")}
            </MuiButton>
          </div>
        );
      default:
        return <span>{t("dashboard.server.config.unsupportedType")}</span>;
    }
  };

  const renderDisplayValue = () => {
    switch (item.type) {
      case "boolean":
        return (
          <Chip
            label={item.value ? "启用" : "禁用"}
            color={item.value ? "success" : "default"}
            size="small"
          />
        );
      case "array":
        return (
          <div className="flex flex-wrap gap-1 max-h-16 overflow-y-auto">
            {Array.isArray(item.value) && item.value.length > 0 ? (
              item.value.map((arrayItem, index) => (
                <Chip key={index} label={arrayItem} variant="outlined" size="small" />
              ))
            ) : (
              <span className="text-gray-500 text-sm">无</span>
            )}
          </div>
        );
      default:
        return <span className="text-sm">{item.value || "未设置"}</span>;
    }
  };

  return (
    <div className="border rounded-lg p-4 hover:bg-gray-50 transition-colors h-40 flex flex-col">
      <div className="flex items-start justify-between mb-3">
        <div className="flex-1 min-w-0">
          <h3 className="font-medium text-sm truncate">{item.label}</h3>
          <p className="text-xs text-gray-600 line-clamp-2">
            {item.description}
          </p>
          {item.required && (
            <Chip label="必填" color="error" size="small" sx={{ mt: 0.5 }} />
          )}
        </div>
        <div className="flex items-center space-x-1 ml-2 flex-shrink-0">
          {isEditing ? (
            <>
              <MuiButton variant="outlined" size="small" onClick={handleCancel} sx={{ minWidth: "28px", p: "2px" }}>
                <Cancel className="h-3 w-3" />
              </MuiButton>
              <MuiButton variant="contained" size="small" onClick={handleSave} sx={{ minWidth: "28px", p: "2px" }}>
                <Check className="h-3 w-3" />
              </MuiButton>
            </>
          ) : (
            <MuiButton
              variant="outlined"
              size="small"
              onClick={() => onEditToggle(item.key)}
              sx={{ minWidth: "28px", p: "2px" }}
            >
              <Edit2 className="h-3 w-3" />
            </MuiButton>
          )}
        </div>
      </div>
      <div
        className="flex-1 w-full cursor-pointer flex items-center"
        onDoubleClick={() => !isEditing && onEditToggle(item.key)}
        title={!isEditing ? t("dashboard.server.config.doubleClickToEdit") : ""}
      >
        {renderEditableValue()}
      </div>
    </div>
  );
}
