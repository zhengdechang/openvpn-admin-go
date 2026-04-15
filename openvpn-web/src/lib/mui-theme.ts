import { createTheme } from "@mui/material/styles";

export const muiTheme = createTheme({
  palette: {
    primary: { main: "#0369A1" },
    secondary: { main: "#22C55E" },
    background: { default: "#F0F9FF" },
    error: { main: "#EF4444" },
  },
  shape: { borderRadius: 8 },
  components: {
    MuiButton: {
      defaultProps: { disableElevation: true },
      styleOverrides: {
        root: { textTransform: "none", fontWeight: 600 },
      },
    },
    MuiTextField: {
      defaultProps: { size: "small" },
    },
    MuiSelect: {
      defaultProps: { size: "small" },
    },
    MuiDialog: {
      defaultProps: { maxWidth: "sm", fullWidth: true },
    },
  },
});
