import type { SvelteComponent } from 'svelte';

export default interface ColumnModel<R> {
  key: keyof R;
  label: string;
  sortFunc?: (a: R, b: R, column: ColumnModel<R>) => any | null;
  filterable?: boolean;
  component?: (typeof SvelteComponent) | (typeof SvelteComponent)[] | null;
}
