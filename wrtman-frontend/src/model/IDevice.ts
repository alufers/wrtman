export default interface IDevice {
  hasDHCP: boolean;
  available: boolean;
  address: string;
  hostname: string;
  uptime: string;
  uptimeSeconds: number;
  vendor: string | null;
}
