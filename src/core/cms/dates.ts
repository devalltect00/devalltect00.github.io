/** Validate real historical calendar values without inventing missing month/day or offsetting early years. */
export function isProjectDate(value: string): boolean {
  const m = /^(\d{4})(?:-(\d{2})(?:-(\d{2}))?)?$/.exec(value);
  if (!m || Number(m[1]) < 1) return false;
  if (!m[2]) return true;
  const y = Number(m[1]),
    month = Number(m[2]);
  if (month < 1 || month > 12) return false;
  if (!m[3]) return true;
  const leap = y % 4 === 0 && (y % 100 !== 0 || y % 400 === 0),
    day = Number(m[3]);
  return (
    day >= 1 &&
    day <=
      [31, leap ? 29 : 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31][month - 1]
  );
}
