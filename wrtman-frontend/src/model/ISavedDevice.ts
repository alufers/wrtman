export default interface ISavedDevice {
  id: string;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
  macAddress: string;
  hostname: string;
  wirelessNetwork: string;
  wirelessAPName: string;
  lastSeen: string;
  note: string;
  vendor: string;
}
