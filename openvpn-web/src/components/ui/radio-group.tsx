"use client";

import * as React from "react";
import MuiRadioGroup from "@mui/material/RadioGroup";
import Radio from "@mui/material/Radio";
import FormControlLabel from "@mui/material/FormControlLabel";

/**
 * 兼容旧 Radix API（value / onValueChange）的 MUI RadioGroup。
 */
interface RadioGroupProps {
  value?: string;
  defaultValue?: string;
  onValueChange?: (value: string) => void;
  className?: string;
  children?: React.ReactNode;
}

const RadioGroup = React.forwardRef<HTMLDivElement, RadioGroupProps>(
  ({ value, defaultValue, onValueChange, className, children }, ref) => (
    <MuiRadioGroup
      ref={ref}
      className={className}
      value={value ?? defaultValue ?? ""}
      onChange={(_e, v) => onValueChange?.(v)}
    >
      {children}
    </MuiRadioGroup>
  )
);
RadioGroup.displayName = "RadioGroup";

interface RadioGroupItemProps {
  value: string;
  label?: React.ReactNode;
  disabled?: boolean;
  id?: string;
  className?: string;
}

const RadioGroupItem = ({ value, label, disabled, id, className }: RadioGroupItemProps) => (
  <FormControlLabel
    value={value}
    disabled={disabled}
    className={className}
    control={<Radio id={id} size="small" />}
    label={label ?? ""}
  />
);
RadioGroupItem.displayName = "RadioGroupItem";

export { RadioGroup, RadioGroupItem };
