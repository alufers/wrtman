import type ColumnModel from "./ColumnModel";

export function sortText(a: any, b: any, c: ColumnModel<any>) {
  let aTxt = a[c.key] || "zzzzzzzzzzzzzzzz";
  let bTxt = b[c.key] || "zzzzzzzzzzzzzzzz";
  return aTxt.localeCompare(bTxt);
}

export function sortDates(a: any, b: any, c: ColumnModel<any>) {
  if (!a[c.key]) return -1;
  if (!a[c.key]) return 1;
  return a[c.key].getTime() - b[c.key].getTime();
}

export function sortNumbers(a: any, b: any, c: ColumnModel<any>) {
  let aNum = parseFloat(a[c.key]);
  if (isNaN(aNum)) {
    aNum = -Infinity;
  }
  let bNum = parseFloat(b[c.key]);
  if (isNaN(bNum)) {
    bNum = -Infinity;
  }
  return aNum - bNum;
}
