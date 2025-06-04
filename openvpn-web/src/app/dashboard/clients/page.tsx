"use client";

import React, { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { useRouter } from "next/navigation";
import MainLayout from "@/components/layout/main-layout";
import { OpenVPNClient } from "@/types/types";
import { openvpnAPI } from "@/services/api";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function ClientsPage() {
  const router = useRouter();
  const { t } = useTranslation();
  const [clients, setClients] = useState<OpenVPNClient[]>([]);
  const [filteredClients, setFilteredClients] = useState<OpenVPNClient[]>([]);
  const [searchTerm, setSearchTerm] = useState("");
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [currentClient, setCurrentClient] = useState<OpenVPNClient | null>(
    null
  );
  const [formData, setFormData] = useState<Partial<OpenVPNClient>>({
    name: "",
    email: "",
    status: "pending",
    notes: "",
  });

  // 加载客户端数据
  useEffect(() => {
    const fetchClients = async () => {
      try {
        const data = await openvpnAPI.getClientList();
        setClients(data);
        setFilteredClients(data);
      } catch (error) {
        toast.error("加载客户端列表失败");
      }
    };
    fetchClients();
  }, []);

  // 搜索过滤
  useEffect(() => {
    if (searchTerm.trim() === "") {
      setFilteredClients(clients);
    } else {
      const filtered = clients.filter(
        (client) =>
          client.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          client.email.toLowerCase().includes(searchTerm.toLowerCase())
      );
      setFilteredClients(filtered);
    }
  }, [searchTerm, clients]);

  // 处理表单输入变化
  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  // 处理状态选择变化
  const handleStatusChange = (value: string) => {
    setFormData((prev) => ({
      ...prev,
      status: value as OpenVPNClient["status"],
    }));
  };

  // 打开添加客户端对话框
  const handleAddClick = () => {
    setFormData({
      name: "",
      email: "",
      status: "pending",
      notes: "",
    });
    setIsAddDialogOpen(true);
  };

  // 打开编辑客户端对话框
  const handleEditClick = (client: OpenVPNClient) => {
    setCurrentClient(client);
    setFormData({
      name: client.name,
      email: client.email,
      status: client.status,
      notes: client.notes || "",
    });
    setIsEditDialogOpen(true);
  };

  // 打开删除客户端对话框
  const handleDeleteClick = (client: OpenVPNClient) => {
    setCurrentClient(client);
    setIsDeleteDialogOpen(true);
  };

  // 处理添加客户端
  const handleAddClient = () => {
    if (!formData.name || !formData.email) {
      toast.error(t("dashboard.error.add"));
      return;
    }

    const newClient: OpenVPNClient = {
      id: Date.now().toString(),
      name: formData.name,
      email: formData.email,
      status: formData.status as OpenVPNClient["status"],
      createdAt: new Date().toISOString(),
      notes: formData.notes,
    };

    setClients((prev) => [...prev, newClient]);
    setFilteredClients((prev) => [...prev, newClient]);
    setIsAddDialogOpen(false);
    toast.success(t("dashboard.success.added"));
  };

  // 处理编辑客户端
  const handleEditClient = () => {
    if (!currentClient || !formData.name || !formData.email) {
      toast.error(t("dashboard.error.update"));
      return;
    }

    const updatedClients = clients.map((client) =>
      client.id === currentClient.id
        ? {
            ...client,
            name: formData.name,
            email: formData.email,
            status: formData.status as OpenVPNClient["status"],
            notes: formData.notes,
          }
        : client
    );

    setClients(updatedClients);
    setFilteredClients(updatedClients);
    setIsEditDialogOpen(false);
    toast.success(t("dashboard.success.updated"));
  };

  // 处理删除客户端
  const handleDeleteClient = () => {
    if (!currentClient) {
      toast.error(t("dashboard.error.delete"));
      return;
    }

    const updatedClients = clients.filter(
      (client) => client.id !== currentClient.id
    );

    setClients(updatedClients);
    setFilteredClients(updatedClients);
    setIsDeleteDialogOpen(false);
    toast.success(t("dashboard.success.deleted"));
  };

  // 格式化日期
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  // 获取状态标签样式
  const getStatusStyle = (status: OpenVPNClient["status"]) => {
    switch (status) {
      case "active":
        return "bg-green-100 text-green-800";
      case "inactive":
        return "bg-red-100 text-red-800";
      case "pending":
        return "bg-yellow-100 text-yellow-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <MainLayout>
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">{t("dashboard.title")}</h1>
          <Button onClick={handleAddClick}>{t("dashboard.addClient")}</Button>
        </div>

        <Card className="mb-8">
          <CardContent className="pt-6">
            <div className="flex items-center mb-4">
              <Input
                placeholder={t("dashboard.search")}
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="max-w-md"
              />
            </div>

            {filteredClients.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                {t("dashboard.noClients")}
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{t("dashboard.table.name")}</TableHead>
                    <TableHead>{t("dashboard.table.email")}</TableHead>
                    <TableHead>{t("dashboard.table.status")}</TableHead>
                    <TableHead>{t("dashboard.table.createdAt")}</TableHead>
                    <TableHead>{t("dashboard.table.lastConnected")}</TableHead>
                    <TableHead>{t("dashboard.table.ipAddress")}</TableHead>
                    <TableHead>{t("dashboard.table.actions")}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredClients.map((client) => (
                    <TableRow key={client.id}>
                      <TableCell className="font-medium">
                        {client.name}
                      </TableCell>
                      <TableCell>{client.email}</TableCell>
                      <TableCell>
                        <span
                          className={`px-2 py-1 rounded-full text-xs font-medium ${getStatusStyle(
                            client.status
                          )}`}
                        >
                          {t(`dashboard.status.${client.status}`)}
                        </span>
                      </TableCell>
                      <TableCell>{formatDate(client.createdAt)}</TableCell>
                      <TableCell>
                        {client.lastConnected
                          ? formatDate(client.lastConnected)
                          : "-"}
                      </TableCell>
                      <TableCell>{client.ipAddress || "-"}</TableCell>
                      <TableCell>
                        <div className="flex space-x-2">
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleEditClick(client)}
                          >
                            {t("dashboard.actions.edit")}
                          </Button>
                          <Button
                            variant="destructive"
                            size="sm"
                            onClick={() => handleDeleteClick(client)}
                          >
                            {t("dashboard.actions.delete")}
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </CardContent>
        </Card>

        {/* 添加客户端对话框 */}
        <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t("dashboard.addClient")}</DialogTitle>
              <DialogDescription>
                {t("dashboard.form.namePlaceholder")}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="name" className="text-right">
                  {t("dashboard.form.name")}
                </Label>
                <Input
                  id="name"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  className="col-span-3"
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="email" className="text-right">
                  {t("dashboard.form.email")}
                </Label>
                <Input
                  id="email"
                  name="email"
                  type="email"
                  value={formData.email}
                  onChange={handleInputChange}
                  className="col-span-3"
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="status" className="text-right">
                  {t("dashboard.form.status")}
                </Label>
                <Select
                  value={formData.status}
                  onValueChange={handleStatusChange}
                >
                  <SelectTrigger className="col-span-3">
                    <SelectValue placeholder={t("dashboard.form.status")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="active">
                      {t("dashboard.status.active")}
                    </SelectItem>
                    <SelectItem value="inactive">
                      {t("dashboard.status.inactive")}
                    </SelectItem>
                    <SelectItem value="pending">
                      {t("dashboard.status.pending")}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="notes" className="text-right">
                  {t("dashboard.form.notes")}
                </Label>
                <Textarea
                  id="notes"
                  name="notes"
                  value={formData.notes}
                  onChange={handleInputChange}
                  className="col-span-3"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setIsAddDialogOpen(false)}
              >
                {t("dashboard.form.cancel")}
              </Button>
              <Button onClick={handleAddClient}>
                {t("dashboard.form.save")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* 编辑客户端对话框 */}
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t("dashboard.editClient")}</DialogTitle>
              <DialogDescription>
                {t("dashboard.form.namePlaceholder")}
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="edit-name" className="text-right">
                  {t("dashboard.form.name")}
                </Label>
                <Input
                  id="edit-name"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  className="col-span-3"
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="edit-email" className="text-right">
                  {t("dashboard.form.email")}
                </Label>
                <Input
                  id="edit-email"
                  name="email"
                  type="email"
                  value={formData.email}
                  onChange={handleInputChange}
                  className="col-span-3"
                />
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="edit-status" className="text-right">
                  {t("dashboard.form.status")}
                </Label>
                <Select
                  value={formData.status}
                  onValueChange={handleStatusChange}
                >
                  <SelectTrigger className="col-span-3">
                    <SelectValue placeholder={t("dashboard.form.status")} />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="active">
                      {t("dashboard.status.active")}
                    </SelectItem>
                    <SelectItem value="inactive">
                      {t("dashboard.status.inactive")}
                    </SelectItem>
                    <SelectItem value="pending">
                      {t("dashboard.status.pending")}
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="edit-notes" className="text-right">
                  {t("dashboard.form.notes")}
                </Label>
                <Textarea
                  id="edit-notes"
                  name="notes"
                  value={formData.notes}
                  onChange={handleInputChange}
                  className="col-span-3"
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setIsEditDialogOpen(false)}
              >
                {t("dashboard.form.cancel")}
              </Button>
              <Button onClick={handleEditClient}>
                {t("dashboard.form.save")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* 删除客户端对话框 */}
        <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t("dashboard.deleteConfirm.title")}</DialogTitle>
              <DialogDescription>
                {t("dashboard.deleteConfirm.message")}
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <Button
                variant="outline"
                onClick={() => setIsDeleteDialogOpen(false)}
              >
                {t("dashboard.deleteConfirm.cancel")}
              </Button>
              <Button variant="destructive" onClick={handleDeleteClient}>
                {t("dashboard.deleteConfirm.confirm")}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </MainLayout>
  );
}
