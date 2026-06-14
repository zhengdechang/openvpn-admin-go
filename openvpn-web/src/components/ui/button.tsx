"use client";

import * as React from "react";
import MuiButton from "@mui/material/Button";
import type { SxProps, Theme } from "@mui/material/styles";
import { cn } from "@/lib/utils";

type Variant =
  | "default"
  | "destructive"
  | "outline"
  | "secondary"
  | "ghost"
  | "link";
type Size = "default" | "sm" | "lg" | "icon";

export interface ButtonProps
  extends Omit<React.ButtonHTMLAttributes<HTMLButtonElement>, "color"> {
  variant?: Variant;
  size?: Size;
  asChild?: boolean;
}

function muiVariant(variant: Variant): "contained" | "outlined" | "text" {
  if (variant === "outline") return "outlined";
  if (variant === "ghost" || variant === "link") return "text";
  return "contained";
}

function muiColor(
  variant: Variant
): "primary" | "secondary" | "error" | "inherit" {
  switch (variant) {
    case "destructive":
      return "error";
    case "secondary":
      return "secondary";
    case "outline":
    case "ghost":
      return "inherit";
    default:
      return "primary";
  }
}

function muiSize(size: Size): "small" | "medium" | "large" {
  if (size === "sm" || size === "icon") return "small";
  if (size === "lg") return "large";
  return "medium";
}

function extraSx(variant: Variant, size: Size): SxProps<Theme> {
  const sx: Record<string, unknown> = {};
  if (variant === "outline") {
    sx.borderColor = "var(--argon-border)";
    sx.color = "var(--argon-gray-dark)";
    sx["&:hover"] = {
      borderColor: "var(--argon-primary)",
      backgroundColor: "rgba(94,114,228,0.04)",
    };
  }
  if (variant === "ghost") {
    sx.color = "var(--argon-gray-dark)";
  }
  if (variant === "link") {
    sx.textDecoration = "underline";
    sx.padding = 0;
    sx.minWidth = 0;
  }
  if (size === "icon") {
    sx.minWidth = 0;
    sx.width = 40;
    sx.height = 40;
    sx.padding = 0;
  }
  return sx as SxProps<Theme>;
}

/**
 * 兼容旧 shadcn API（variant/size/asChild/className）的 MUI Button。
 */
const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "default", size = "default", asChild = false, children, ...props }, ref) => {
    const common = {
      ref,
      className,
      variant: muiVariant(variant),
      color: muiColor(variant),
      size: muiSize(size),
      sx: extraSx(variant, size),
      disableElevation: variant === "outline" || variant === "ghost" || variant === "link",
    } as const;

    // asChild：把子元素作为渲染组件（如 <Button asChild><Link/></Button>）
    if (asChild && React.isValidElement(children)) {
      const child = children as React.ReactElement<Record<string, unknown>>;
      const childProps = child.props as Record<string, unknown>;
      return (
        <MuiButton
          {...common}
          component={child.type as React.ElementType}
          {...childProps}
          className={cn(className, childProps.className as string | undefined)}
        >
          {childProps.children as React.ReactNode}
        </MuiButton>
      );
    }

    return (
      <MuiButton {...common} {...props}>
        {children}
      </MuiButton>
    );
  }
);
Button.displayName = "Button";

export { Button };
