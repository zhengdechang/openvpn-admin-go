"use client";

import { createTheme } from "@mui/material/styles";

// LuCI Argon 调色板（与 luci-theme-argon 一致）
const ARGON = {
  primary: "#5e72e4",
  primaryDark: "#483d8b",
  default: "#172b4d",
  gray: "#8898aa",
  grayDark: "#32325d",
  bg: "#f4f5f7",
  border: "#e9ecef",
  success: "#2dce89",
  info: "#11cdef",
  warning: "#fb6340",
  danger: "#f5365c",
};

const theme = createTheme({
  cssVariables: true,
  palette: {
    mode: "light",
    primary: { main: ARGON.primary, dark: ARGON.primaryDark, contrastText: "#fff" },
    secondary: { main: ARGON.grayDark, contrastText: "#fff" },
    success: { main: ARGON.success, contrastText: "#fff" },
    info: { main: ARGON.info, contrastText: "#fff" },
    warning: { main: ARGON.warning, contrastText: "#fff" },
    error: { main: ARGON.danger, contrastText: "#fff" },
    text: { primary: ARGON.default, secondary: ARGON.gray },
    background: { default: ARGON.bg, paper: "#ffffff" },
    divider: ARGON.border,
  },
  shape: { borderRadius: 8 },
  typography: {
    fontFamily:
      '"Google Sans", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Microsoft Yahei", "WenQuanYi Micro Hei", "Helvetica Neue", Helvetica, sans-serif',
    button: { textTransform: "none", fontWeight: 600 },
  },
  components: {
    // Argon 按钮：圆角 + 柔和阴影，不大写
    MuiButton: {
      styleOverrides: {
        root: { borderRadius: 6, textTransform: "none" },
      },
      variants: [
        {
          props: { variant: "contained", color: "primary" },
          style: {
            boxShadow: "0 4px 10px rgba(94,114,228,0.25)",
            "&:hover": { boxShadow: "0 6px 14px rgba(94,114,228,0.35)" },
          },
        },
      ],
    },
    // Argon 面板：白底、圆角、柔和阴影
    MuiPaper: {
      styleOverrides: {
        rounded: { borderRadius: 8 },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          boxShadow: "0 0 1rem 0 rgba(136, 152, 170, 0.15)",
          border: `1px solid ${ARGON.border}`,
        },
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          backgroundColor: "#fff",
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        head: {
          backgroundColor: "#f6f9fc",
          color: ARGON.grayDark,
          fontWeight: 600,
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: { borderRadius: 10 },
      },
    },
  },
});

export default theme;
