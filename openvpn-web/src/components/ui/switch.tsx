"use client";

import * as React from "react";
import MuiSwitch from "@mui/material/Switch";

export interface SwitchProps {
  checked?: boolean;
  defaultChecked?: boolean;
  disabled?: boolean;
  className?: string;
  onCheckedChange?: (checked: boolean) => void;
  name?: string;
  id?: string;
}

/**
 * 兼容旧 Radix API（checked / onCheckedChange）的 MUI Switch。
 */
const Switch = React.forwardRef<HTMLButtonElement, SwitchProps>(
  ({ checked, defaultChecked, disabled, className, onCheckedChange, name, id }, ref) => (
    <MuiSwitch
      ref={ref as React.Ref<HTMLButtonElement>}
      checked={checked}
      defaultChecked={defaultChecked}
      disabled={disabled}
      className={className}
      name={name}
      id={id}
      onChange={(_e, value) => onCheckedChange?.(value)}
    />
  )
);
Switch.displayName = "Switch";

export { Switch };
