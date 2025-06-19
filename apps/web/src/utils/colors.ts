export function isColorLight(color: string) {
  let r, g, b;

  if (color.match(/^rgb/)) {
    const rgb = color.match(/^rgba?\((\d+),\s*(\d+),\s*(\d+)(?:,\s*(\d+(?:\.\d+)?))?\)$/);
    if (!rgb) return null;
    r = +rgb[1];
    g = +rgb[2];
    b = +rgb[3];
  } else {
    const hex = color.match(/^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i);
    if (!hex) return null;
    r = parseInt(hex[1], 16);
    g = parseInt(hex[2], 16);
    b = parseInt(hex[3], 16);
  }

  if (!r && !g && !b) return null;

  const hsp = Math.sqrt(0.299 * (r * r) + 0.587 * (g * g) + 0.114 * (b * b));

  return hsp > 127.5;
}
