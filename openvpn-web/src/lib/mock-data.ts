import { OpenVPNClient } from "../types/types";

export const mockOpenVPNClients: OpenVPNClient[] = [
  {
    id: "1",
    name: "John Doe",
    email: "john.doe@example.com",
    status: "active",
    createdAt: "2023-01-15T08:30:00Z",
    lastConnected: "2023-04-20T14:25:00Z",
    ipAddress: "192.168.1.100",
    notes: "Regular user with full access",
  },
  {
    id: "2",
    name: "Jane Smith",
    email: "jane.smith@example.com",
    status: "inactive",
    createdAt: "2023-02-10T11:45:00Z",
    lastConnected: "2023-03-15T09:30:00Z",
    ipAddress: "192.168.1.101",
    notes: "Temporary access granted",
  },
];
