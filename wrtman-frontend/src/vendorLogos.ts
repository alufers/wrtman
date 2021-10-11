export const vendorLogos = [
  {
    filename: `espressif.png`,
    vendorStrings: ["Espressif"],
  },
  {
    filename: `google.png`,
    vendorStrings: ["Google"],
  },
  {
    filename: `samsung.png`,
    vendorStrings: ["Samsung"],
    hostnameStrings: ["Samsung", "Galaxy"],
  },
  {
    filename: `raspberry.png`,
    vendorStrings: ["Raspberry Pi"],
  },
  {
    filename: `sony.png`,
    vendorStrings: ["Sony "],
  },
  {
    filename: `xiaomi.png`,
    vendorStrings: ["Xiaomi"],
  },
  {
    filename: `intel.png`,
    vendorStrings: ["Intel"],
  },
  {
    filename: `mikrotik.png`,
    vendorStrings: ["Routerboard.com"],
  },
  {
    filename: `realtek.png`,
    vendorStrings: ["Realtek"],
  },
  {
    filename: `hp.png`,
    vendorStrings: ["HP Inc.", "Hewlett Packard"],
  },
  {
    filename: `tp-link.png`,
    vendorStrings: ["Tp-Link"],
  },
  {
    filename: `asus.png`,
    vendorStrings: ["ASUSTek"],
  },
  {
    filename: `dell.png`,
    vendorStrings: ["Dell Inc."],
  },
  {
    filename: `apple.png`,
    vendorStrings: ["Apple"],
    hostnameStrings: ["Apple", "iPad"],
  },
  {
    filename: `ublox.png`,
    vendorStrings: ["u-blox"],
  },
  {
    filename: `huawei.png`,
    vendorStrings: ["Huawei Technologies Co."],
  },
];

export function getVendorLogoURL(row: { vendor: string; hostname: string }) {
  let filename = null;

  filename = vendorLogos.find((vl) =>
    vl?.vendorStrings?.some((s) =>
      row.vendor.toLowerCase().includes(s.toLowerCase())
    )
  )?.filename;

  if (!filename) {
    filename = vendorLogos.find((vl) =>
      vl?.hostnameStrings?.some((s) =>
        row.hostname.toLowerCase().includes(s.toLowerCase())
      )
    )?.filename;
  }
  if (filename) {
    return `/assets/logos/${filename}`;
  }
  return null;
}
