"use client";

import * as React from "react";
import MuiSelect from "@mui/material/Select";
import MenuItem from "@mui/material/MenuItem";

/**
 * 兼容旧 Radix Select 复合 API 的 MUI 实现。
 * 用法保持不变：
 *   <Select value onValueChange>
 *     <SelectTrigger className><SelectValue placeholder/></SelectTrigger>
 *     <SelectContent>
 *       <SelectItem value>label</SelectItem>
 *     </SelectContent>
 *   </Select>
 * SelectTrigger / SelectValue / SelectContent / SelectItem 仅作为标记节点，
 * 由 Select 解析出 placeholder、className 与选项后渲染原生 MUI Select。
 */

interface SelectItemProps {
  value: string;
  children?: React.ReactNode;
  disabled?: boolean;
}
export const SelectItem: React.FC<SelectItemProps> = () => null;

interface SelectValueProps {
  placeholder?: string;
}
export const SelectValue: React.FC<SelectValueProps> = () => null;

interface SelectTriggerProps {
  className?: string;
  children?: React.ReactNode;
}
export const SelectTrigger: React.FC<SelectTriggerProps> = () => null;

interface SelectContentProps {
  children?: React.ReactNode;
}
export const SelectContent: React.FC<SelectContentProps> = () => null;

export const SelectGroup: React.FC<{ children?: React.ReactNode }> = ({ children }) => <>{children}</>;
export const SelectLabel: React.FC<{ children?: React.ReactNode }> = ({ children }) => <>{children}</>;
export const SelectSeparator: React.FC = () => null;

interface SelectProps {
  value?: string;
  defaultValue?: string;
  onValueChange?: (value: string) => void;
  disabled?: boolean;
  children?: React.ReactNode;
}

export function Select({ value, defaultValue, onValueChange, disabled, children }: SelectProps) {
  let placeholder: string | undefined;
  let triggerClassName: string | undefined;
  const items: { value: string; label: React.ReactNode; disabled?: boolean }[] = [];

  React.Children.forEach(children, (child) => {
    if (!React.isValidElement(child)) return;
    if (child.type === SelectTrigger) {
      const p = child.props as SelectTriggerProps;
      triggerClassName = p.className;
      React.Children.forEach(p.children, (c) => {
        if (React.isValidElement(c) && c.type === SelectValue) {
          placeholder = (c.props as SelectValueProps).placeholder;
        }
      });
    }
    if (child.type === SelectContent) {
      React.Children.forEach((child.props as SelectContentProps).children, (item) => {
        if (React.isValidElement(item) && item.type === SelectItem) {
          const ip = item.props as SelectItemProps;
          items.push({ value: ip.value, label: ip.children, disabled: ip.disabled });
        }
      });
    }
  });

  return (
    <MuiSelect
      className={triggerClassName}
      size="small"
      displayEmpty
      disabled={disabled}
      value={value ?? defaultValue ?? ""}
      onChange={(e) => onValueChange?.(e.target.value as string)}
      renderValue={(selected) => {
        const v = selected as string;
        if (!v) {
          return <span style={{ opacity: 0.6 }}>{placeholder}</span>;
        }
        return items.find((it) => it.value === v)?.label ?? v;
      }}
    >
      {placeholder ? (
        <MenuItem value="" disabled>
          {placeholder}
        </MenuItem>
      ) : null}
      {items.map((it) => (
        <MenuItem key={it.value} value={it.value} disabled={it.disabled}>
          {it.label}
        </MenuItem>
      ))}
    </MuiSelect>
  );
}
