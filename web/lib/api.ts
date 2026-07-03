// API wrapper for the Go TR-069 ACS. Base URL is overridable via NEXT_PUBLIC_ACS_API.
export const ACS_API =
  process.env.NEXT_PUBLIC_ACS_API ?? 'http://localhost:7547';

export interface Device {
  serialNumber: string;
  manufacturer: string;
  oui: string;
  productClass: string;
  ip: string;
  softwareVersion: string;
  hardwareVersion: string;
  lastEvent: string;
  lastInformAt: string;
  informCount: number;
  createdAt: string;
}

export async function fetchDevices(): Promise<Device[]> {
  const res = await fetch(`${ACS_API}/api/devices`, { cache: 'no-store' });
  if (!res.ok) {
    throw new Error(`ACS API ${res.status}`);
  }
  return res.json();
}
