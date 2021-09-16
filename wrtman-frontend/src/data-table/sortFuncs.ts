import type ColumnModel from "./ColumnModel";

export function sortText(a: any, b: any, c: ColumnModel<any>) {
  return a[c.key].localeCompare(b[c.key]);
}

export function sortDates(a: any, b: any, c: ColumnModel<any>) {
  return a[c.key].getTime() - b[c.key].getTime();
}
