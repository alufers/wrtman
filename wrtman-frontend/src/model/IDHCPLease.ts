export default interface IDHCPLease {
  macAddress: string;
  ipAddress: string;
  hostname: string;
  expiryTime: string;
  vendor: string;
  wirelessNetworkType: string | null;
  ssid: string | null;
  signalStrength: number | null;
  apHostname: string | null;
}
