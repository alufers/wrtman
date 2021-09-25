export function formatDuration(durationSeconds: number): string {
  let inThePast = false;
  if (durationSeconds < 0) {
    durationSeconds = -durationSeconds;
    inThePast = true;
  }
  let result = "";
  let days = Math.floor(durationSeconds / (24 * 60 * 60));
  if (days > 0) {
    durationSeconds %= 24 * 60 * 60;
    result += days.toString().padStart(2, "0") + "d ";
  }
  let hours = Math.floor(durationSeconds / (60 * 60));

  durationSeconds %= 60 * 60;
  result += hours.toString().padStart(2, "0") + "h ";

  let minutes = Math.floor(durationSeconds / 60);

  durationSeconds %= 60;
  result += minutes.toString().padStart(2, "0") + "m ";
  result += durationSeconds.toString().padStart(2, "0") + "s";
  if (inThePast) {
    result += " ago";
  }
  return result;
}
